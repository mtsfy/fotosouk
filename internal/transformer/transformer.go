package transformer

import (
	"context"
	"image"
)

type Transformer interface {
	Crop(ctx context.Context, img image.Image, x, y, width, height int) (image.Image, error)
}

type GoImageTransformer struct{}

func (t *GoImageTransformer) Crop(ctx context.Context, img image.Image, x, y, width, height int) (image.Image, error) {

	return nil, nil
}
