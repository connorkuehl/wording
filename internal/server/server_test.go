package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/connorkuehl/wording/internal/wording"
)

func TestCreateGame(t *testing.T) {
	now := time.Now().UTC()
	svc := NewMockService(t)
	svr := New(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/games", nil)

	q := r.URL.Query()
	q.Add("answer", "potato")
	q.Add("expires_at", fmt.Sprintf("%d", now.Add(12*time.Hour).Unix()))
	q.Add("guess_limit", "6")
	r.URL.RawQuery = q.Encode()

	svc.EXPECT().
		CreateGame(mock.Anything, "potato", 6, now.Add(12*time.Hour).Truncate(time.Second)).
		Return(&wording.Game{
			AdminToken: "wretched-apostle",
			Answer:     "potato",
			ExpiresAt:  now.Add(12 * time.Hour),
			GuessLimit: 6,
		}, nil).
		Once()

	svr.CreateGame(w, r)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.DeepEqual(t, []string{"/manage/wretched-apostle"}, w.Result().Header["Location"])
}
