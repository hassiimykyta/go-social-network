package handlers

import (
	"encoding/json"
	"go-rest-chi/internal/auth"
	"go-rest-chi/internal/helpers"
	"go-rest-chi/internal/resp"
	"go-rest-chi/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	svc services.PostService
}

func NewPostHandler(s services.PostService) *PostHandler { return &PostHandler{svc: s} }

type createPostReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updatePostReq struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "invalid json")
		return
	}

	userID := auth.UserIDFromCtx(r.Context())
	if userID == 0 {
		resp.Error(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "missing user")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		resp.Error(w, r, http.StatusBadRequest, "MISSING_FIELDS", "title is required")
		return
	}

	post, err := h.svc.Create(r.Context(), req.Title, req.Description, userID)
	if err != nil {
		resp.Error(w, r, http.StatusInternalServerError, "CREATE_POST_FAIL", "cannot create post")
		return
	}
	resp.OK(w, r, post)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := helpers.ParseInt(r.URL.Query().Get("limit"), 10, 100)
	offset := helpers.ParseInt(r.URL.Query().Get("offset"), 0, 1_000_000)

	items, err := h.svc.ListPaginated(r.Context(), int32(limit), int32(offset))
	if err != nil {
		resp.Error(w, r, http.StatusInternalServerError, "LIST_POSTS_FAIL", "cannot list posts")
		return
	}
	resp.OK(w, r, map[string]any{
		"items": items,
		"page":  map[string]any{"limit": limit, "offset": offset},
	})
}

func (h *PostHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	uidStr := chi.URLParam(r, "id")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil || uid <= 0 {
		resp.Error(w, r, http.StatusBadRequest, "BAD_USER_ID", "invalid user id")
		return
	}

	items, err := h.svc.GetAllByUser(r.Context(), uid)
	if err != nil {
		resp.Error(w, r, http.StatusInternalServerError, "LIST_USER_POSTS_FAIL", "cannot list user posts")
		return
	}
	resp.OK(w, r, map[string]any{
		"items": items,
	})
}

func (h *PostHandler) UpdatePartial(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		resp.Error(w, r, http.StatusBadRequest, "BAD_POST_ID", "invalid post id")
		return
	}

	var req updatePostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "invalid json")
		return
	}

	if req.Title == nil && req.Description == nil {
		resp.Error(w, r, http.StatusBadRequest, "EMPTY_PATCH", "no fields to update")
		return
	}

	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			resp.Error(w, r, http.StatusBadRequest, "TITLE_EMPTY", "title cannot be empty")
			return
		}
		req.Title = &t
	}
	if req.Description != nil {
		d := strings.TrimSpace(*req.Description)
		req.Description = &d
	}

	post, err := h.svc.UpdatePartitial(r.Context(), id, req.Title, req.Description)
	if err != nil {
		resp.Error(w, r, http.StatusInternalServerError, "UPDATE_POST_FAIL", "cannot update post")
		return
	}
	resp.OK(w, r, post)
}

func (h *PostHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		resp.Error(w, r, http.StatusBadRequest, "BAD_POST_ID", "invalid post id")
		return
	}
	if err := h.svc.SoftDelete(r.Context(), id); err != nil {
		resp.Error(w, r, http.StatusInternalServerError, "DELETE_POST_FAIL", "cannot delete post")
		return
	}
	resp.OK(w, r, map[string]bool{"deleted": true})
}
