package middleware

import (
	"net/http"
	"xmpp-llm-bridge/pkg/xmpp"
)

type Middleware = func(http.Handler) http.Handler
type XMPPMiddleware = func(xmpp.Handler) xmpp.Handler
