package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/f4tal-err0r/discord_faas/pkgs/config"
)

type Local struct {
	path string
}

func NewLocal(opts config.Storage) (Local, error) {
	if _, err := os.Stat(opts.Path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(opts.Path, 0755); err != nil {
				return Local{}, fmt.Errorf("failed to create storage directory: %w", err)
			}
		} else {
			return Local{}, fmt.Errorf("failed to stat storage directory: %w", err)
		}
	}
	return Local{
		path: opts.Path,
	}, nil
}

func (l Local) AddArtifact(ctx context.Context, name string, data io.Reader, size int64) error {
	f, err := os.Create(filepath.Join(l.path, name))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (l Local) ListArtifacts(ctx context.Context, path string) ([]string, error) {
	f, err := os.ReadDir(filepath.Join(l.path, path))
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	var files []string
	for _, file := range f {
		files = append(files, file.Name())
	}
	return files, nil
}

func (l Local) GetArtifact(ctx context.Context, name string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(l.path, name))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return f, nil
}

func (l Local) DeleteArtifact(ctx context.Context, name string) error {
	if err := os.Remove(filepath.Join(l.path, name)); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
