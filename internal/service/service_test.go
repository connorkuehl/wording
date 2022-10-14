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
	oldNow := now
	defer func() { now = oldNow }()
	now = func() time.Time {
		n, err := time.Parse(time.Stamp, "Oct 14 6:28:00")
		assert.NilError(t, err)
		return n
	}

	thisInstant := now()

	tokGen := NewMockTokenGenerator(t)
	mockStore := NewMockStore(t)

	tokGen.EXPECT().NewToken().Return("wretched-apostle")

	mockStore.EXPECT().
		CreateGame(mock.Anything, "wretched-apostle", "answer", 3, thisInstant.Add(24*time.Hour)).
		Return(&wording.Game{
			AdminToken: "wretched-apostle",
			Answer:     "answer",
			ExpiresAt:  thisInstant.Add(24 * time.Hour),
			GuessLimit: 3,
		}, nil).
		Once()

	svc := New(mockStore, tokGen, nil)

	got, err := svc.CreateGame(
		context.TODO(),
		"answer",
		3,
		24*time.Hour,
	)
	assert.NilError(t, err)

	want := &wording.Game{
		AdminToken: "wretched-apostle",
		Answer:     "answer",
		ExpiresAt:  thisInstant.Add(24 * time.Hour),
		GuessLimit: 3,
	}

	assert.DeepEqual(t, want, got)
}
