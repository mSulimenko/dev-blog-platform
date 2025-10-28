package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/dto"
	authv1 "github.com/mSulimenko/dev-blog-platform/protos/gen/go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api authv1.AuthClient
	log *zap.SugaredLogger
}

func NewAuthClient(ctx context.Context, log *zap.SugaredLogger, addr string, timeout time.Duration, retriesCount uint) (*Client, error) {
	log.Infow("Creating gRPC auth client", "addr", addr, "timeout", timeout, "retries", retriesCount)

	retryOpts := []retry.CallOption{
		retry.WithMax(retriesCount),
		retry.WithPerRetryTimeout(timeout),
		retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		log.Errorw("Failed to create gRPC client", "addr", addr, "error", err)
		return nil, fmt.Errorf("create client conn: %w", err)
	}

	log.Info("gRPC auth client created successfully")
	return &Client{
		api: authv1.NewAuthClient(cc),
		log: log,
	}, nil
}

func (c *Client) Validate(ctx context.Context, token string) (*dto.ValidationResp, error) {
	grpcResp, err := c.api.Validate(ctx, &authv1.ValidateRequest{Token: token})
	if err != nil {
		c.log.Errorw("failed to validate token", "error", err)
		return nil, fmt.Errorf("failed validation token: %w", err)
	}

	resp := dto.ValidationResp{
		Valid:  grpcResp.Valid,
		UserId: grpcResp.UserId,
		Role:   grpcResp.Role,
	}

	return &resp, nil

}
