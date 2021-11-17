package generate

import (
	"context"
	"testing"

	"github.com/operator-framework/combo/pkg/combination"
	testdata "github.com/operator-framework/combo/test/assets/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err        error
	evaluation []string
}

func TestEvaluate(t *testing.T) {
	for _, tt := range []struct {
		name         string
		template     string
		combinations combination.Stream
		expected     expected
	}{
		{
			name:     "can process a template",
			template: testdata.EvaluateInput,
			expected: expected{
				err:        nil,
				evaluation: testdata.EvaluateOutput,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
		{
			name:     "processes an empty template",
			template: ``,
			expected: expected{
				err:        nil,
				evaluation: []string{},
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			generator := NewGenerator(tt.template, tt.combinations)

			evaluation, err := generator.Evaluate(ctx)
			if err != nil {
				t.Fatal("received an error during evaluation:", err)
			}

			assert.Equal(t, tt.expected.err, err)

			require.ElementsMatch(t, tt.expected.evaluation, evaluation)
		})
	}
}
