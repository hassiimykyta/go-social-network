package repositories

import (
	"context"
	"database/sql"
	"fmt"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/dbgen"
	"go-rest-chi/internal/helpers"
	"go-rest-chi/internal/models"
)

type CreateMediaParams struct {
	OwnerID    int64
	Kind       string
	StorageKey string
	MimeType   string
	SizeBytes  int64
	Width      *int32
	Height     *int32
	DurationMs *int32
}

type MediaRepository interface {
	Create(ctx context.Context, p CreateMediaParams) (models.Media, error)
	AttachToPost(ctx context.Context, postID, mediaID int64, position int) error
}

type mediaRepo struct {
	q *dbgen.Queries
}

func NewMediaRepository(sql *appdb.SQL) MediaRepository {
	return &mediaRepo{q: sql.Q}

}

func toMediaModel(m dbgen.Medium) models.Media {
	width := helpers.PtrFromNull(m.Width.Valid, m.Width.Int32)
	height := helpers.PtrFromNull(m.Height.Valid, m.Height.Int32)
	duration := helpers.PtrFromNull(m.DurationMs.Valid, m.DurationMs.Int32)

	out := models.Media{
		ID:         m.ID,
		OwnerID:    m.OwnerID,
		Kind:       m.Kind,
		StorageKey: m.StorageKey,
		MimeType:   m.MimeType,
		Width:      width,
		Height:     height,
		DurationMs: duration,
		SizeBytes:  m.SizeBytes,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  m.DeletedAt,
	}

	return out
}

// Create implements MediaRepository.
func (m *mediaRepo) Create(ctx context.Context, p CreateMediaParams) (models.Media, error) {

	width := helpers.ToNull(p.Width, func(v int32) sql.NullInt32 {
		return sql.NullInt32{Int32: v, Valid: true}
	})

	height := helpers.ToNull(p.Height, func(v int32) sql.NullInt32 {
		return sql.NullInt32{Int32: v, Valid: true}
	})

	duration := helpers.ToNull(p.DurationMs, func(v int32) sql.NullInt32 {
		return sql.NullInt32{Int32: v, Valid: true}
	})

	row, err := m.q.CreateMedia(ctx, dbgen.CreateMediaParams{
		OwnerID:    p.OwnerID,
		Kind:       p.Kind,
		StorageKey: p.StorageKey,
		MimeType:   p.MimeType,
		SizeBytes:  p.SizeBytes,
		Width:      width,
		Height:     height,
		DurationMs: duration,
	})
	if err != nil {
		return models.Media{}, fmt.Errorf("CreateMedia: %w", err)
	}
	return toMediaModel(row), nil
}

// AttachToPost implements MediaRepository.
func (m *mediaRepo) AttachToPost(ctx context.Context, postID int64, mediaID int64, position int) error {
	return m.q.AttachMediaToPost(ctx, dbgen.AttachMediaToPostParams{
		PostID:   postID,
		MediaID:  mediaID,
		Position: int32(position),
	})
}
