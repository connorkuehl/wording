package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/connorkuehl/wording/internal/wording"
)

//go:generate mockery --name Service --case underscore --with-expecter --testonly --inpackage
type Service interface {
	CreateGame(ctx context.Context, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
}

type Server struct {
	svc Service
}

func New(svc Service) *Server {
	return &Server{
		svc: svc,
	}
}

func (s *Server) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	// TODO
	// This whole section of pulling query parameters out is awful.
	keys, ok := r.URL.Query()["answer"]
	if !ok || len(keys) == 0 {
		http.Error(w, "answer is missing", http.StatusBadRequest)
		return
	}
	answer := keys[0]

	keys, ok = r.URL.Query()["guess_limit"]
	if !ok || len(keys) == 0 {
		http.Error(w, "guess_limit is missing", http.StatusBadRequest)
		return
	}
	guessLimit, err := strconv.Atoi(keys[0])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	keys, ok = r.URL.Query()["expires_at"]
	if !ok || len(keys) == 0 {
		http.Error(w, "expires_at is missing", http.StatusBadRequest)
		return
	}
	expiresAt, err := strconv.ParseInt(keys[0], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	expiresAtTime := time.Unix(expiresAt, 0).UTC()

	game, err := s.svc.CreateGame(ctx, answer, guessLimit, expiresAtTime)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/manage/%s", game.AdminToken), http.StatusMovedPermanently)
}

func (s *Server) ManageGame(w http.ResponseWriter, r *http.Request) {

}
