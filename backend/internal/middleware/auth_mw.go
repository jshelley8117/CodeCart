package middleware

import (
	"context"
	"net/http"
	"strings"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/jshelley8117/CodeCart/internal/common"
	"github.com/jshelley8117/CodeCart/internal/utils"
	"go.uber.org/zap"
)

func AuthMiddleware(fbAuth *firebaseauth.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			z := utils.FromContext(r.Context(), zap.NewNop())
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or malformed authorization header", http.StatusUnauthorized)
				return
			}
			rawToken := strings.TrimPrefix(authHeader, "Bearer ")

			// VerifyIDToken checks signature, expiry, and audience automatically
			token, err := fbAuth.VerifyIDToken(r.Context(), rawToken)
			if err != nil {
				z.Error("firebase token verification failed", zap.Error(err))
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract role custom claim if present; defaults to empty string
			role, _ := token.Claims["role"].(string)

			ctx := context.WithValue(r.Context(), common.ContextKeyFirebaseUID, token.UID)
			ctx = context.WithValue(ctx, common.ContextKeyRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(common.ContextKeyRole).(string)
		if !ok || role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
