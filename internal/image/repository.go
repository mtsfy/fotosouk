package image

import (
	"context"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/models"
)

type ImageRepository interface {
	Create(ctx context.Context, img *models.Image) (*models.Image, error)
	ListByUser(ctx context.Context, userID int) ([]*models.Image, error)
}

type PgRepo struct{}

func (r *PgRepo) Create(ctx context.Context, img *models.Image) (*models.Image, error) {
	result := database.DB.WithContext(ctx).Create(img)
	return img, result.Error
}

func (r *PgRepo) ListByUser(ctx context.Context, userID int) ([]*models.Image, error) {
	var images []*models.Image
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&images)
	return images, result.Error
}
