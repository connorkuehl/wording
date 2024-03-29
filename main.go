package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/connorkuehl/wording/internal/generator"
	"github.com/connorkuehl/wording/internal/randword"
	"github.com/connorkuehl/wording/internal/server"
	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/store"
)

func main() {
	var config struct {
		environment string
		baseURL     string
		dbDSN       string
		bind        string
		wordGenSvc  string
	}

	fromEnvOr := func(key, fallback string) string {
		s := os.Getenv(key)
		if s == "" {
			return fallback
		}
		return s
	}

	flag.StringVar(&config.environment, "environment", fromEnvOr("WORDING_ENVIRONMENT", "dev"), "Environment")
	flag.StringVar(&config.baseURL, "base-url", os.Getenv("WORDING_BASE_URL"), "Base URL to prefix links with")
	flag.StringVar(&config.dbDSN, "db-dsn", os.Getenv("WORDING_DB_DSN"), "Postgres DSN")
	flag.StringVar(&config.bind, "bind-addr", os.Getenv("WORDING_BIND_ADDR"), "Bind address")
	flag.StringVar(&config.wordGenSvc, "word-gen-svc", os.Getenv("WORDING_WORD_GEN_SVC"), "Word generator API")
	flag.Parse()

	log.WithFields(log.Fields{
		"bind-addr":    config.bind,
		"base-url":     config.baseURL,
		"word-gen-svc": config.wordGenSvc,
	}).Info("starting up")

	store, err := store.NewPostgresStore(config.dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	log.Info("connected to database")

	adminTokenGenerator := generator.NewUUIDGenerator()
	gameTokenGenerator := generator.NewFallibleGenerator(
		generator.NewHumanReadable(randword.NewClient(config.wordGenSvc)),
		generator.NewUUIDGenerator(),
	)

	var svc service.Service = service.New(store, adminTokenGenerator, gameTokenGenerator)
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
