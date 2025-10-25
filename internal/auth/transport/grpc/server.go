package grpc

import (
	"context"
	"errors"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
	authv1 "github.com/mSulimenko/dev-blog-platform/protos/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Auth(ctx context.Context, token string) (userId, role string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	authService AuthService
}

func Register(gRPCServer *grpc.Server, auth AuthService) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{authService: auth})
}

func (s *serverAPI) Validate(ctx context.Context, in *authv1.ValidateRequest,
) (*authv1.ValidateResponse, error) {

	userId, role, err := s.authService.Auth(ctx, in.Token)
	if err != nil {
		if errors.Is(err, models.ErrInvalidToken) ||
			errors.Is(err, models.ErrUserNotFound) ||
			errors.Is(err, models.ErrTokenExpired) {
			return &authv1.ValidateResponse{
				Valid: false,
			}, nil
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.ValidateResponse{
		Valid:  true,
		UserId: userId,
		Role:   role,
	}, nil

}
