package handlers

import (
	"encoding/json"
	"go-rest-chi/internal/resp"
	"go-rest-chi/internal/services"
	"net/http"
	"strings"
	"unicode/utf8"
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

func lookLikeEmail(s string) bool {
	return strings.Count(s, "@") == 1 && strings.Contains(s, ".")
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

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if !lookLikeEmail(email) {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "invalid json")

	}

	password := req.Password
	if utf8.RuneCountInString(password) < 8 {
		resp.Error(w, r, http.StatusBadRequest, "BAD_JSON", "pasword must contain 8 or more symbols")

	}

	username := strings.TrimSpace(req.Username)

	pub, err := h.users.Register(r.Context(), email, username, password)
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

	identifier := strings.TrimSpace(req.Identifier)
	password := strings.TrimSpace(req.Password)

	if identifier == "" && password == "" {
		resp.Error(w, r, http.StatusBadRequest, "MISSING_FIELDS", "identifier and password are required")
		return

	}

	pub, err := h.users.Login(r.Context(), identifier, password)
	if err != nil {
		resp.Error(w, r, http.StatusConflict, "REGISTER_FAILED", err.Error())
		return
	}

	resp.OK(w, r, pub)

}
