package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

type Storage interface {
	AddArtifact(ctx context.Context, name string, data io.Reader, size int64) error
	GetArtifact(ctx context.Context, name string) (io.ReadCloser, error)
	ListArtifacts(ctx context.Context, path string) ([]string, error)
	DeleteArtifact(ctx context.Context, name string) error
}

func NewStorage(opts config.Storage) (Storage, error) {
	switch opts.Type {
	case "s3":
		return NewS3Client(opts)
	case "local":
		return NewLocal(opts)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", opts.Type)
	}
}
