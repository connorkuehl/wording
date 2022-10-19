package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
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
		environment string
		baseURL     string
		dbDSN       string
		bind        string
		wordGenSvc  string
		sentryDSN   string
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
	flag.StringVar(&config.sentryDSN, "sentry-dsn", os.Getenv("WORDING_SENTRY_DSN"), "sentry.io DSN")
	flag.Parse()

	log.WithFields(log.Fields{
		"bind-addr":    config.bind,
		"base-url":     config.baseURL,
		"word-gen-svc": config.wordGenSvc,
		"use-sentry":   config.sentryDSN != "",
	}).Info("starting up")

	if config.sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:         config.sentryDSN,
			Environment: config.environment,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	store, err := store.NewPostgresStore(config.dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	log.Info("connected to database")

	rand.Seed(time.Now().UnixNano())
	adminTokenGenerator := generator.NewUUIDGenerator()
	gameTokenGenerator := generator.NewRandomWordClient(
		config.wordGenSvc,
		generator.NewHumanReadableGenerator(rand.Int),
	)

	svc := service.New(store, adminTokenGenerator, gameTokenGenerator)
	srv := server.New(config.baseURL, service.NewSentry(svc))

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
