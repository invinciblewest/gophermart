package middleware

import (
	"context"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/usecase"
	"net/http"
)

func AuthMiddleware(authUseCase usecase.AuthUseCase) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}

			userID, err := authUseCase.ParseToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, helper.UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
