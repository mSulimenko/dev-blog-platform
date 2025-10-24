package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/dto"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UsersService struct {
	usersRepo UsersRepository
	log       *zap.SugaredLogger
}

func NewUsersService(usersRepo UsersRepository, logger *zap.SugaredLogger) *UsersService {
	return &UsersService{
		usersRepo: usersRepo,
		log:       logger,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, userReq *dto.UserCreateRequest) (string, error) {
	const op = "users.CreateUser"
	s.log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("generating password hash", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	user := &models.User{
		Username:     userReq.Username,
		Email:        userReq.Email,
		PasswordHash: string(passHash),
	}
	err = s.usersRepo.CreateUser(ctx, user)

	if err != nil {
		s.log.Error("failed to create user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil

}

func (s *UsersService) GetUser(ctx context.Context, id string) (*dto.UserResp, error) {
	const op = "users.GetUser"
	s.log.Infow("getting user", "userID", id)

	if id == "" {
		return nil, fmt.Errorf("%s: %w", op, models.ErrInvalidUserID)
	}

	user, err := s.usersRepo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Errorw("failed to get user", "userID", id, "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userResp := &dto.UserResp{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return userResp, nil
}

func (s *UsersService) ListUsers(ctx context.Context) ([]*dto.UserResp, error) {
	const op = "users.ListUsers"
	s.log.Info("listing all users")

	users, err := s.usersRepo.ListUsers(ctx)
	if err != nil {
		s.log.Errorw("failed to list users", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var usersResp []*dto.UserResp
	for _, user := range users {
		usersResp = append(usersResp, &dto.UserResp{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return usersResp, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, id string, userReq *dto.UserUpdateRequest) error {
	const op = "users.UpdateUser"
	s.log.Infow("updating user", "userID", id)

	if id == "" {
		return fmt.Errorf("%s: %w", op, models.ErrInvalidUserID)
	}

	existingUser, err := s.usersRepo.GetUserByID(ctx, id)
	if err != nil {
		s.log.Errorw("user not found", "userID", id, "error", err)
		return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
	}

	if userReq.Username != nil {
		existingUser.Username = *userReq.Username
	}
	if userReq.Email != nil {
		existingUser.Email = *userReq.Email
	}
	if userReq.Password != nil {
		passHash, err := bcrypt.GenerateFromPassword([]byte(*userReq.Password), bcrypt.DefaultCost)
		if err != nil {
			s.log.Errorw("generating password hash", "error", err)
			return fmt.Errorf("%s: %w", op, err)
		}
		existingUser.PasswordHash = string(passHash)
	}

	err = s.usersRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		s.log.Errorw("failed to update user", "userID", id, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *UsersService) DeleteUser(ctx context.Context, id string) error {
	const op = "users.DeleteUser"
	s.log.Infow("deleting user", "userID", id)

	if id == "" {
		return fmt.Errorf("%s: %w", op, models.ErrInvalidUserID)
	}

	err := s.usersRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			s.log.Errorw("user not found", "userID", id, "error", err)
			return fmt.Errorf("%s: %w", op, models.ErrUserNotFound)
		}
		s.log.Errorw("failed to delete user", "userID", id, "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
