package image

import (
	"context"

	"github.com/mtsfy/fotosouk/internal/models"
)

type ImageRepository interface {
	Create(ctx context.Context, img *models.Image) error
}

type PgRepo struct{}

func (r *PgRepo) Create(ctx context.Context, img *models.Image) error {
	return nil
}
