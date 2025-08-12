package middleware

import (
	"__template__/internal/ports"
	"__template__/pkg/middleware"
	"net/http"
)

type SkipperMiddleware struct {
	middlewares []middleware.Middleware
	skipPath    []string
	next        http.Handler
}

func (m *SkipperMiddleware) skip(r *http.Request) bool {
	for _, path := range m.skipPath {
		if path == r.URL.Path {
			return true
		}
	}
	return false
}

func (m *SkipperMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next := m.next
	if !m.skip(r) {
		for i := len(m.middlewares) - 1; i >= 0; i-- {
			next = m.middlewares[i](next)
		}
	}
	next.ServeHTTP(w, r)
}

func NewSkipperMiddleware(config ports.Config, middlewares ...middleware.Middleware) middleware.Middleware {
	config.SetDefault("skip", []string{})

	return func(next http.Handler) http.Handler {
		return &SkipperMiddleware{
			middlewares: middlewares,
			skipPath:    config.GetStringSlice("skip"),
			next:        next,
		}
	}
}
