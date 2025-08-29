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
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailTaken    = errors.New("email is already in use")
	ErrUsernameTaken = errors.New("username is already in use")
)

type UserRepository interface {
	Create(ctx context.Context, email string, username string, hash string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByUsername(ctx context.Context, username string) (models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

type userRepo struct {
	q *dbgen.Queries
}

func NewUserRepository(db *appdb.SQL) UserRepository {
	return &userRepo{q: db.Q}
}

func toModel(u dbgen.User) models.User {
	return models.User{
		Id:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}

func (r *userRepo) Create(ctx context.Context, email string, username string, hash string) (models.User, error) {
	u, err := r.q.CreateUser(ctx, dbgen.CreateUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
	})
	if err != nil {
		if helpers.IsUnique(err) {
			switch {
			case helpers.IsOnConstraint(err, "users_email"):
				return models.User{}, ErrEmailTaken
			case helpers.IsOnConstraint(err, "users_username"):
				return models.User{}, ErrUsernameTaken
			}
		}
		return models.User{}, fmt.Errorf("CreateUser: %v", err)
	}
	return toModel(u), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("GetUserByEmail : %v", err)

	}
	return toModel(u), nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (models.User, error) {
	u, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("GetUserByUsername : %v", err)
	}
	return toModel(u), nil
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.q.ExistsUserByEmail(ctx, email)
}

func (r *userRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.q.ExistsUserByUsername(ctx, username)
}
