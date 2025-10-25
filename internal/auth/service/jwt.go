package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
	"time"
)

type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *UsersService) newToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(s.secretDur).Unix()

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *AuthService) validateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, models.ErrInvalidToken
	}

	if time.Now().After(time.Unix(claims.ExpiresAt.Unix(), 0)) {
		return nil, models.ErrTokenExpired
	}

	return claims, nil
}
