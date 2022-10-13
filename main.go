package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/connorkuehl/wording/internal/server"
	"github.com/connorkuehl/wording/internal/service"
)

func main() {
	svc := service.New(nil, nil, nil)
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
