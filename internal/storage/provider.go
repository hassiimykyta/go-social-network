package storage

import (
	"fmt"
	"go-rest-chi/internal/config"
)

func NewFromConfig(cfg config.StorageConfig) (Storage, error) {
	switch cfg.Driver {
	case "local":
		return NewLocalFS(cfg.LocalDir, cfg.PublicBaseURL), nil
	case "s3":
		return nil, fmt.Errorf("s3 driver not implemented yet")
	default:
		return nil, fmt.Errorf("unknown storage driver %v", cfg.Driver)
	}

}
