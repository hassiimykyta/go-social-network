package services

import (
	"context"
	"fmt"
	"go-rest-chi/internal/models"
	"go-rest-chi/internal/repositories"
)

type PostService interface {
	Create(ctx context.Context, title string, description string, userId int64) (models.PostPublic, error)
	GetAllByUser(ctx context.Context, userId int64) ([]models.PostPublic, error)
	SoftDelete(ctx context.Context, id int64) error
	UpdatePartitial(ctx context.Context, id int64, title *string, description *string) (models.PostPublic, error)
	ListPaginated(ctx context.Context, limit, offset int32) ([]models.PostPublic, error)
}

type postService struct {
	repo repositories.PostRepository
}

func NewPostService(r repositories.PostRepository) PostService {
	return &postService{repo: r}
}

// Create implements PostService.
func (p *postService) Create(ctx context.Context, title string, description string, userId int64) (models.PostPublic, error) {
	post, err := p.repo.Create(ctx, title, description, userId)
	if err != nil {
		return models.PostPublic{}, fmt.Errorf("CreatePost : %v", err)
	}

	return post.Public(), nil
}

// GetAllByUser implements PostService.
func (p *postService) GetAllByUser(ctx context.Context, userId int64) ([]models.PostPublic, error) {
	posts, err := p.repo.GetAllByUser(ctx, userId)
	if err != nil {
		return []models.PostPublic{}, fmt.Errorf("GetAllPostsByIser : %v", err)
	}

	out := make([]models.PostPublic, 0, len(posts))
	for _, pr := range posts {
		out = append(out, pr.Public())
	}

	return out, nil
}

// ListPaginated implements PostService.
func (p *postService) ListPaginated(ctx context.Context, limit int32, offset int32) ([]models.PostPublic, error) {
	posts, err := p.repo.ListPaginated(ctx, limit, offset)
	if err != nil {
		return []models.PostPublic{}, fmt.Errorf("GetListPaginated : %v", err)
	}

	out := make([]models.PostPublic, 0, len(posts))
	for _, pr := range posts {
		out = append(out, pr.Public())
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
