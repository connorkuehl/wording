package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/connorkuehl/wording/internal/service"
	"github.com/connorkuehl/wording/internal/view"
	"github.com/connorkuehl/wording/internal/wording"
)

const playerTokenCookie = "WordingToken"

//go:generate mockery --name Service --case underscore --with-expecter --testonly --inpackage
type Service interface {
	CreateGame(ctx context.Context, answer string, guessLimit int) (*wording.Game, error)
	Game(ctx context.Context, adminToken string) (*wording.Game, error)
	GameByToken(ctx context.Context, token string) (*wording.Game, error)
	SubmitGuess(ctx context.Context, gameToken, playerToken, guess string) error
	GameState(ctx context.Context, gameToken, playerToken string) (*wording.GameState, error)
	NewPlayerToken(ctx context.Context) string
	Stats(ctx context.Context) (wording.Stats, error)
	DeleteGame(ctx context.Context, adminToken string) error
	GameStats(ctx context.Context, adminToken string) (wording.Stats, error)
}

// Server is the HTTP "edge" of the web application.
type Server struct {
	baseURL string
	svc     Service
}

// New creates a new Server.
func New(baseURL string, svc Service) *Server {
	return &Server{
		baseURL: baseURL,
		svc:     svc,
	}
}

// Home renders the root of the site, which is the game creation page.
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

// CreateGame handles incoming POST forms for creating a new game.
func (s *Server) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	var (
		answer      string
		numAttempts int
	)

	_ = r.ParseForm()

	answer = r.PostFormValue("answer")
	if answer == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" answer is missing", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(r.PostFormValue("num_attempts"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+" num_attempts is not a number", http.StatusBadRequest)
		return
	}
	numAttempts = i

	game, err := s.svc.CreateGame(ctx, answer, numAttempts)

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

// ManageGame renders the manage game page, where a game's stats are shown and it
// can optionally be deleted.
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

	stats, err := s.svc.GameStats(ctx, adminToken)
	if err != nil {
		// TODO
		log.Println(err)
	}

	err = view.ManageGame{
		BaseURL:        s.baseURL,
		AdminToken:     game.AdminToken,
		Token:          game.Token,
		Answer:         game.Answer,
		GuessesAllowed: game.GuessLimit,
		GuessesMade:    stats.GuessesMade,
		CorrectGuesses: stats.GamesWon,
	}.RenderTo(w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// PlayGame renders the play game page.
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

	for _, attempt := range state.Attempts {
		for i := range attempt {
			attempt[i].Value = strings.ToUpper(attempt[i].Value)
		}
	}

	padding := game.GuessLimit - len(state.Attempts)
	for i := 0; i < padding; i++ {
		state.Attempts = append(state.Attempts, wording.Attempt{})
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

// Guess handles the POST form data for a player submitting a guess for a game.
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

	_ = r.ParseForm()

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
		log.Println(err)
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
