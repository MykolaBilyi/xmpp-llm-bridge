package middleware

import (
	"__template__/internal/ports"
	"net/http"

	"github.com/rs/cors"
)

func NewCORSMiddleware(config ports.Config) func(next http.Handler) http.Handler {
	config.SetDefault("allowOrigins", []string{"*"})
	config.SetDefault("allowMethods", []string{http.MethodGet, http.MethodPost, http.MethodHead})
	config.SetDefault("allowHeaders", []string{"Origin", "Authorization", "Content-Type"})
	config.SetDefault("exposeHeaders", []string{})
	config.SetDefault("allowCredentials", false)
	config.SetDefault("maxAge", 5)

	return cors.New(cors.Options{
		AllowedOrigins:   config.GetStringSlice("allowOrigins"),
		AllowedMethods:   config.GetStringSlice("allowMethods"),
		AllowedHeaders:   config.GetStringSlice("allowHeaders"),
		ExposedHeaders:   config.GetStringSlice("exposeHeaders"),
		AllowCredentials: config.GetBool("allowCredentials"),
		MaxAge:           int(config.GetDuration("maxAge").Seconds()),
	}).Handler
}
