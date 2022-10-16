package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"gotest.tools/assert"

	"github.com/connorkuehl/wording/internal/wording"
)

func TestCreateGame(t *testing.T) {
	tokGen := NewMockTokenGenerator(t)
	admTokGen := NewMockTokenGenerator(t)
	mockStore := NewMockStore(t)

	admTokGen.EXPECT().NewToken().Return("wretched-apostle")
	tokGen.EXPECT().NewToken().Return("hungry-hippo")

	mockStore.EXPECT().
		CreateGame(mock.Anything, "wretched-apostle", "hungry-hippo", "answer", 3).
		Return(&wording.Game{
			AdminToken: "wretched-apostle",
			Token:      "hungry-hippo",
			Answer:     "answer",
			GuessLimit: 3,
		}, nil).
		Once()
	mockStore.EXPECT().
		IncrementStats(mock.Anything, wording.IncrementStats{Stats: wording.Stats{GamesCreated: 1}}).
		Return(nil)

	svc := New(mockStore, admTokGen, tokGen)

	got, err := svc.CreateGame(
		context.TODO(),
		"answer",
		3,
	)
	assert.NilError(t, err)

	want := &wording.Game{
		AdminToken: "wretched-apostle",
		Token:      "hungry-hippo",
		Answer:     "answer",
		GuessLimit: 3,
	}

	assert.DeepEqual(t, want, got)
}
