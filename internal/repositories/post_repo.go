package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/dbgen"
	"go-rest-chi/internal/models"
)

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrPostNotUpdate = errors.New("unable update post")
)

type PostRepository interface {
	Create(ctx context.Context, title string, description string, userId int64) (models.Post, error)
	GetAllByUser(ctx context.Context, userId int64) ([]models.Post, error)
	SoftDelete(ctx context.Context, id int64) error
	UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.Post, error)
	ListPaginated(ctx context.Context, limit, offset int32) ([]models.Post, error)
}

type postRepo struct {
	q *dbgen.Queries
}

func toPostModel(p dbgen.Post) models.Post {
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

func NewPostRepository(db *appdb.SQL) PostRepository {
	return &postRepo{q: db.Q}
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

	return toPostModel(row), nil
}

// GetAllByUser implements PostRepository.
func (p *postRepo) GetAllByUser(ctx context.Context, userId int64) ([]models.Post, error) {
	rows, err := p.q.GetAllPostsByUser(ctx, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Post{}, ErrPostNotFound
		}
		return []models.Post{}, fmt.Errorf("GetAllPostByUser : %v", err)

	}
	out := make([]models.Post, 0, len(rows))
	for _, pr := range rows {
		out = append(out, toPostModel(pr))
	}
	return out, nil
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
		out = append(out, toPostModel(pr))
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

	return toPostModel(row), nil
}
