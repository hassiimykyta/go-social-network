package router

import (
	"go-rest-chi/internal/router/routes"

	"github.com/go-chi/chi/v5"
)

func MountAPI(r *chi.Mux, d Deps) {
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			routes.MountAuth(v1, d.Services.JWT, d.Services.Users)
			routes.MountPosts(v1, d.Services.JWT, d.Services.Posts)
			routes.MountMedia(v1, d.Services.JWT, d.Services.Media)
		})
	})
}
