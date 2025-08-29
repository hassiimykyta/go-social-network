package auth

import (
	"net/http"
	"strings"
)

type ctxKey string

const CtxUserID ctxKey = "uid"
const CtxRole ctxKey = "role"

func Middleware(s *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "missing bearer", http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
			claims, err := s.Verify(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := WithUserID(r.Context(), claims.UserId)
			if claims.Role != "" {
				ctx = WithRole(ctx, claims.Role)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
