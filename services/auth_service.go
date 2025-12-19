package service

import (
	"context"
	"fmt"
	"time"

	"TMS/models"
	"TMS/repos"

	"TMS/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repos.UserRepo

	jwt *jwt.Jwt
}

func NewAuthService(repo *repos.UserRepo, jwt *jwt.Jwt) *AuthService {
	return &AuthService{
		UserRepo: repo,
		jwt:      jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	user, err := s.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return models.LoginResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.LoginResponse{}, err
	}
	userId := fmt.Sprintf(user.Id)

	accessToken, err := s.jwt.GenerateJwtToken(userId, 15*time.Minute)
	if err != nil {
		return models.LoginResponse{}, err
	}
	refreshToken, err := s.jwt.GenerateJwtToken(userId, 4*time.Hour)
	if err != nil {
		return models.LoginResponse{}, err
	}
	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userId,
	}, nil
}
