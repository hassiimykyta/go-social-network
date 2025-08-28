package routes

import (
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountAuth(r chi.Router, usersSvc services.UserService) {
	h := handlers.NewAuthHandler(usersSvc)

	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
}
