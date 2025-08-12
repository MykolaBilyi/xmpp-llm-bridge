package client

import (
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"

	"mellium.im/xmlstream"
	"mellium.im/xmpp"
	"mellium.im/xmpp/stanza"
)

func NewHandler(config ports.Config, logger ports.Logger) xmpp.Handler {
	return xmpp.HandlerFunc(func(t xmlstream.TokenReadEncoder, start *xml.StartElement) error {

		// This is a workaround for https://mellium.im/issue/196
		// until a cleaner permanent fix is devised (see https://mellium.im/issue/197)
		d := xml.NewTokenDecoder(xmlstream.MultiReader(xmlstream.Token(*start), t))
		if _, err := d.Token(); err != nil {
			return err
		}

		// Ignore anything that's not a message. In a real system we'd want to at
		// least respond to IQs.
		if start.Name.Local != "message" {
			return nil
		}

		msg := entities.MessageBody{}
		err := d.DecodeElement(&msg, start)
		if err != nil && err != io.EOF {
			logger.Error("Error decoding message", ports.Fields{"error": err})
			return nil
		}

		// Don't reflect messages unless they are chat messages and actually have a
		// body.
		// In a real world situation we'd probably want to respond to IQs, at least.
		if msg.Body == "" || msg.Type != stanza.ChatMessage {
			return nil
		}

		reply := entities.MessageBody{
			Message: stanza.Message{
				To: msg.From.Bare(),
			},
			Body: msg.Body,
		}
		logger.Debug("Replying to message", ports.Fields{"id": msg.ID, "to": reply.To, "body": reply.Body})
		err = t.Encode(reply)
		if err != nil {
			logger.Error("Error responding to message", ports.Fields{"id": msg.ID, "error": err})
		}
		return nil
	})
}
