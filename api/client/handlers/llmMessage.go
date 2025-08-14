package handlers

import (
	"context"
	"encoding/xml"
	"io"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"

	"mellium.im/xmlstream"
	"mellium.im/xmpp/mux"
	"mellium.im/xmpp/stanza"
)

type LlmMessageHandler struct {
	llmClient ports.LLMClient
	logger    ports.Logger
}

var _ mux.MessageHandler = &LlmMessageHandler{}

func NewLlmMessageHandler(llmClient ports.LLMClient, logger ports.Logger) *LlmMessageHandler {
	return &LlmMessageHandler{
		llmClient: llmClient,
		logger:    logger,
	}
}

func (h *LlmMessageHandler) HandleMessage(m stanza.Message, t xmlstream.TokenReadEncoder) error {
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

	response, err := h.llmClient.GetChatCompletion(context.Background(), ports.ChatCompletionRequest(msg.Body))
	if err != nil {
		h.logger.Error("error talking to llm", ports.Fields{"error": err})
		return nil
	}

	reply := entities.MessageBody{
		Message: stanza.Message{
			To: from,
		},
		Body: string(response),
	}
	h.logger.Debug("Replying to message", ports.Fields{"id": msg.ID, "to": reply.To, "body": reply.Body})
	return t.Encode(reply)
}
