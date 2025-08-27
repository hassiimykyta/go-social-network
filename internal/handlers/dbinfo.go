package handlers

import (
	"context"
	"net/http"
	"time"

	appdb "go-rest-chi/internal/db"
	"go-rest-chi/internal/resp"
)

type DBInfoHandler struct {
	DB *appdb.SQL
}

func NewDBInfoHandler(db *appdb.SQL) *DBInfoHandler {
	return &DBInfoHandler{DB: db}
}

func (h *DBInfoHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.Check(ctx); err != nil {
		resp.Error(w, r, http.StatusServiceUnavailable, "DB_CHECK_FAIL", "database check failed")
		return
	}
	resp.OK(w, r, map[string]bool{"ok": true})
}

func (h *DBInfoHandler) Version(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	v, err := h.DB.Version(ctx)
	if err != nil {
		resp.Error(w, r, http.StatusServiceUnavailable, "DB_VERSION_FAIL", "cannot read database version")
		return
	}
	resp.OK(w, r, map[string]string{"version": v})
}
