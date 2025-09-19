package handlers

import (
	"context"
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/xmpp/stanza"
)

type EchoHandler struct {
	loggerProvider *providers.LoggerProvider
}

var _ xmpp.Handler = &EchoHandler{}

func NewEchoHandler(loggerProvider *providers.LoggerProvider) *EchoHandler {
	return &EchoHandler{
		loggerProvider: loggerProvider,
	}
}

func (h *EchoHandler) HandleXMPP(ctx context.Context, t xml.TokenReader, w xmpp.StreamWriter) (bool, error) {
	logger := h.loggerProvider.Value(ctx)

	d := xml.NewTokenDecoder(t)

	msg := entities.MessageBody{}
	err := d.Decode(&msg)
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

	return true, w.Write(xmpp.Message(reply))
}
