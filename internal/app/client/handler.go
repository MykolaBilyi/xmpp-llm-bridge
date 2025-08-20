package client

import (
	"xmpp-llm-bridge/internal/app/client/handlers"
	"xmpp-llm-bridge/internal/app/client/middleware"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"
	"xmpp-llm-bridge/pkg/xmpp/mux"

	"mellium.im/xmpp/stanza"
)

func NewHandler(
	config ports.Config,
	loggerProvider *providers.LoggerProvider,
	requestIdProvider *providers.RequestIdProvider,
	session ports.XMPPSession,
	llmService ports.LLMService,
) xmpp.Handler {
	muxHandler := mux.New(
		stanza.NSClient,
		// mux.Message(stanza.ChatMessage, handlers.NewEchoHandler(loggerProvider, session)),
		mux.Message(stanza.ChatMessage, handlers.NewLlmForwardHandler(loggerProvider, session, llmService)),
		handlers.NewDebugHandler(loggerProvider),
	)

	return middleware.WithRequestID(muxHandler, requestIdProvider, loggerProvider)
}
