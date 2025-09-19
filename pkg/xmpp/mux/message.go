package mux

import (
	"context"
	"encoding/xml"
	myxml "xmpp-llm-bridge/pkg/xml"
	"xmpp-llm-bridge/pkg/xmpp"
	"xmpp-llm-bridge/pkg/xmpp/stanza"
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

func (m *messageHandler) HandleXMPP(ctx context.Context, t xml.TokenReader, w xmpp.StreamWriter) (bool, error) {
	t, start, err := myxml.ExtractStartElement(t)
	if err != nil {
		return false, err
	}
	if stanza.Stanza(start.Name.Local) != stanza.Message {
		return false, nil // Not a message, skip
	}

	var typ stanza.MessageType
	for _, attr := range start.Attr {
		if attr.Name.Local == "type" {
			typ = stanza.MessageType(attr.Value)
			break
		}
	}

	if typ != m.typ {
		return false, nil // Not the expected message type, skip
	}

	return m.handler.HandleXMPP(ctx, t, w)
}
