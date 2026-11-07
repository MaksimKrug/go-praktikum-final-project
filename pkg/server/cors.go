package server

import (
	"net/http"

	"github.com/rs/cors"
)

// WithCORS applies one CORS policy to all routes, before authentication.
func WithCORS(next http.Handler, allowedOrigins []string) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           600,
	}).Handler(next)
}
