package routes

import (
	"go-rest-chi/internal/auth"
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountPosts(r chi.Router, jwtSvc *auth.Service, posrSvc services.PostService) {
	h := handlers.NewPostHandler(posrSvc)

	r.Route("/posts", func(rr chi.Router) {
		rr.Get("/", h.List)
		rr.Group(func(priv chi.Router) {
			priv.Use(auth.Middleware(jwtSvc))
			priv.Post("/", h.Create)
			priv.Patch("/{id}", h.UpdatePartial)
			priv.Delete("/{id}", h.SoftDelete)
		})

	})
}
