package app

import (
	"github.com/go-chi/chi"
	"go-blog/internal/config"
	"net/http"
)

func Run() error {

	config.ReadConfig()


	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("welcome"))
	})
	return http.ListenAndServe(":3000", r)
}