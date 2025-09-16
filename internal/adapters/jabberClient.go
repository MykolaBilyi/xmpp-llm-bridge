package adapters

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"

	myxmpp "xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/sasl"
	"mellium.im/xmpp"
	"mellium.im/xmpp/dial"
	"mellium.im/xmpp/jid"
	"mellium.im/xmpp/stanza"
	"mellium.im/xmpp/stream"
)

type JabberClient struct {
	config         ports.Config
	loggerProvider *providers.LoggerProvider
	server         string
	jid            jid.JID
	session        *xmpp.Session
}

var _ ports.XMPPSession = (*JabberClient)(nil)

func NewJabberClient(
	ctx context.Context,
	config ports.Config,
	loggerProvider *providers.LoggerProvider,
) (*JabberClient, error) {
	config.SetDefault("connectionTimeout", "10s")
	config.SetDefault("lang", "en")

	address := config.GetString("jid")
	j, err := jid.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("error parsing address %q: %w", address, err)
	}

	return &JabberClient{
		config:         config,
		loggerProvider: loggerProvider,
		jid:            j,
		server:         j.Domainpart(),
	}, nil
}

func (j *JabberClient) streamConfig(*xmpp.Session, *xmpp.StreamConfig) xmpp.StreamConfig {
	return xmpp.StreamConfig{
		Lang: j.config.GetString("lang"),
		Features: []xmpp.StreamFeature{
			xmpp.BindResource(),
			xmpp.StartTLS(&tls.Config{
				ServerName: j.server,
				MinVersion: tls.VersionTLS12,
			}),
			xmpp.SASL(
				"",
				j.config.GetString("password"),
				sasl.ScramSha1Plus,
				sasl.ScramSha1,
				sasl.Plain,
			),
		},
	}
}

func (j *JabberClient) Connect(ctx context.Context) error {
	logger := j.loggerProvider.Value(ctx)

	d := dial.Dialer{}
	logger.Debug("connecting", ports.Fields{"host": j.server})
	dialCtx, dialCtxCancel := context.WithTimeout(ctx, j.config.GetDuration("connectionTimeout"))
	connection, err := d.DialServer(dialCtx, "tcp", j.jid, j.server)
	if err != nil {
		dialCtxCancel()
		return fmt.Errorf("error dialing session: %w", err)
	}

	logger.Debug("logging in", ports.Fields{"jid": j.jid.String()})
	j.session, err = xmpp.NewSession(
		ctx,
		j.jid.Domain(),
		j.jid,
		connection,
		0,
		xmpp.NewNegotiator(j.streamConfig),
	)
	dialCtxCancel()
	if err != nil {
		if errors.Is(err, stream.SeeOtherHost) {
			j.server = err.(stream.Error).Content

			logger.Info("see-other-host", ports.Fields{"host": j.server})
			_ = j.session.Close()
			_ = connection.Close()
			return j.Connect(ctx)
		}
		return fmt.Errorf("error establishing a session: %w", err)
	}
	logger.Info("connected to server", ports.Fields{"host": j.server, "jid": j.jid.String()})

	return nil
}

func (j *JabberClient) Handle(ctx context.Context, handler myxmpp.Handler) error {
	err := j.session.Send(ctx, stanza.Presence{Type: stanza.AvailablePresence}.Wrap(nil))
	if err != nil {
		return fmt.Errorf("error sending initial presence: %w", err)
	}

	return j.session.Serve(myxmpp.HandleWithContext(ctx, handler))
}

func (j *JabberClient) Send(ctx context.Context, s myxmpp.Stanza) error {
	xmlBytes, err := xml.Marshal(s.Get())
	if err != nil {
		return fmt.Errorf("error marshalling stanza: %w", err)
	}
	reader := bytes.NewReader(xmlBytes)

	err = j.session.Send(ctx, xml.NewDecoder(reader))
	if err != nil {
		return fmt.Errorf("error sending stanza: %w", err)
	}

	return nil
}

func (j *JabberClient) Close(ctx context.Context) error {
	if j.session == nil {
		return nil
	}
	return j.session.Close()
}
