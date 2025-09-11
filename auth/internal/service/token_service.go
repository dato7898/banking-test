package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateJWT(userID int) (string, error)
	ValidateJWT(token string) (*jwt.MapClaims, error)
}

type tokenService struct {
	secret string
}

func NewTokenService(secret string) TokenService {
	return &tokenService{secret: secret}
}

func (s *tokenService) GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *tokenService) ValidateJWT(token string) (*jwt.MapClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims := parsed.Claims.(jwt.MapClaims)
	return &claims, nil
}
