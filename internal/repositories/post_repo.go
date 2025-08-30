package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/dbgen"
	"go-rest-chi/internal/helpers"
	"go-rest-chi/internal/models"
)

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrPostNotUpdate = errors.New("unable update post")
)

type PostRepository interface {
	Create(ctx context.Context, title string, description string, userId int64) (models.Post, error)
	SoftDelete(ctx context.Context, id int64) error
	UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.Post, error)
	ListWithMediaPaginated(ctx context.Context, userId *int64, limit, offset int32) ([]models.PostMedia, error)
}

type postRepo struct {
	q *dbgen.Queries
}

func toPostModelRow(p dbgen.Post) models.Post {
	return models.Post{
		Id:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		UserId:      p.UserID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}
}

func toMediaModelFromRow(r dbgen.ListPostsWithMediaPaginatedRow) (models.Media, bool) {
	if !r.MediaID.Valid {
		return models.Media{}, false
	}

	width := helpers.PtrFromNull(r.MediaWidth.Valid, r.MediaWidth.Int32)
	height := helpers.PtrFromNull(r.MediaHeight.Valid, r.MediaWidth.Int32)
	duration := helpers.PtrFromNull(r.MediaDurationMs.Valid, r.MediaDurationMs.Int32)
	kind := helpers.ValueOr(r.MediaKind.Valid, r.MediaKind.String, "")
	mimeType := helpers.ValueOr(r.MediaMimeType.Valid, r.MediaMimeType.String, "")
	storageKey := helpers.ValueOr(r.MediaStorageKey.Valid, r.MediaStorageKey.String, "")

	m := models.Media{
		ID:         r.MediaID.Int64,
		OwnerID:    r.UserID,
		Kind:       kind,
		MimeType:   mimeType,
		StorageKey: storageKey,
		Width:      width,
		Height:     height,
		DurationMs: duration,
	}
	return m, true
}

func NewPostRepository(db *appdb.SQL) PostRepository {
	return &postRepo{q: db.Q}
}

// ListWithMedia implements PostRepository.
func (p *postRepo) ListWithMediaPaginated(ctx context.Context, userId *int64, limit int32, offset int32) ([]models.PostMedia, error) {

	uid := helpers.ToNull(userId, func(v int64) sql.NullInt64 {
		return sql.NullInt64{Int64: v, Valid: true}
	})

	rows, err := p.q.ListPostsWithMediaPaginated(ctx, dbgen.ListPostsWithMediaPaginatedParams{
		UserID: uid,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostMedia{}, ErrPostNotFound
		}
		return []models.PostMedia{}, fmt.Errorf("GetAllPostMediaPaginated : %v", err)
	}

	byID := make(map[int64]*models.PostMedia, len(rows))
	order := make([]int64, 0, len(rows))

	for _, pm := range rows {
		item, ok := byID[pm.ID]
		if !ok {
			item = &models.PostMedia{
				Post: models.Post{
					Id:          pm.ID,
					Title:       pm.Title,
					Description: pm.Description,
					UserId:      pm.UserID,
					CreatedAt:   pm.CreatedAt,
					UpdatedAt:   pm.UpdatedAt,
					DeletedAt:   pm.DeletedAt,
				},
				Medias: make([]models.Media, 0, 2),
			}
			byID[pm.ID] = item
			order = append(order, pm.ID)
		}

		if media, ok := toMediaModelFromRow(pm); ok {
			item.Medias = append(item.Medias, media)
		}

	}

	out := make([]models.PostMedia, 0, len(byID))
	for _, id := range order {
		out = append(out, *byID[id])
	}
	return out, nil
}

// Create implements PostRepository.
func (p *postRepo) Create(ctx context.Context, title string, description string, userId int64) (models.Post, error) {
	row, err := p.q.CreatePost(ctx, dbgen.CreatePostParams{
		Title:       title,
		Description: description,
		UserID:      userId,
	})
	if err != nil {
		return models.Post{}, fmt.Errorf("CreatePost: %v", err)
	}

	return toPostModelRow(row), nil
}

// SoftDelete implements PostRepository.
func (p *postRepo) SoftDelete(ctx context.Context, id int64) error {
	return p.q.SoftDeletePost(ctx, id)
}

// UpdatePartitial implements PostRepository.
func (p *postRepo) UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.Post, error) {

	tns := helpers.ToNull(title, func(v string) sql.NullString {
		return sql.NullString{String: v, Valid: true}
	})

	dns := helpers.ToNull(description, func(v string) sql.NullString {
		return sql.NullString{String: v, Valid: true}
	})

	row, err := p.q.UpdatePostPartial(ctx, dbgen.UpdatePostPartialParams{
		ID:          id,
		Title:       tns,
		Description: dns,
	})
	if err != nil {
		return models.Post{}, ErrPostNotUpdate
	}

	return toPostModelRow(row), nil
}
