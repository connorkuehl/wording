package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/connorkuehl/wording/internal/generator"
	"github.com/connorkuehl/wording/internal/server"
	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/store"
)

func main() {
	var config struct {
		baseURL string
		dsn     string
	}

	flag.StringVar(&config.baseURL, "base-url", "http://localhost:8080", "Base URL to prefix links with")
	flag.StringVar(&config.dsn, "db-dsn", os.Getenv("WORDING_DB_DSN"), "Postgres DSN")

	flag.Parse()

	store, err := store.NewPostgresStore(config.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

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

	log.Fatal(http.ListenAndServe(":8080", router))
}
