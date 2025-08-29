package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Save(ctx context.Context, key string, r io.Reader, size int64, mime string) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	URL(ctx context.Context, key string) (string, error)
	PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}
