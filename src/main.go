package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tempgalias/src/config"
	"tempgalias/src/database"
	"tempgalias/src/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// https://github.com/veryhappytree/go-boilerplate/
// https://stackoverflow.com/questions/64510093/gorm-migration-using-golang-migrate-migrate

func main() {
	config.LoadConfig()
	database.SetupDatabase()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//public routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
		r.Post("/get-alias", routes.GetAliasHandler)
	})

	http.ListenAndServe(":3000", r)

	gracefulShutdown(
		func() error {
			database.DB.Close()
			return nil
		},
		// func() error {
		// 	return redis.Client.Close()
		// },
		// func() error {
		// 	return rabbit.Service.Channel.Close()
		// },
		func() error {
			os.Exit(0)
			return nil
		},
	)
}

func gracefulShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				slog.Error("gracefulShutdown op failed", "error", err)
				panic(err)
			}
		}
	}
}
