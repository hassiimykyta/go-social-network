package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func MountLocalMediaStatic(r chi.Router, localDir string) {
	fs := http.StripPrefix("/media/", http.FileServer(http.Dir(localDir)))

	r.Get("/media/*", func(w http.ResponseWriter, r *http.Request) {
		p := chi.URLParam(r, "*")
		clean := filepath.Clean(p)

		if clean == "." || clean == "/" || clean == ".." || clean == "" {
			http.NotFound(w, r)
			return
		}

		if len(clean) > 0 && clean[0] == '.' {
			http.NotFound(w, r)
			return
		}

		if _, err := os.Stat(filepath.Join(localDir, clean)); err != nil {
			http.NotFound(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	})
}
