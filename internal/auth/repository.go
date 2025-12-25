package auth

import (
	"context"
	"time"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUserName(ctx context.Context, username string) (bool, error)

	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}

type PgRepo struct{}

func (r *PgRepo) Create(ctx context.Context, u *models.User) (*models.User, error) {
	result := database.DB.WithContext(ctx).Create(u)
	return u, result.Error
}

func (r *PgRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := database.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}

func (r *PgRepo) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	var count int64
	result := database.DB.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0, result.Error
}

func (r *PgRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var u *models.User
	result := database.DB.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Find(&u)
	return u, result.Error
}

func (r *PgRepo) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return database.DB.WithContext(ctx).Create(token).Error
}

func (r *PgRepo) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := database.DB.WithContext(ctx).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&rt).Error
	return &rt, err
}

func (r *PgRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	return database.DB.WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true).Error
}
