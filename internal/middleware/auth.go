package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/martian/v3/log"
)

type UserProvider interface {
	VerifyToken(ctx context.Context, idToken string) (string, error)
}

type ctxKey string

const userIDCtxKey ctxKey = "userID"

// AuthMiddleware creates a middleware that authenticates requests using the Authorization header.
func AuthMiddleware(provider UserProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if len(authz) < len(prefix) || !strings.EqualFold(authz[:len(prefix)], prefix) {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}
			idToken := authz[len(prefix):]
			if idToken == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			userID, err := provider.VerifyToken(r.Context(), idToken)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Store userID in context
			ctx := context.WithValue(r.Context(), userIDCtxKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the authenticated user ID from the context.
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(userIDCtxKey).(string)
	if !ok {
		log.Errorf("user ID not found in context")
		return ""
	}
	return userID
}
