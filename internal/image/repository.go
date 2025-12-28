package image

import (
	"context"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/models"
)

type ImageRepository interface {
	Create(ctx context.Context, img *models.Image) (*models.Image, error)
}

type PgRepo struct{}

func (r *PgRepo) Create(ctx context.Context, img *models.Image) (*models.Image, error) {
	result := database.DB.WithContext(ctx).Create(img)
	return img, result.Error
}
