package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"TMS/models"
	"TMS/repos"

	"TMS/utils/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repos.UserRepoInterface
	Redis    *redis.Client

	jwt *jwt.Jwt
}

func NewAuthService(repo repos.UserRepoInterface, redis *redis.Client, jwt *jwt.Jwt) *AuthService {
	return &AuthService{
		UserRepo: repo,
		Redis:    redis,
		jwt:      jwt,
	}
}

type CookiesModel struct {
	AccessCookie  *http.Cookie
	RefreshCookie *http.Cookie
}

func (c *CookiesModel) ChangeSameSiteForDevelopment(r *http.Request) {
	origin := r.Header.Get("Origin")
	if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
		c.AccessCookie.SameSite = http.SameSiteNoneMode
		c.RefreshCookie.SameSite = http.SameSiteNoneMode
	}
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, *CookiesModel, error) {
	user, err := s.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, fmt.Errorf("invalid password")
	}

	userId := fmt.Sprintf("%d", user.ID)

	accessToken, err := s.jwt.GenerateJwtToken(userId, 15*time.Minute)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwt.GenerateJwtToken(userId, 4*time.Hour)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	accessKey := userId + "---" + accessToken
	if err := s.Redis.Set(ctx, accessKey, userId, 15*time.Minute).Err(); err != nil {
		return nil, nil, fmt.Errorf("Redis failure")
	}
	refrshKey := userId + "---" + refreshToken
	if err := s.Redis.Set(ctx, refrshKey, userId, 4*time.Hour).Err(); err != nil {
		return nil, nil, fmt.Errorf("Redis failure")
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(4 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	cookieModel = &CookiesModel{
		AccessCookie:  accessCookie,
		RefreshCookie: refreshCookie,
	}

	return models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,

		UserID:        user.ID,
		AccessCookie:  accessCookie.String(),
		RefreshCookie: refreshCookie.String(),
	}, cookieModel, nil
}
