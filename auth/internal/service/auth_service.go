package service

import (
	"auth/internal/repository"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(username, password string) (string, error)
	ValidateToken(token string) (string, error)
}

type authService struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	tokenService TokenService
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, tokenService TokenService) AuthService {
	return &authService{userRepo: userRepo, sessionRepo: sessionRepo, tokenService: tokenService}
}

func (s *authService) Login(username, password string) (string, error) {
	u, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	hash := sha256.Sum256([]byte(password))
	if hex.EncodeToString(hash[:]) != u.Password {
		return "", errors.New("invalid credentials")
	}

	jwt, err := s.tokenService.GenerateJWT(u.ID)
	if err != nil {
		return "", errors.New("Internal Server Error")
	}

	opaqueToken := uuid.New().String()
	if err := s.sessionRepo.Save(opaqueToken, jwt, u.ID); err != nil {
		log.Println("Error saving token:", err)
		return "", err
	}
	return opaqueToken, nil
}

func (s *authService) ValidateToken(opaqueToken string) (string, error) {
	jwt, userID, err := s.sessionRepo.Get(opaqueToken)
	if err != nil {
		return "", errors.New("invalid token")
	}

	parsed, _ := s.tokenService.ValidateJWT(jwt)
	if parsed != nil {
		_ = s.sessionRepo.Refresh(opaqueToken)
		return jwt, nil
	}

	newJWT, err := s.tokenService.GenerateJWT(userID)
	if err != nil {
		return "", err
	}

	_ = s.sessionRepo.Save(opaqueToken, newJWT, userID)
	return newJWT, nil
}
