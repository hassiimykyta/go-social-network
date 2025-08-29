package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalFS struct {
	baseDir    string
	publicBase string
}

func NewLocalFS(baseDir, publicBase string) *LocalFS {
	return &LocalFS{
		baseDir:    filepath.Clean(baseDir),
		publicBase: strings.TrimRight(publicBase, "/"),
	}
}

func (l *LocalFS) FullPath(key string) string {
	clean := filepath.Clean(key)
	return filepath.Join(l.baseDir, clean)
}

func (l *LocalFS) Save(ctx context.Context, key string, r io.Reader, _ int64, _ string) error {
	path := l.FullPath(key)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}

func (l *LocalFS) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	return os.Open(l.FullPath(key))
}

func (l *LocalFS) URL(ctx context.Context, key string) (string, error) {
	u := fmt.Sprintf("%s/media/%s", l.publicBase, strings.ReplaceAll(key, "\\", "/"))
	return u, nil
}

func (l *LocalFS) PresignGet(ctx context.Context, key string, _ time.Duration) (string, error) {
	return l.URL(ctx, key)
}

func (l *LocalFS) Delete(ctx context.Context, key string) error {
	return os.Remove(l.FullPath(key))
}
