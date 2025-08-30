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
	ListPaginated(ctx context.Context, limit, offset int32) ([]models.Post, error)
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

	var width *int32
	if r.MediaWidth.Valid {
		v := r.MediaWidth.Int32
		width = &v
	}

	var height *int32
	if r.MediaHeight.Valid {
		v := r.MediaHeight.Int32
		height = &v
	}

	var duration *int32
	if r.MediaDurationMs.Valid {
		v := r.MediaDurationMs.Int32
		duration = &v
	}

	m := models.Media{
		ID:         r.MediaID.Int64,
		OwnerID:    r.UserID,
		Kind:       helpers.DerefStr(r.MediaKind),
		MimeType:   helpers.DerefStr(r.MediaMimeType),
		StorageKey: helpers.DerefStr(r.MediaStorageKey),
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
	var uid sql.NullInt64
	if userId != nil {
		uid = sql.NullInt64{Int64: *userId, Valid: true}
	}
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

// ListPaginated implements PostRepository.
func (p *postRepo) ListPaginated(ctx context.Context, limit int32, offset int32) ([]models.Post, error) {
	rows, err := p.q.ListPostsPaginated(ctx, dbgen.ListPostsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Post{}, ErrPostNotFound
		}
		return []models.Post{}, fmt.Errorf("GetAllPostPaginated : %v", err)

	}
	out := make([]models.Post, 0, len(rows))
	for _, pr := range rows {
		out = append(out, toPostModelRow(pr))
	}
	return out, nil
}

// SoftDelete implements PostRepository.
func (p *postRepo) SoftDelete(ctx context.Context, id int64) error {
	return p.q.SoftDeletePost(ctx, id)
}

// UpdatePartitial implements PostRepository.
func (p *postRepo) UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.Post, error) {
	var tns, dns sql.NullString
	if title != nil {
		tns = sql.NullString{String: *title, Valid: true}
	}
	if description != nil {
		dns = sql.NullString{String: *description, Valid: true}
	}

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
