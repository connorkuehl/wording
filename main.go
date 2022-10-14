package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/connorkuehl/wording/internal/server"
	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/store"
)

func main() {
	var config struct {
		dsn string
	}

	flag.StringVar(&config.dsn, "db-dsn", os.Getenv("WORDING_DB_DSN"), "Postgres DSN")

	flag.Parse()

	store, err := store.NewPostgresStore(config.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	svc := service.New(store, nil)
	srv := server.New(svc)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", srv.Home)
	router.Get("/manage/{admin_token}", srv.ManageGame)
	router.Post("/games", srv.CreateGame)

	log.Fatal(http.ListenAndServe(":8080", router))
}
