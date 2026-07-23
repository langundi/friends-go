package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/langundi/friends-go/internal/services"
	"github.com/langundi/friends-go/internal/utils"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
)

var (
	ErrNoHeader      = errors.New("Authorization header required.")
	ErrFormat        = errors.New("Invalid authorization format.")
	ErrInvalidToken  = errors.New("Invalid or expired token.")
	ErrInvalidClaims = errors.New("Invalid token claims.")
)

func AuthMiddleware(authSerice *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.UnauthorizedError(w, r, ErrNoHeader)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.UnauthorizedError(w, r, ErrFormat)
			}

			tokenString := parts[1]

			// Validate token
			claims, err := authSerice.ValidateToken(tokenString)
			if err != nil {
				utils.UnauthorizedError(w, r, err)
				return
			}

			// Extract user ID from claims
			userID, ok := claims["sub"]
			if !ok {
				utils.UnauthorizedError(w, r, ErrInvalidClaims)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(UserIDKey).(float64)
	return int64(userID), ok
}
