package app

import (
	"go-blog/internal/auth"
	"go-blog/internal/blog"
	//"go-blog/internal/tag"
	"go-blog/internal/user"
	//"go-blog/internal/blog/reply"
	"go-blog/internal/config"
	"go-blog/internal/db"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
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
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", auth.Login)
			r.Post("/register", auth.Register)

			r.Group(func(r chi.Router) {
				r.Use(auth.AuthenticationMiddleware)
				r.Get("/profile", auth.Profile)
				r.Post("/profile", auth.UpdateProfile)
				r.Post("/profile/password", auth.ChangeProfilePassword)
			})
		})

		r.Route("/blog", func(r chi.Router) {
			r.Get("/", blog.RouterList)
			r.With(auth.AuthenticationMiddleware).Post("/", blog.RouterCreate)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", blog.RouterRead)

				r.With(auth.AuthenticationMiddleware).Put("/", blog.RouterUpdate)
				r.With(auth.AuthenticationMiddleware).Delete("/", blog.RouterDelete)

				//r.Route("/reply", func(r chi.Router) {
				//	r.Get("/", reply.RouterList)
				//  r.With(auth.AuthenticationMiddleware).Post("/", reply.RouterCreate)
				//	r.Route("/{id}", func(r chi.Router) {
				//      r.Use(auth.AuthenticationMiddleware)
				//		r.Put("/", reply.RouterUpdate)
				//		r.Delete("/", reply.RouterDelete)
				//	})
				//})
			})
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", user.RouterList)

			r.Route("/{id}", func(r chi.Router) {
				//r.Get("/", user.RouterRead)
				//r.Get("/blog", user.RouterListBlog)
			})
		})

		//r.Route("/tag", func(r chi.Router) {
		//	r.Get("/", tag.RouterList)
		//	r.With(auth.AuthenticationMiddleware).Post("/", tag.RouterCreate)
		//
		//	r.Route("/{tag}", func(r chi.Router) {
		//		r.Get("/", tag.RouterRead)
		//		r.Group(func(r chi.Router) {
		//			r.Use(auth.AuthenticationMiddleware)
		//			r.Put("/", tag.RouterUpdate)
		//			r.Delete("/", tag.RouterDelete)
		//			r.Get("/blog", tag.RouterListBlog)
		//		})
		//
		//	})
		//})
	})

	log.Printf("Running at port 80")
	return http.ListenAndServe(":80", r)
}
