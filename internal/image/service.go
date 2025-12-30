package image

import (
	"bytes"
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

func (s *ImageService) UploadImage(ctx context.Context, userID int, filename string, r io.Reader, fileSize int64, mimeType string) (*models.Image, error) {
	fmt.Println(filename, mimeType, fileSize)

	if float64(fileSize) > math.Pow(10, 7) {
		return nil, errors.New("image size is too big")
	}

	mt := strings.ToLower(mimeType)

	if mt != "image/jpeg" && mt != "image/png" {
		return nil, errors.New("unsupported image type")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("unable to decode image file")
	}

	fmt.Println(img.Bounds(), format)
	path := fmt.Sprintf("images/%s", filename)
	url, err := s.storage.Upload(ctx, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	newImg := &models.Image{
		UserID:   userID,
		Width:    img.Bounds().Dx(),
		Height:   img.Bounds().Dy(),
		Filename: filename,
		Url:      url,
		MimeType: mimeType,
		FileSize: int64(len(data)),
	}

	return s.repo.Create(ctx, newImg)
}

func (s *ImageService) GetAllImages(ctx context.Context, userID int) ([]*models.Image, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *ImageService) GetImageDetail(ctx context.Context, userID int, imgID int) (*models.Image, error) {
	return s.repo.GetImageByID(ctx, userID, imgID)
}
