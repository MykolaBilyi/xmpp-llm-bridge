package client

import (
	"xmpp-llm-bridge/internal/app/client/handlers"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"
	"xmpp-llm-bridge/pkg/xmpp/mux"
	"xmpp-llm-bridge/pkg/xmpp/stanza"

	stnz "mellium.im/xmpp/stanza"
)

func NewHandler(
	config ports.Config,
	loggerProvider *providers.LoggerProvider,
	requestIdProvider *providers.RequestIdProvider,
	session ports.XMPPSession,
	llmService ports.LLMService,
) xmpp.Handler {
	// TODO add error handler middleware
	// TODO add whitelist middleware
	// TODO add rate limiting middleware
	muxHandler := mux.New(
		// FIXME use own constants for stanza types
		stnz.NSClient,
		mux.Message(stanza.ChatMessage, handlers.NewEchoHandler(loggerProvider)),
		// mux.Message(stanza.ChatMessage, handlers.NewLlmForwardHandler(loggerProvider, session, llmService)),
		handlers.NewDebugHandler(loggerProvider),
	)

	// return middleware.WithRequestID(muxHandler, requestIdProvider, loggerProvider)
	return muxHandler
}
