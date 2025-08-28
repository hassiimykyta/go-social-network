package services

import (
	"context"
	"errors"
	"go-rest-chi/internal/models"
	"go-rest-chi/internal/repositories"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrBadCredentials  = errors.New("invalid credentials")
)

type UserService interface {
	Register(ctx context.Context, email string, username string, password string) (models.UserPublic, error)
	Login(ctx context.Context, identifier, password string) (models.UserPublic, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{repo: r}
}

func lookLikeEmail(s string) bool {
	return strings.Count(s, "@") == 1 && strings.Contains(s, ".")
}

// Register implements UserService.
func (u *userService) Register(ctx context.Context, email string, username string, password string) (models.UserPublic, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if !lookLikeEmail(email) {
		return models.UserPublic{}, ErrInvalidEmail
	}

	if utf8.RuneCountInString(password) < 8 {
		return models.UserPublic{}, ErrInvalidPassword
	}

	username = strings.TrimSpace(username)

	if ok, err := u.repo.ExistsByUsername(ctx, username); err == nil && ok {
		return models.UserPublic{}, repositories.ErrUsernameTaken
	}

	if ok, err := u.repo.ExistsByEmail(ctx, email); err == nil && ok {
		return models.UserPublic{}, repositories.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.UserPublic{}, err
	}

	usr, err := u.repo.Create(ctx, email, username, string(hash))
	if err != nil {
		return models.UserPublic{}, err
	}

	return usr.Public(), nil
}

// Login implements UserService.
func (u *userService) Login(ctx context.Context, identifier, password string) (models.UserPublic, error) {
	id := strings.TrimSpace(identifier)
	if id == "" && strings.TrimSpace(password) == "" {
		return models.UserPublic{}, ErrBadCredentials
	}

	var usr models.User
	var err error

	if strings.Contains(id, "@") {
		usr, err = u.repo.GetByEmail(ctx, id)
	} else {
		usr, err = u.repo.GetByUsername(ctx, id)
	}

	if err != nil {
		return models.UserPublic{}, ErrBadCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password)) != nil {
		return models.UserPublic{}, ErrInvalidPassword
	}

	return usr.Public(), nil
}
