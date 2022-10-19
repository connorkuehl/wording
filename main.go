package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/connorkuehl/wording/internal/generator"
	"github.com/connorkuehl/wording/internal/server"
	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/store"
)

func main() {
	var config struct {
		baseURL string
		dsn     string
		bind    string
	}

	flag.StringVar(&config.baseURL, "base-url", "http://localhost:8080", "Base URL to prefix links with")
	flag.StringVar(&config.dsn, "db-dsn", os.Getenv("WORDING_DB_DSN"), "Postgres DSN")
	flag.StringVar(&config.bind, "bind-addr", os.Getenv("WORDING_BIND_ADDR"), "Bind address")
	flag.Parse()

	log.WithFields(log.Fields{
		"bind-addr": config.bind,
		"base-url":  config.baseURL,
	}).Info("starting up")

	store, err := store.NewPostgresStore(config.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	log.Info("connected to database")

	rand.Seed(time.Now().UnixNano())
	adminTokenGenerator := generator.NewUUIDGenerator()
	gameTokenGenerator := generator.NewHumanReadableGenerator(rand.Int)

	svc := service.New(store, adminTokenGenerator, gameTokenGenerator)
	srv := server.New(config.baseURL, svc)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", srv.Home)
	router.Get("/manage/{admin_token}", srv.ManageGame)
	router.Post("/games", srv.CreateGame)
	router.Get("/game/{token}", srv.PlayGame)
	router.Post("/game/{token}", srv.Guess)
	router.Post("/manage/{admin_token}/delete", srv.DeleteGame)
	router.Get("/health", func(_ http.ResponseWriter, _ *http.Request) {
	})

	log.Info("listening")
	log.Fatal(http.ListenAndServe(config.bind, router))
}
