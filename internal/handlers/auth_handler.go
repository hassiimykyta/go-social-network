package handlers

import (
	"encoding/json"
	"go-rest-chi/internal/resp"
	"go-rest-chi/internal/services"
	"net/http"
	"strings"
)

type AuthHandler struct {
	users services.UserService
}

func NewAuthHandler(s services.UserService) *AuthHandler {
	return &AuthHandler{
		users: s,
	}
}

type registerReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginReq struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "invalid json")
		return
	}

	if strings.TrimSpace(req.Email) == "" && strings.TrimSpace(req.Username) == "" && strings.TrimSpace(req.Password) == "" {
		resp.Error(w, r, http.StatusBadRequest, "MISSING_FIELDS", "email, username and password are required")
		return

	}

	pub, err := h.users.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		resp.Error(w, r, http.StatusConflict, "REGISTER_FAILED", err.Error())
		return
	}

	resp.OK(w, r, pub)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "invalid json")
		return
	}

	if strings.TrimSpace(req.Identifier) == "" && strings.TrimSpace(req.Password) == "" {
		resp.Error(w, r, http.StatusBadRequest, "MISSING_FIELDS", "identifier and password are required")
		return

	}

	pub, err := h.users.Login(r.Context(), req.Identifier, req.Password)
	if err != nil {
		resp.Error(w, r, http.StatusConflict, "REGISTER_FAILED", err.Error())
		return
	}

	resp.OK(w, r, pub)

}
