package repositories

import (
	"context"
	"database/sql"
	"fmt"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/dbgen"
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
	var deleted *sql.NullTime
	_ = deleted

	out := models.Media{
		ID:         m.ID,
		OwnerID:    m.OwnerID,
		Kind:       m.Kind,
		StorageKey: m.StorageKey,
		MimeType:   m.MimeType,
		SizeBytes:  m.SizeBytes,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  m.DeletedAt,
	}

	if m.Width.Valid {
		v := m.Width.Int32
		out.Width = &v
	}
	if m.Height.Valid {
		v := m.Height.Int32
		out.Height = &v
	}
	if m.DurationMs.Valid {
		v := m.DurationMs.Int32
		out.DurationMs = &v
	}
	return out
}

func p2nulli32(p *int32) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *p, Valid: true}
}

// Create implements MediaRepository.
func (m *mediaRepo) Create(ctx context.Context, p CreateMediaParams) (models.Media, error) {
	row, err := m.q.CreateMedia(ctx, dbgen.CreateMediaParams{
		OwnerID:    p.OwnerID,
		Kind:       p.Kind,
		StorageKey: p.StorageKey,
		MimeType:   p.MimeType,
		SizeBytes:  p.SizeBytes,
		Width:      p2nulli32(p.Width),
		Height:     p2nulli32(p.Height),
		DurationMs: p2nulli32(p.DurationMs),
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
