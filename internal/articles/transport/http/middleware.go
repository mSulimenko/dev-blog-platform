package httphandler

import (
	"context"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/transport/grpc"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware(grpcAuthClient *grpc.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Invalid authorization header format")
				return
			}

			token := parts[1]

			validationResp, err := grpcAuthClient.Validate(r.Context(), token)
			if err != nil {
				h.log.Errorw("Token validation failed", "error", err)
				sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Invalid token")
				return
			}

			if !validationResp.Valid {
				sendError(w, http.StatusUnauthorized, ErrCodeValidation, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", validationResp.UserId)
			ctx = context.WithValue(ctx, "user_role", validationResp.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
