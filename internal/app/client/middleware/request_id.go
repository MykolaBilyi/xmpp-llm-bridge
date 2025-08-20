package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"xmpp-llm-bridge/internal/entities"
	"xmpp-llm-bridge/internal/ports"
	"xmpp-llm-bridge/internal/providers"
	"xmpp-llm-bridge/pkg/xmpp"

	"mellium.im/xmlstream"
)

type RequestIDMiddleware struct {
	handler           xmpp.Handler
	requestIdProvider *providers.RequestIdProvider
	loggerProvider    *providers.LoggerProvider
}

func WithRequestID(handler xmpp.Handler, requestIdProvider *providers.RequestIdProvider, loggerProvider *providers.LoggerProvider) xmpp.Handler {
	return &RequestIDMiddleware{
		handler:           handler,
		requestIdProvider: requestIdProvider,
		loggerProvider:    loggerProvider,
	}
}

func (m *RequestIDMiddleware) HandleXMPP(ctx context.Context, t xmlstream.TokenReadEncoder, start *xml.StartElement) (bool, error) {
	requestId := findRequestID(start)

	ctx = m.requestIdProvider.WithValue(ctx, requestId)
	ctx = m.loggerProvider.WithLogger(ctx, m.loggerProvider.Value(ctx).WithFields(ports.Fields{
		"id": string(requestId),
	}))

	return m.handler.HandleXMPP(ctx, t, start)
}

func findRequestID(start *xml.StartElement) entities.RequestId {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			return entities.RequestId(attr.Value)
		}
	}
	return generateRequestID()
}

func generateRequestID() entities.RequestId {
	bytes := make([]byte, 8) // 16 character hex string
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to a simple counter if crypto/rand fails (shouldn't happen)
		return entities.RequestId("fallback-id")
	}
	return entities.RequestId(hex.EncodeToString(bytes))
}
