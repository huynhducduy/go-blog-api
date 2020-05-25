package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"go-blog/internal/blog"
	"go-blog/internal/config"
	"go-blog/internal/db"
	"net/http"
)

func Run() error {

	config.ReadConfig()

	db.OpenConnection()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Route("/blog", func(r chi.Router) {
				r.Get("/", blog.RouterList)
				//r.Post("/", blog.RouterCreate)
				//
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", blog.RouterRead)
				//	r.Put("/", blog.RouterUpdate)
				//	r.Delete("/", blog.RouterDelete)
				})
			})
		})
	})

	log.Printf("Running at port 80")
	return http.ListenAndServe(":80", r)
}