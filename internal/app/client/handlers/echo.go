package handlers

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/stanza"
)

type EchoHandler struct {
	loggerProvider *providers.LoggerProvider
	session        ports.XMPPSession
}

var _ xmpp.Handler = &EchoHandler{}

func NewEchoHandler(loggerProvider *providers.LoggerProvider, session ports.XMPPSession) *EchoHandler {
	return &EchoHandler{
		loggerProvider: loggerProvider,
		session:        session,
	}
}

func (h *EchoHandler) HandleXMPP(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error) {
	logger := h.loggerProvider.Value(ctx)

	// TODO DRY decoding logic
	d := xml.NewTokenDecoder(xmlstream.MultiReader(xmlstream.Token(*start), t))
	if _, err := d.Token(); err != nil {
		return false, err
	}

	msg := entities.MessageBody{}
	err := d.DecodeElement(&msg, start)
	if err != nil && err != io.EOF {
		logger.Error("error decoding message", ports.Fields{"error": err})
		return false, err
	}

	if msg.Body == "" {
		return false, nil
	}

	reply := entities.MessageBody{
		Message: stanza.Message{
			Type: stanza.ChatMessage,
			From: msg.To,
			To:   msg.From.Bare(),
		},
		Body: msg.Body,
	}

	logger.Debug("echo", ports.Fields{"to": reply.To, "body": reply.Body})

	// TODO DRY sending logic
	xmlBytes, _ := xml.Marshal(reply)
	reader := bytes.NewReader(xmlBytes)
	return true, h.session.Send(ctx, xml.NewDecoder(reader))
}
