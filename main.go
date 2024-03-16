package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"gitlab.com/arnaud-web/neli-webservices/api/auth/bearer"
	"gitlab.com/arnaud-web/neli-webservices/api/auth/credentials"
	"gitlab.com/arnaud-web/neli-webservices/api/check"
	"gitlab.com/arnaud-web/neli-webservices/api/content"
	"gitlab.com/arnaud-web/neli-webservices/api/user"
	"gitlab.com/arnaud-web/neli-webservices/api/user/admin"
	"gitlab.com/arnaud-web/neli-webservices/api/user/leader"
	"gitlab.com/arnaud-web/neli-webservices/api/user/member"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/api/capture"
	)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

func main() {
	tokenAuth := jwtauth.New(bearer.TokenAlgorithm, []byte(*config.TokenSecret), nil)
	refreshAuth := jwtauth.New(bearer.TokenAlgorithm, []byte(*config.RefreshSecret), nil)

	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	a := chi.NewRouter()

	a.Use(cors.Handler)

	a.Route("/", func(a chi.Router) {
		a.Post("/login", credentials.Login)
		a.Post("/set-password", credentials.Set)
		a.Post("/reset-password", credentials.Reset)

		a.Route("/refresh", func(a chi.Router) {
			a.Use(jwtauth.Verifier(refreshAuth))
			a.Use(check.Authenticator)
			a.Post("/", bearer.Refresh)
		})
	})

	r := chi.NewRouter()

	r.Use(cors.Handler)

	r.Use(middleware.Logger)

	r.Mount("/auth", a)

	fs := http.StripPrefix("/assets", http.FileServer(http.Dir(*config.AssetsPath)))
	r.Get("/assets*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path  == "/assets/" {
			http.NotFound(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	}))

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(check.Authenticator)

		r.Post("/change-password", user.Password)

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", user.Get)

			r.Route("/leader", func(r chi.Router) {
				r.Use(check.Admin)
				r.Get("/", leader.List)
				r.Post("/", leader.Create)
			})

			r.Route("/member", func(r chi.Router) {
				r.Use(check.Leader)
				r.Get("/", member.List)
				r.Post("/", member.Create)
			})

			r.Route("/administrator", func(r chi.Router) {
				r.Use(check.SuperAdmin)
				r.Get("/", admin.List)
				r.Post("/", admin.Create)
			})

			r.Route("/{userId:[0-9]+}", func(r chi.Router) {
				r.Put("/", user.Edit)
				r.Delete("/", user.Delete)
			})
		})

		r.Route("/video-content", func(r chi.Router) {

			r.Get("/{videoContentId:[0-9]+}/share", content.Shares)

			r.Route("/max-duration", func(r chi.Router) {
				r.Use(check.SuperAdmin)
				r.Put("/", content.MaxDuration)
				r.Get("/", content.GetMaxDuration)
			})

			r.Route("/all", func(r chi.Router) {
				r.Use(check.Admin)
				r.Get("/", content.All)
			})

			r.Route("/{videoContentId:[0-9]+}/ready", func(r chi.Router) {
				r.Use(check.Hub)
				r.Post("/", content.Ready)
			})

			r.Route("/", func(r chi.Router) {
				r.Use(check.Leader)
				r.Get("/", content.Get)

				r.Post("/{videoContentId:[0-9]+}/share/{userId:[0-9]+}", content.Share)
				r.Post("/{videoContentId:[0-9]+}/share", content.MultipleShares)

				r.Post("/id", content.New)

				r.Put("/{videoContentId:[0-9]+}", content.Edit)

				r.Delete("/{videoContentId:[0-9]+}", content.Delete)
			})

			r.Route("/cleaning", func(r chi.Router) {
				r.Use(check.SuperAdmin)
				r.Post("/", content.Cleaning)
			})
		})

		r.Route("/max-duration", func(r chi.Router) {
			r.Use(check.SuperAdmin)
			r.Put("/", content.MaxDuration)
		})

		r.Route("/capture", func(r chi.Router) {
			r.Use(check.Leader)
			r.Post("/play/{videoContentId:[0-9]+}", capture.Play)
			r.Post("/stop", capture.Stop)
			r.Get("/status", capture.Status)
		})
	})
	log.Println(*config.AssetsPath)


	logger.Println("Running server on port", *config.Port)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", *config.Port), r))
}
