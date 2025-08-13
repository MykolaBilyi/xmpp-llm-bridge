package handlers

import (
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
)

type EchoMessageHandler struct {
	logger ports.Logger
}

var _ mux.MessageHandler = &EchoMessageHandler{}

func NewEchoMessageHandler(logger ports.Logger) *EchoMessageHandler {
	return &EchoMessageHandler{
		logger: logger,
	}
}

func (h *EchoMessageHandler) HandleMessage(m stanza.Message, t xmlstream.TokenReadEncoder) error {
	d := xml.NewTokenDecoder(t)
	from := m.From
	if m.Type != stanza.GroupChatMessage {
		from = m.From.Bare()
	}

	msg := entities.MessageBody{}
	err := d.Decode(&msg)
	if err != nil && err != io.EOF {
		h.logger.Error("error decoding message", ports.Fields{"error": err})
		return nil
	}

	reply := entities.MessageBody{
		Message: stanza.Message{
			To: from,
		},
		Body: msg.Body,
	}
	h.logger.Debug("Replying to message", ports.Fields{"id": msg.ID, "to": reply.To, "body": reply.Body})
	return t.Encode(reply)
}
