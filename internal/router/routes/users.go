package routes

import (
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountUsers(r chi.Router, posrSvc services.PostService) {
	h := handlers.NewPostHandler(posrSvc)

	r.Route("/users", func(rr chi.Router) {
		rr.Get("/{id}/posts", h.ListByUser)
	})
}
