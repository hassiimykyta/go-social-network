package routes

import (
	"go-rest-chi/internal/auth"
	"go-rest-chi/internal/handlers"
	"go-rest-chi/internal/services"

	"github.com/go-chi/chi/v5"
)

func MountMedia(r chi.Router, jwt *auth.Service, mcv services.MediaService) {
	h := handlers.NewMediaHandler(mcv)

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(jwt))
		r.Route("/media", func(r chi.Router) {
			r.Post("/posts/{id}", h.UploadPostMedia)
		})
	})
}
