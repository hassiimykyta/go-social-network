package services

import (
	"context"
	"fmt"
	"go-rest-chi/internal/helpers"
	"go-rest-chi/internal/models"
	"go-rest-chi/internal/repositories"
	"go-rest-chi/internal/storage"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrNotImplemented  = fmt.Errorf("not implemented")
	ErrUnsupportedMime = fmt.Errorf("unsupported content type")
)

type MediaService interface {
	SavePostMedia(ctx context.Context, userID, postID int64, filename string, r io.Reader, size int64, mimeType string) (models.MediaPublic, error)

	UploadUserAvatar(ctx context.Context, userID int64, filename string, r io.Reader, size int64, mimeType string) (models.MediaPublic, error)
}

type mediaService struct {
	repo repositories.MediaRepository
	st   storage.Storage
	ttl  time.Duration
}

func NewMediaService(repo repositories.MediaRepository, st storage.Storage, presignTTL time.Duration) MediaService {
	return &mediaService{repo: repo, st: st, ttl: presignTTL}
}

// SavePostImage implements MediaService.
func (med *mediaService) SavePostMedia(ctx context.Context, userID int64, postID int64, filename string, r io.Reader, size int64, mimeType string) (models.MediaPublic, error) {
	kind := helpers.InferKind(mimeType)

	if kind == "" {
		return models.MediaPublic{}, ErrUnsupportedMime
	}

	filename = helpers.SanitizeFilename(filename)

	if ext := strings.ToLower(filepath.Ext(filename)); ext == "" {
		if exts, _ := mime.ExtensionsByType(mimeType); len(exts) > 0 {
			filename += exts[0]
		}
	}

	key := storage.BuildPostKey(userID, filename, time.Now())

	if err := med.st.Save(ctx, key, r, size, mimeType); err != nil {
		return models.MediaPublic{}, fmt.Errorf("storage save: %w", err)
	}

	media, err := med.repo.Create(ctx, repositories.CreateMediaParams{
		OwnerID:    userID,
		Kind:       kind,
		StorageKey: key,
		MimeType:   mimeType,
		SizeBytes:  size,
	})

	if err != nil {
		return models.MediaPublic{}, fmt.Errorf("media create: %w", err)
	}

	if err := med.repo.AttachToPost(ctx, postID, media.ID, 0); err != nil {
		return models.MediaPublic{}, fmt.Errorf("attach to post: %w", err)

	}

	url, _ := med.st.PresignGet(ctx, media.StorageKey, med.ttl)
	return media.PublicWithURL(url), nil

}

// UploadUserAvatar implements MediaService.
func (med *mediaService) UploadUserAvatar(ctx context.Context, userID int64, filename string, r io.Reader, size int64, mimeType string) (models.MediaPublic, error) {
	panic("unimplemented")
}
