package image

import (
	"context"
	"io"

	"github.com/mtsfy/fotosouk/internal/models"
)

type ImageService struct {
	repo ImageRepository
}

func NewImageService(repo ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) UploadImage(ctx context.Context, userID int, filename string, file io.Reader, fileSize int64, mimeType string) (*models.Image, error) {
	return nil, nil
}
