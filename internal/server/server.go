package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/connorkuehl/wording/internal/store"
	"github.com/connorkuehl/wording/internal/view"
	"github.com/connorkuehl/wording/internal/wording"
)

//go:generate mockery --name Service --case underscore --with-expecter --testonly --inpackage
type Service interface {
	CreateGame(ctx context.Context, answer string, guessLimit int, expiresAt time.Time) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
}

type Server struct {
	svc Service
}

func New(svc Service) *Server {
	return &Server{
		svc: svc,
	}
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	_, err := io.Copy(w, strings.NewReader(view.Home()))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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

	http.Redirect(w, r, fmt.Sprintf("/manage/%s", game.AdminToken), http.StatusSeeOther)
}

func (s *Server) ManageGame(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	adminToken := chi.URLParam(r, "admin_token")
	if adminToken == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	game, err := s.svc.Game(ctx, adminToken)
	if errors.Is(err, store.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bodyFmt := `<html>
	<p>The answer is <b>%s</b>. Players will have %d guesses.</p>
	<p>The game expires at %s.</p>
	</html>`

	w.Write([]byte(fmt.Sprintf(bodyFmt, game.Answer, game.GuessLimit, game.ExpiresAt.Format(time.UnixDate))))
}
