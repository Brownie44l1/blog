package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Brownie44l1/blog/internal/auth"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"error": "` + message + `"}`))
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization format. Must be 'Bearer <token>'")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			log.Printf("DEBUG: Received token: %s...", tokenString[:50])
			log.Printf("DEBUG: JWT Secret length: %d", len(jwtSecret))

			claims, err := auth.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				log.Printf("❌ Token validation error: %v", err) // Add this
				respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			log.Printf("✅ Middleware: Token validated for user ID: %d", claims.UserID) // Add this

			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)

			r = r.WithContext(ctx)

			log.Printf("✅ Middleware: Calling next handler") // Add this
			next.ServeHTTP(w, r)
			log.Printf("✅ Middleware: Handler completed") // Add this
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int64)
	return userID, ok
}
