package routes

import (
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountPosts(r chi.Router, posrSvc services.PostService) {
	h := handlers.NewPostHandler(posrSvc)

	r.Route("/posts", func(rr chi.Router) {
		rr.Get("/", h.List)
		rr.Post("/", h.Create)
		rr.Patch("/{id}", h.UpdatePartial)
		rr.Delete("/{id}", h.SoftDelete)
	})
}
