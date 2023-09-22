package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Tus1688/go-blob-storage/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

var authorizationHeader string

func main() {
	authorizationHeader = os.Getenv("AUTHORIZATION_HEADER")
	if authorizationHeader == "" {
		log.Fatal("AUTHORIZATION_HEADER not found")
	}

	log.Print("server running on port 3000")
	r := initRouter()

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func initRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(EnforceAuthentication)

	r.Post(
		"/file", func(w http.ResponseWriter, r *http.Request) {
			file, m, err := r.FormFile("file")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()

			filename := uuid.New().String() + filepath.Ext(m.Filename)

			dst, err := os.Create("/usr/share/nginx/html/" + filename)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = render.JSON(w, http.StatusOK, map[string]string{"filename": filename})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
	)
	r.Delete(
		"/file/{filename}", func(w http.ResponseWriter, r *http.Request) {
			filename := chi.URLParam(r, "filename")

			err := os.Remove("/usr/share/nginx/html/" + filename)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			err = render.JSON(w, http.StatusOK, map[string]string{"filename": filename})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		},
	)

	return r
}

func EnforceAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader != authorizationHeader {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
