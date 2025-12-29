package image

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"strings"

	"github.com/mtsfy/fotosouk/internal/models"
	"github.com/mtsfy/fotosouk/internal/storage"
)

type ImageService struct {
	repo    ImageRepository
	storage storage.Storage
}

func NewImageService(repo ImageRepository, storage storage.Storage) *ImageService {
	return &ImageService{repo, storage}
}

func (s *ImageService) UploadImage(ctx context.Context, userID int, filename string, file io.Reader, fileSize int64, mimeType string) (*models.Image, error) {
	fmt.Println(filename, mimeType, fileSize)

	if float64(fileSize) > math.Pow(10, 7) {
		return nil, errors.New("image size is too big")
	}

	mt := strings.ToLower(mimeType)

	if mt != "image/jpeg" && mt != "image/png" {
		return nil, errors.New("unsupported image type")
	}

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, errors.New("unable to decode image file")
	}

	fmt.Println(img.Bounds(), format)

	newImg := &models.Image{
		UserID:   userID,
		Width:    img.Bounds().Dx(),
		Height:   img.Bounds().Dy(),
		Filename: filename,
		Url:      "", // for now
		MimeType: mimeType,
		FileSize: fileSize,
	}

	return s.repo.Create(ctx, newImg)
}
