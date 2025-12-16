package auth

import (
	"context"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) (*user.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUserName(ctx context.Context, username string) (bool, error)
}

type PgRepo struct{}

func (r *PgRepo) Create(ctx context.Context, u *user.User) (*user.User, error) {
	result := database.DB.WithContext(ctx).Create(u)
	return u, result.Error
}

func (r *PgRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := database.DB.WithContext(ctx).Model(&user.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}

func (r *PgRepo) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	var count int64
	result := database.DB.WithContext(ctx).Model(&user.User{}).Where("username = ?", username).Count(&count)
	return count > 0, result.Error
}
