package app

import (
	"github.com/go-chi/chi"
	"go-blog/internal/config"
	"go-blog/internal/db"
	"net/http"
)

func Run() error {

	config.ReadConfig()

	db.OpenConnection()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})
	return http.ListenAndServe(":3000", r)
}