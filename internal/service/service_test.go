package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/connorkuehl/wording/internal/wording"
)

func TestCreateGame(t *testing.T) {
	now := time.Now()
	nower := func() time.Time { return now }
	tokGen := NewMockTokenGenerator(t)
	mockStore := NewMockStore(t)

	tokGen.EXPECT().Adjective().Return("wretched")
	tokGen.EXPECT().Noun().Return("apostle")

	mockStore.EXPECT().
		CreateGame(mock.Anything, "wretched-apostle", "answer", 3, now.Add(24*time.Hour)).
		Return(&wording.Game{
			AdminToken: "wretched-apostle",
			Answer:     "answer",
			ExpiresAt:  now.Add(24 * time.Hour),
			GuessLimit: 3,
		}, nil).
		Once()

	svc := New(mockStore, tokGen, nower)

	got, err := svc.CreateGame(
		context.TODO(),
		"answer",
		3,
		now.Add(24*time.Hour),
	)
	assert.NilError(t, err)

	want := &wording.Game{
		AdminToken: "wretched-apostle",
		Answer:     "answer",
		ExpiresAt:  now.Add(24 * time.Hour),
		GuessLimit: 3,
	}

	assert.DeepEqual(t, want, got)
}
