package middleware

import (
	"context"
	"net/http"
	"strings"

	"rmf2-map-server/internal/repository"
)

type orgKey struct{}

func OrgFromContext(ctx context.Context) string {
	org, _ := ctx.Value(orgKey{}).(string)
	return org
}

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
	"http://localhost:3000": true,
}

// noAuthPaths lists routes that do not require a Bearer token.
var noAuthPaths = map[string]bool{
	"/":       true,
	"/health": true,
}

func New(keys repository.APIKeys, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Auth
		if !noAuthPaths[r.URL.Path] {
			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			org, ok := keys[token]
			if !ok || token == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"detail":"unauthorized"}`))
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), orgKey{}, org))
		}

		next.ServeHTTP(w, r)
	})
}
