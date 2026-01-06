package image

import (
	"context"

	"github.com/mtsfy/fotosouk/internal/database"
	"github.com/mtsfy/fotosouk/internal/models"
)

type ImageRepository interface {
	Create(ctx context.Context, img *models.Image) (*models.Image, error)
	ListByUser(ctx context.Context, userID int) ([]*models.Image, error)
	GetImageByID(ctx context.Context, userID int, imgID int) (*models.Image, error)
	Update(ctx context.Context, image *models.Image, userID int) (*models.Image, error)
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

func (r *PgRepo) GetImageByID(ctx context.Context, userID int, imgID int) (*models.Image, error) {
	var img *models.Image
	result := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", imgID, userID).
		First(&img)
	return img, result.Error
}

func (r *PgRepo) Update(ctx context.Context, image *models.Image, userID int) (*models.Image, error) {
	var newImg *models.Image
	result := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", image.ID, userID).Updates(models.Image{
		Width:  image.Width,
		Height: image.Height,

		Filename: image.Filename,
		Url:      image.Url,
		MimeType: image.MimeType,
		FileSize: image.FileSize,
		Format:   image.Format,
	}).First(&newImg)
	return newImg, result.Error
}
