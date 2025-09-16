package mux

import (
	"context"
	"encoding/xml"
	"xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/stanza"
)

type messageHandler struct {
	typ     stanza.MessageType
	handler xmpp.Handler
}

func Message(typ stanza.MessageType, handler xmpp.Handler) xmpp.Handler {
	return &messageHandler{
		typ:     typ,
		handler: handler,
	}
}

func (m *messageHandler) HandleXMPP(
	ctx context.Context,
	t xmlstream.TokenReadEncoder,
	start *xml.StartElement,
) (bool, error) {
	if start.Name.Local != "message" {
		return false, nil // Not a message, skip
	}

	msg, err := stanza.NewMessage(*start)
	if err != nil {
		return false, err
	}

	if msg.Type != m.typ {
		return false, nil // Not the expected message type, skip
	}

	return m.handler.HandleXMPP(ctx, t, start)
}
