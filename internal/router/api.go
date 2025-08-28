package router

import (
	"go-rest-chi/internal/router/routes"

	"github.com/go-chi/chi/v5"
)

func MountAPI(r *chi.Mux, d Deps) {
	r.Route("/api", func(api chi.Router) {
		api.Route("/v1", func(v1 chi.Router) {
			routes.MountHealth(v1, d.DB)
			routes.MountDB(v1, d.DB)
			routes.MountAuth(v1, d.Services.Users)
		})
	})
}
