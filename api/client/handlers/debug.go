package handlers

import (
	"bytes"
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/ports"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/stanza"
)

// DebugHandler logs all incoming XMPP stanzas for debugging purposes
type DebugHandler struct {
	logger ports.Logger
}

func NewDebugHandler(logger ports.Logger) *DebugHandler {
	return &DebugHandler{
		logger: logger,
	}
}

func toString(tr xml.TokenReader) (string, error) {
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	for {
		tok, err := tr.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if err := enc.EncodeToken(tok); err != nil {
			return "", err
		}
	}
	if err := enc.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (h *DebugHandler) HandleXMPP(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	str, err := toString(xmlstream.MultiReader(xmlstream.Token(*start), t))
	if err != nil {
		return err
	}

	h.logger.Debug("Received stanza", ports.Fields{
		"name":    start.Name.Local,
		"content": str,
	})

	return nil
}

func (h *DebugHandler) HandleMessage(msg stanza.Message, t xmlstream.TokenReadEncoder) error {
	str, err := toString(t)
	if err != nil {
		return err
	}
	h.logger.Debug("Received message", ports.Fields{
		"from":    msg.From.String(),
		"to":      msg.To.String(),
		"type":    msg.Type,
		"ID":      msg.ID,
		"lang":    msg.Lang,
		"content": str,
	})
	return nil
}

func (h *DebugHandler) HandlePresence(pres stanza.Presence, t xmlstream.TokenReadEncoder) error {
	h.logger.Debug("Received message", ports.Fields{
		"from": pres.From.String(),
		"to":   pres.To.String(),
		"type": pres.Type,
		"ID":   pres.ID,
		"lang": pres.Lang,
	})
	return nil
}

func (h *DebugHandler) HandleIQ(iq stanza.IQ, t xmlstream.TokenReadEncoder, start *xml.StartElement) error {
	h.logger.Debug("Received message", ports.Fields{
		"from": iq.From.String(),
		"to":   iq.To.String(),
		"type": iq.Type,
		"ID":   iq.ID,
		"lang": iq.Lang,
	})
	return nil
}
