package handlers

import (
	"fmt"
	"go-rest-chi/internal/auth"
	"go-rest-chi/internal/resp"
	"go-rest-chi/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MediaHandler struct {
	svc services.MediaService
}

func NewMediaHandler(svc services.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

func (h *MediaHandler) UploadPostMedia(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(200 << 20); err != nil {
		resp.Error(w, r, http.StatusBadRequest, "BAD_MULTIPART", "invalid form")
		return
	}

	postIDStr := chi.URLParam(r, "id")
	if postIDStr == "" {
		resp.Error(w, r, http.StatusBadRequest, "POST_ID_REQUIRED", "post_id is required")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		resp.Error(w, r, http.StatusBadRequest, "FILE_REQUIRED", "file is required")
		return
	}
	defer file.Close()

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil || postID <= 0 {
		resp.Error(w, r, http.StatusBadRequest, "BAD_POST_ID", "invalid post_id")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	size := header.Size

	userID := auth.UserIDFromCtx(r.Context())
	if userID == 0 {
		resp.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "missing user")
		return
	}

	pub, err := h.svc.SavePostMedia(r.Context(), userID, postID, header.Filename, file, size, mimeType)
	if err != nil {
		resp.Error(w, r, http.StatusBadRequest, "UPLOAD_FAIL", fmt.Sprintf("cannot save: %v", err))
		return
	}

	resp.OK(w, r, pub)
}
