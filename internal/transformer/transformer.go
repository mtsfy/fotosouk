package transformer

import (
	"context"
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
)

type Transformer interface {
	Crop(ctx context.Context, img image.Image, x, y, width, height int) (image.Image, error)
}

type GoImageTransformer struct{}

func (t *GoImageTransformer) Crop(ctx context.Context, img image.Image, x, y, width, height int) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be positive")
	}
	if x < 0 || y < 0 {
		return nil, errors.New("x and y must be non-negative")
	}

	srcBounds := img.Bounds()
	cropRect := image.Rect(srcBounds.Min.X+x, srcBounds.Min.Y+y, srcBounds.Min.X+x+width, srcBounds.Min.Y+y+height)

	if !cropRect.In(srcBounds) {
		return nil, errors.New("crop rectangle out of bounds")
	}

	if si, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}); ok {
		return si.SubImage(cropRect), nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), img, cropRect.Min, draw.Src)
	return dst, nil
}
