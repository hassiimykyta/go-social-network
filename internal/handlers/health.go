package handlers

import (
	"context"
	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/resp"
	"net/http"
	"time"
)

type HealthHandler struct{ DB *appdb.SQL }

func NewHealthHandler(db *appdb.SQL) *HealthHandler { return &HealthHandler{DB: db} }

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	resp.OK(w, r, map[string]bool{"alive": true})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.Check(ctx); err != nil {
		resp.Error(w, r, http.StatusServiceUnavailable, "DB_NOT_READY", "database is not ready")
		return
	}

	resp.OK(w, r, map[string]bool{"ready": true})
}
