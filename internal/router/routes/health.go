package routes

import (
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func MountHealth(r chi.Router, db *appdb.SQL) {
	h := handlers.NewHealthHandler(db)
	r.Get("/live", h.Liveness)
	r.Get("/ready", h.Readiness)
}
