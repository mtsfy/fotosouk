package image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"math"
	"strings"

	"github.com/mtsfy/fotosouk/internal/models"
	"github.com/mtsfy/fotosouk/internal/storage"
	"github.com/mtsfy/fotosouk/internal/transformer"
)

type ImageService struct {
	repo        ImageRepository
	storage     storage.Storage
	transformer transformer.Transformer
}

func NewImageService(repo ImageRepository, storage storage.Storage, trans transformer.Transformer) *ImageService {
	return &ImageService{repo, storage, trans}
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

type TransformOptions struct {
	Resize struct {
		Width  int
		Height int
	}
	Crop struct {
		Width  int
		Height int
		X      int
		Y      int
	}
	Rotate  int
	Format  string
	Filters struct {
		Grayscale bool
		Sepia     bool
	}
}

func (s *ImageService) TransformImage(ctx context.Context, userID, imageID int, opts *TransformOptions) (*models.Image, error) {
	ogImg, err := s.repo.GetImageByID(ctx, userID, imageID)
	if err != nil {
		return nil, errors.New("image not found")
	}

	imgData, err := s.storage.Download(ctx, ogImg.Url)
	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, errors.New("failed to decode image")
	}

	if opts.Crop.Width > 0 && opts.Crop.Height > 0 {
		img, err = s.transformer.Crop(ctx, img, opts.Crop.X, opts.Crop.Y, opts.Crop.Width, opts.Crop.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to crop: %w", err)
		}
	}

	if opts.Resize.Width > 0 && opts.Resize.Height > 0 {
		img, err = s.transformer.Resize(ctx, img, opts.Resize.Width, opts.Resize.Height)
		if err != nil {
			return nil, fmt.Errorf("failed to resize: %w", err)
		}
	}

	outFormat := opts.Format
	if outFormat == "" {
		outFormat = format
	}

	var buf bytes.Buffer
	switch outFormat {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		ogImg.MimeType = "image/jpeg"
	case "png":
		err = png.Encode(&buf, img)
		ogImg.MimeType = "image/png"
	default:
		return nil, errors.New("unsupported output format")
	}

	if err != nil {
		return nil, err
	}

	_, err = s.storage.Upload(ctx, "images/"+ogImg.Filename, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}

	ogImg.Width = img.Bounds().Dx()
	ogImg.Height = img.Bounds().Dy()
	ogImg.FileSize = int64(buf.Len())

	return s.repo.Update(ctx, ogImg, userID)
}
