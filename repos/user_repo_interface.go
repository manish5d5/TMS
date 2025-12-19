package repos

import (
	"TMS/models"
	"context"
)

type UserRepoInterface interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
