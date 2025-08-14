package adapters

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"xmpp-llm-bridge/internal/ports"

	"mellium.im/sasl"
	"mellium.im/xmpp"
	"mellium.im/xmpp/dial"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
	"mellium.im/xmpp/stream"
)

type JabberClient struct {
	baseCtx       context.Context
	cancelContext context.CancelFunc
	cfg           ports.Config
	logger        ports.Logger
	server        string
	jid           jid.JID
	session       *xmpp.Session
	handler       xmpp.Handler
}

var _ ports.Server = (*JabberClient)(nil)

func NewJabberClient(cfg ports.Config, handler xmpp.Handler, logger ports.Logger) (*JabberClient, error) {
	cfg.SetDefault("connectionTimeout", "10s")
	cfg.SetDefault("lang", "en")

	address := cfg.GetString("jid")
	j, err := jid.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("error parsing address %q: %w", address, err)
	}

	baseCtx, cancel := context.WithCancel(context.Background())

	return &JabberClient{
		cfg:           cfg,
		baseCtx:       baseCtx,
		cancelContext: cancel,
		logger:        logger,
		jid:           j,
		server:        j.Domainpart(),
		handler:       handler,
	}, nil
}

func (j *JabberClient) streamConfig(*xmpp.Session, *xmpp.StreamConfig) xmpp.StreamConfig {
	return xmpp.StreamConfig{
		Lang: j.cfg.GetString("lang"),
		Features: []xmpp.StreamFeature{
			xmpp.BindResource(),
			xmpp.StartTLS(&tls.Config{
				ServerName: j.server,
				MinVersion: tls.VersionTLS12,
			}),
			xmpp.SASL("", j.cfg.GetString("password"), sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
		},
	}
}

func (j *JabberClient) Serve() error {
	d := dial.Dialer{}
	j.logger.Debug("connecting", ports.Fields{"host": j.server})
	dialCtx, dialCtxCancel := context.WithTimeout(j.baseCtx, j.cfg.GetDuration("connectionTimeout"))
	connection, err := d.DialServer(dialCtx, "tcp", j.jid, j.server)
	if err != nil {
		dialCtxCancel()
		return fmt.Errorf("error dialing session: %w", err)
	}
	j.logger.Debug("logging in", ports.Fields{"jid": j.jid.String()})
	j.session, err = xmpp.NewSession(j.baseCtx, j.jid.Domain(), j.jid, connection, 0, xmpp.NewNegotiator(j.streamConfig))
	dialCtxCancel()
	if err != nil {
		if errors.Is(err, stream.SeeOtherHost) {
			j.server = err.(stream.Error).Content

			j.logger.Info("see-other-host", ports.Fields{"host": j.server})
			j.session.Close()
			connection.Close()
			return j.Serve()
		}
		return fmt.Errorf("error establishing a session: %w", err)
	}
	j.logger.Info("connected to server", ports.Fields{"host": j.server, "jid": j.jid.String()})

	err = j.session.Send(j.baseCtx, stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return fmt.Errorf("Error sending initial presence: %w", err)
	}

	return j.session.Serve(j.handler)
}

func (j *JabberClient) Shutdown(ctx context.Context) error {
	j.cancelContext()
	if j.session == nil {
		return nil
	}
	return j.session.Close()
}
