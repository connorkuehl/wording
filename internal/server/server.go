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
	CreateGame(ctx context.Context, answer string, guessLimit int, expiresAfter time.Duration) (*wording.Game, error)
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

	var (
		answer       string
		expiresAfter time.Duration
		numAttempts  int
	)

	r.ParseForm()

	answer = r.PostFormValue("answer")
	if answer == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" answer is missing", http.StatusBadRequest)
		return
	}

	d, err := time.ParseDuration(r.PostFormValue("expires_after"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" expires_after is missing", http.StatusBadRequest)
		return
	}
	expiresAfter = d

	i, err := strconv.Atoi(r.PostFormValue("num_attempts"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" num_attempts is not a number", http.StatusBadRequest)
		return
	}
	numAttempts = i

	game, err := s.svc.CreateGame(ctx, answer, numAttempts, expiresAfter)
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
