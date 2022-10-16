package wording

import (
	"testing"

	"gotest.tools/assert"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		guess  string
		answer string
		want   Attempt
	}{
		{
			guess:  "potato",
			answer: "potato",
			want: Attempt{
				Character{Value: "p", IsCorrect: true},
				Character{Value: "o", IsCorrect: true},
				Character{Value: "t", IsCorrect: true},
				Character{Value: "a", IsCorrect: true},
				Character{Value: "t", IsCorrect: true},
				Character{Value: "o", IsCorrect: true},
			},
		},
		{
			guess:  "bzc",
			answer: "abc",
			want: Attempt{
				Character{Value: "b", IsPartial: true},
				Character{Value: "z", IsCorrect: false, IsPartial: false},
				Character{Value: "c", IsCorrect: true},
			},
		},
		{
			guess:  "eee",
			answer: "aaa",
			want: Attempt{
				Character{Value: "e", IsCorrect: false, IsPartial: false},
				Character{Value: "e", IsCorrect: false, IsPartial: false},
				Character{Value: "e", IsCorrect: false, IsPartial: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.guess, func(t *testing.T) {
			assert.DeepEqual(t, tt.want, Evaluate(tt.answer, tt.guess))
		})
	}
}
