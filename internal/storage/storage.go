package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, path string, file io.Reader) (string, error)
}

type S3Storage struct {
	bucket string
}

func (s *S3Storage) Upload(path string, file io.Reader) (string, error) {
	return "", nil
}
