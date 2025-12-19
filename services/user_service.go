package service

import (
	"context"
	"fmt"
	"strings"

	"TMS/models"
	"TMS/repos"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo repos.UserRepoInterface
}

func NewUserService(repo repos.UserRepoInterface) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	user.UUID = uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service: failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	err := s.Repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create user %s: %w", user.Email, err)
	}

	return &user, nil
}

// GetUserByID fetches a user by internal DB ID
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("service: user with id %d not found", id)
		}
		return nil, fmt.Errorf("service: failed to get user by id %d: %w", id, err)
	}

	return user, nil
}

// GetUserByEmail fetches a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("service: user with email %s not found", email)
		}
		return nil, fmt.Errorf("service: failed to get user by email %s: %w", email, err)
	}

	return user, nil
}
