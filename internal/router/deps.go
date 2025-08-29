package router

import (
	"go-rest-chi/internal/auth"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/services"
)

type Services struct {
	Users services.UserService
	Posts services.PostService
	JWT   *auth.Service
}

type Deps struct {
	DB       *appdb.SQL
	Services Services
}
