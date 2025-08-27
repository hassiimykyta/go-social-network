package routes

import (
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func MountDB(r chi.Router, db *appdb.SQL) {
	h := handlers.NewDBInfoHandler(db)
	r.Get("/db/check", h.Check)
	r.Get("/db/version", h.Version)
}
