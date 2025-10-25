package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
	"go.uber.org/zap"
)

const (
	roleUnverified = "unverified"
	roleUser       = "user"
	roleAdmin      = "unverified"
)

type UserProvider interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type AuthService struct {
	userProvider UserProvider
	log          *zap.SugaredLogger
	secret       string
}

func NewAuthService(userProvider UserProvider, logger *zap.SugaredLogger, secret string) *AuthService {
	return &AuthService{
		userProvider: userProvider,
		log:          logger,
		secret:       secret,
	}
}

func (a *AuthService) Auth(ctx context.Context, tokenString string) (userId, role string, err error) {
	const op = "users.CreateUser"

	if tokenString == "" {
		return "", roleUnverified, nil
	}

	token, err := a.validateToken(tokenString)
	if err != nil {
		a.log.Errorw("failed to validate token", "error", err)
		return "", roleUnverified, models.ErrInvalidToken
	}

	a.log.Infow("authorizing user", "userId", token.UserID)

	user, err := a.userProvider.GetUserByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			a.log.Errorw("orphaned token - user not found",
				"user_id", token.UserID,
				"token_email", token.Email)
			return "", "", fmt.Errorf("%s: user not found", op)
		}

		a.log.Errorw("failed to get user", "userId", token.UserID, "error", err)
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return token.UserID, user.Role, nil

}
