package combination

import (
	"context"
	"errors"
	"testing"

	testdata "github.com/operator-framework/combo/test/assets/combination"
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
		input: testdata.CombinationInput,
		expected: expected{
			combinations: testdata.CombinationOutput,
			err:          nil,
		},
	},
	{
		name:  "standard set of long args",
		input: testdata.LongCombinationInput,
		expected: expected{
			combinations: testdata.LongCombinationOutput,
			err:          nil,
		},
	},
}

func TestNext(t *testing.T) {
	// Test 1: Uses testIterative()
	for _, tt := range combinationTests {
		t.Run(tt.name, func(t *testing.T) {
			combinationStream := NewStream(
				WithArgs(tt.input),
			)

			var got []map[string]string

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			for {
				next, err := combinationStream.Next(ctx)
				if !errors.Is(err, tt.expected.err) {
					t.Fatal("[TEST nextIterative()] error received while processing combination stream:", err)
				}

				if next == nil {
					break
				}

				got = append(got, next)
			}
			require.ElementsMatch(t, got, tt.expected.combinations, "[TEST nextIterative()] Combos generated incorrectly")
		})
	}

	// Test 2: Uses nextPreSolvedCombinations()
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
				if !errors.Is(err, tt.expected.err) {
					t.Fatal("[TEST nextPreSolvedCombinations()] error received while processing combination stream:", err)
				}

				if next == nil {
					break
				}

				got = append(got, next)
			}
			require.ElementsMatch(t, got, tt.expected.combinations, "[TEST nextPreSolvedCombinations()] Combos generated incorrectly")
		})
	}
}
