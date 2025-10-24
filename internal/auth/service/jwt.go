package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
	"time"
)

func (s *UsersService) newToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(s.secretDur).Unix()

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
