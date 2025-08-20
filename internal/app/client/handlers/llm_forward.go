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

type LlmForwardHandler struct {
	loggerProvider *providers.LoggerProvider
	session        ports.XMPPSession
	llmService     ports.LLMService
}

func NewLlmForwardHandler(
	loggerProvider *providers.LoggerProvider,
	session ports.XMPPSession,
	llmService ports.LLMService,
) xmpp.Handler {
	return &LlmForwardHandler{
		loggerProvider: loggerProvider,
		session:        session,
		llmService:     llmService,
	}
}

func (h *LlmForwardHandler) HandleXMPP(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error) {
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
	logger.Debug("incoming message", ports.Fields{"from": msg.From.String(), "body": msg.Body})

	if msg.Body == "" {
		return false, nil
	}
	go h.askLLM(ctx, msg)

	// TODO DRY sending logic
	xmlBytes, _ := xml.Marshal(entities.ComposingMessage{
		Message: stanza.Message{
			To: msg.From,
		},
	})
	reader := bytes.NewReader(xmlBytes)
	return true, h.session.Send(ctx, xml.NewDecoder(reader))
}

func (h *LlmForwardHandler) askLLM(ctx context.Context, msg entities.MessageBody) {
	logger := h.loggerProvider.Value(ctx)

	response, err := h.llmService.GetChatCompletion(ctx, ports.ChatCompletionRequest(msg.Body))
	if err != nil {
		logger.Error("error talking to llm", ports.Fields{"error": err})
		return
	}

	reply := entities.MessageBody{
		Message: stanza.Message{
			Type: stanza.ChatMessage,
			From: msg.To,
			To:   msg.From.Bare(),
		},
		Body: string(response),
	}
	logger.Debug("llm reply", ports.Fields{"body": reply.Body})

	// TODO DRY sending logic
	xmlBytes, _ := xml.Marshal(reply)
	reader := bytes.NewReader(xmlBytes)
	err = h.session.Send(ctx, xml.NewDecoder(reader))
	if err != nil {
		logger.Error("error sending reply", ports.Fields{"error": err})
	}
}
