package routes

import (
	"go-rest-chi/internal/auth"
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountAuth(r chi.Router, jwtSvc *auth.Service, usersSvc services.UserService) {
	h := handlers.NewAuthHandler(jwtSvc, usersSvc)

	r.Route("/auth", func(rr chi.Router) {
		rr.Post("/register", h.Register)
		rr.Post("/login", h.Login)
		rr.Post("/token/refresh", h.Refresh)
	})
}
