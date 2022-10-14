package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/connorkuehl/wording/internal/wording"
)

func TestCreateGame(t *testing.T) {
	now := time.Now().UTC()
	svc := NewMockService(t)
	svr := New("http://localhost:8080", svc)

	form := url.Values{
		"answer":        {"potato"},
		"expires_after": {(12 * time.Hour).String()},
		"num_attempts":  {"6"},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/games", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	svc.EXPECT().
		CreateGame(mock.Anything, "potato", 6, (12*time.Hour)).
		Return(&wording.Game{
			AdminToken: "wretched-apostle",
			Answer:     "potato",
			ExpiresAt:  now.Add(12 * time.Hour),
			GuessLimit: 6,
		}, nil).
		Once()

	svr.CreateGame(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code, w.Body)
	assert.DeepEqual(t, []string{"/manage/wretched-apostle"}, w.Result().Header["Location"])
}
