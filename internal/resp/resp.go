package resp

import (
	"net/http"

	"github.com/go-chi/render"
)

type Success struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

type Err struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, data)
}

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Success{Status: "ok", Data: data})
}

func Error(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	render.Status(r, status)
	render.JSON(w, r, Err{Status: "error", Code: code, Message: msg})
}
