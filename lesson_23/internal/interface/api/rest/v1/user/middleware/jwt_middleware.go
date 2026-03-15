package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/artsadert/lesson_23/internal/domain/entities"
	"github.com/go-chi/jwtauth/v5"
	// Required for *jwt.Token
)

// DualVerifier creates a middleware that validates JWT tokens against both
// new and old access secrets. It accepts tokens signed with either secret.
func DualVerifier(config *entities.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			tokenStr := extractToken(r)
			if tokenStr == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			// Try new secret first
			token, err := config.JWTAccessSecret.Decode(tokenStr)
			if err == nil {
				// Token is valid with new secret
				exp, _ := token.Expiration()
				iat, _ := token.IssuedAt()
				log.Println(exp, iat)
				if exp.Before(time.Now()) {
					http.Error(w, "access token expired", http.StatusUnauthorized)
					return
				}

				ctx := jwtauth.NewContext(r.Context(), token, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// If new secret fails, try old secret
			token, err = config.JWTOLDAccessSecret.Decode(tokenStr)
			if err == nil {
				// Token is valid with old secret
				exp, _ := token.Expiration()
				if exp.Before(time.Now()) {
					http.Error(w, "access token expired", http.StatusUnauthorized)
					return
				}

				ctx := jwtauth.NewContext(r.Context(), token, nil)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Both secrets failed
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		})
	}
}

// RefreshVerifier validates refresh tokens (only against refresh secret)
func RefreshVerifier(config *entities.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			token, err := config.JWTRefreshSecret.Decode(tokenStr)
			if err != nil {
				http.Error(w, "invalid refresh token", http.StatusUnauthorized)
				return
			}

			exp, _ := token.Expiration()
			if exp.Before(time.Now()) {
				http.Error(w, "refresh token expired", http.StatusUnauthorized)
				return
			}

			ctx := jwtauth.NewContext(r.Context(), token, nil)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken retrieves the Bearer token from the Authorization header
func extractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return ""
	}
	parts := strings.Split(bearer, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
