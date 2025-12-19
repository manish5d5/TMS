package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int64     `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	if u.Name == "" {
		return errors.New("name is required")
	}

	if len(u.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple RFC-compliant regex (good enough for backend validation)
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("password must contain at least one number")
	}

	if !strings.ContainsAny(password, "!@#$%^&*") {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
