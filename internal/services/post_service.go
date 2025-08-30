package services

import (
	"context"
	"fmt"
	"go-rest-chi/internal/models"
	"go-rest-chi/internal/repositories"
	"go-rest-chi/internal/storage"
)

type PostService interface {
	Create(ctx context.Context, title string, description string, userId int64) (models.PostPublic, error)
	SoftDelete(ctx context.Context, id int64) error
	UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.PostPublic, error)
	ListPaginated(ctx context.Context, userId *int64, limit, offset int32) ([]models.PostMediaPublic, error)
}

type postService struct {
	repo repositories.PostRepository
	st   storage.Storage
}

func NewPostService(r repositories.PostRepository, st storage.Storage) PostService {
	return &postService{repo: r, st: st}
}

// Create implements PostService.
func (p *postService) Create(ctx context.Context, title string, description string, userId int64) (models.PostPublic, error) {
	post, err := p.repo.Create(ctx, title, description, userId)
	if err != nil {
		return models.PostPublic{}, fmt.Errorf("CreatePost : %v", err)
	}

	return post.Public(), nil
}

// ListPaginated implements PostService.
func (p *postService) ListPaginated(ctx context.Context, userId *int64, limit int32, offset int32) ([]models.PostMediaPublic, error) {
	items, err := p.repo.ListWithMediaPaginated(ctx, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ListWithMediaPaginated: %w", err)
	}

	out := make([]models.PostMediaPublic, 0, len(items))
	for _, it := range items {
		pub := models.PostMediaPublic{
			Post:   it.Post.Public(),
			Medias: make([]models.MediaPublic, 0, len(it.Medias)),
		}

		for _, m := range it.Medias {
			url, _ := p.st.URL(ctx, m.StorageKey)
			pub.Medias = append(pub.Medias, models.MediaPublic{
				ID:         m.ID,
				Kind:       m.Kind,
				MimeType:   m.MimeType,
				URL:        url,
				Width:      m.Width,
				Height:     m.Height,
				DurationMs: m.DurationMs,
			})
		}

		out = append(out, pub)
	}
	return out, nil
}

// SoftDelete implements PostService.
func (p *postService) SoftDelete(ctx context.Context, id int64) error {
	return p.repo.SoftDelete(ctx, id)
}

// UpdatePartitial implements PostService.
func (p *postService) UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.PostPublic, error) {
	post, err := p.repo.UpdatePartitial(ctx, id, title, description)
	if err != nil {
		return models.PostPublic{}, fmt.Errorf("UpdatePost[%d] : %v", id, err)
	}

	return post.Public(), nil
}
