package router

import (
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/services"
)

type Services struct {
	Users services.UserService
}

type Deps struct {
	DB       *appdb.SQL
	Services Services
}
