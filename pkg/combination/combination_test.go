package combination

import (
	"context"
	"testing"

	"github.com/operator-framework/combo/test/assets/combinationTestData"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err          error
	combinations []map[string]string
}

var combinationTests = []struct {
	name     string
	input    map[string][]string
	expected expected
}{
	{
		name:  "empty map input",
		input: map[string][]string{},
		expected: expected{
			combinations: []map[string]string{},
			err:          ErrNoArgsSet,
		},
	},
	{
		name:  "standard set of args",
		input: combinationTestData.CombinationInput,
		expected: expected{
			combinations: combinationTestData.CombinationOutput,
			err:          nil,
		},
	},
}

func TestAll(t *testing.T) {
	for _, tt := range combinationTests {
		t.Run(tt.name, func(t *testing.T) {
			combinationStream := NewStream(
				WithArgs(tt.input),
				WithSolveAhead(),
			)
			got, err := combinationStream.All()
			if err != tt.expected.err {
				t.Fatal("error received while retreiving all combinations:", err)
			}
			require.ElementsMatch(t, got, tt.expected.combinations, "Combos generated incorrectly")
		})
	}
}

func TestNext(t *testing.T) {
	for _, tt := range combinationTests {
		t.Run(tt.name, func(t *testing.T) {
			combinationStream := NewStream(
				WithArgs(tt.input),
				WithSolveAhead(),
			)

			var got []map[string]string

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			for {
				next, err := combinationStream.Next(ctx)
				if err != tt.expected.err {
					t.Fatal("error received while processing combination stream:", err)
				}

				if next == nil {
					break
				}

				got = append(got, next)
			}
			require.ElementsMatch(t, got, tt.expected.combinations, "Combos generated incorrectly")
		})
	}
}
