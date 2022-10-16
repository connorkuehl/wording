package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/view"
	"github.com/connorkuehl/wording/internal/wording"
)

const playerTokenCookie = "WordingToken"

//go:generate mockery --name Service --case underscore --with-expecter --testonly --inpackage
type Service interface {
	CreateGame(ctx context.Context, answer string, guessLimit int, expiresAfter time.Duration) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
	SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error
	GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error)
	NewPlayerToken(ctx context.Context) string
	Stats(ctx context.Context) (wording.Stats, error)
	DeleteGame(ctx context.Context, adminToken string) error
}

type Server struct {
	baseURL string
	svc     Service
}

func New(baseURL string, svc Service) *Server {
	return &Server{
		baseURL: baseURL,
		svc:     svc,
	}
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	stats, err := s.svc.Stats(ctx)
	if err != nil {
		// TODO
		log.Println("read stats:", err)
	}

	err = view.Home{Stats: stats}.RenderTo(w)
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

	var invalidInput wording.InputViolations
	if errors.As(err, &invalidInput) {
		http.Error(w, http.StatusText(http.StatusBadRequest)+fmt.Sprintf(": %v", err), http.StatusBadRequest)
		return
	}
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
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = view.ManageGame{
		BaseURL:        s.baseURL,
		AdminToken:     game.AdminToken,
		Token:          game.Token,
		Answer:         game.Answer,
		GuessesAllowed: game.GuessLimit,
		ExpiresAt:      game.ExpiresAt,
	}.RenderTo(w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) PlayGame(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	token := chi.URLParam(r, "token")
	if token == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	game, err := s.svc.GameByToken(ctx, token)
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var id string
	idCookie, err := r.Cookie(playerTokenCookie)
	if err != nil {
		id = s.svc.NewPlayerToken(ctx)
	} else {
		id = idCookie.Value
	}

	state, err := s.svc.GameState(ctx, game.Token, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = view.PlayGame{
		Token:     token,
		Length:    len(game.Answer),
		GameState: state,
	}.RenderTo(w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) Guess(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	token := chi.URLParam(r, "token")

	var id string
	idCookie, err := r.Cookie(playerTokenCookie)
	if err != nil {
		id = s.svc.NewPlayerToken(ctx)
		http.SetCookie(w, &http.Cookie{Name: playerTokenCookie, Value: id})
	} else {
		id = idCookie.Value
	}

	r.ParseForm()

	guess := r.PostForm.Get("guess")

	err = s.svc.SubmitGuess(ctx, token, id, guess)
	var violations wording.InputViolations
	if errors.As(err, &violations) {
		http.Error(w, http.StatusText(http.StatusBadRequest)+fmt.Sprintf(": %v", err), http.StatusBadRequest)
		return
	}
	if errors.Is(err, service.ErrGuessLimitReached) || errors.Is(err, service.ErrCannotContinue) {
		http.Redirect(w, r, fmt.Sprintf("/game/%s", token), http.StatusSeeOther)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/game/%s", token), http.StatusSeeOther)
}

func (s *Server) DeleteGame(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	token := chi.URLParam(r, "admin_token")

	err := s.svc.DeleteGame(ctx, token)
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
