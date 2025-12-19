package repos

import (
	"TMS/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {

	query := `
		INSERT INTO users (uuid, name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.DB.QueryRow(ctx, query,
		user.UUID, user.Name, user.Email, user.PasswordHash, user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", user.Email, err)
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, uuid, name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.UUID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" { // pgx returns this string for no rows
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, uuid, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.UUID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email %s: %w", email, err)
	}

	return user, nil
}
