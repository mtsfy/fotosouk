package transformer

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
)

type Transformer interface {
	Crop(ctx context.Context, imgData []byte, width, height int, format string) ([]byte, error)
	Resize(ctx context.Context, imgData []byte, width, height int, format string) ([]byte, error)
	Rotate(ctx context.Context, imgData []byte, degrees int, format string) ([]byte, error)
	Grayscale(ctx context.Context, imgData []byte, format string) ([]byte, error)
}

type ImageTransformer struct{}

func (t *ImageTransformer) Crop(ctx context.Context, imgData []byte, width, height int, format string) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}
	cropped := imaging.CropCenter(img, width, height)
	return encodeImage(cropped, format)
}

func (t *ImageTransformer) Resize(ctx context.Context, imgData []byte, width, height int, format string) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	return encodeImage(resized, format)
}

func (t *ImageTransformer) Rotate(ctx context.Context, imgData []byte, degrees int, format string) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	degrees = degrees % 360
	if degrees < 0 {
		degrees += 360
	}

	var rotated image.Image
	switch degrees {
	case 90:
		rotated = imaging.Rotate90(img)
	case 180:
		rotated = imaging.Rotate180(img)
	case 270:
		rotated = imaging.Rotate270(img)
	case 0:
		rotated = img
	default:
		rotated = imaging.Rotate(img, float64(degrees), color.Transparent)
	}

	return encodeImage(rotated, format)
}

func (t *ImageTransformer) Grayscale(ctx context.Context, imgData []byte, format string) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}
	gray := imaging.Grayscale(img)
	return encodeImage(gray, format)
}

func encodeImage(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	format = strings.ToLower(format)

	switch format {
	case "png", "image/png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	case "jpeg", "jpg", "image/jpeg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
		return buf.Bytes(), err
	default:
		return nil, errors.New("unsupported format: use png or jpeg")
	}
}

func GetImageSize(imgData []byte) (width, height int, err error) {
	img, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return 0, 0, err
	}
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}
