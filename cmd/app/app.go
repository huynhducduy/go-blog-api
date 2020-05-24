package app

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Run() error {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return http.ListenAndServe(":3000", r)
}