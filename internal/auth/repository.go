package auth

import (
	"context"
	"time"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUserName(ctx context.Context, username string) (bool, error)

	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
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

func (r *PgRepo) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var u *user.User
	result := database.DB.WithContext(ctx).Model(&user.User{}).Where("username = ?", username).Find(&u)
	return u, result.Error
}

func (r *PgRepo) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

func (r *PgRepo) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := database.DB.WithContext(ctx).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&rt).Error
	return &rt, err
}

func (r *PgRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	return database.DB.WithContext(ctx).
		Model(&RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}
