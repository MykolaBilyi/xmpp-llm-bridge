package client

import (
	"xmpp-llm-bridge/internal/app/client/handlers"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"
	"xmpp-llm-bridge/pkg/xmpp/mux"

	"mellium.im/xmpp/stanza"
)

func NewHandler(
	config ports.Config,
	loggerProvider *providers.LoggerProvider,
	session ports.XMPPSession,
	llmService ports.LLMService,
) xmpp.Handler {
	// TODO: add requestId middleware
	return mux.New(
		stanza.NSClient,
		// mux.Message(stanza.ChatMessage, handlers.NewEchoHandler(loggerProvider, sessionProvider)),
		mux.Message(stanza.ChatMessage, handlers.NewLlmForwardHandler(loggerProvider, session, llmService)),
		handlers.NewDebugHandler(loggerProvider),
	)
}
