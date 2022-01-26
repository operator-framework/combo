package template

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/operator-framework/combo/pkg/combination"
	testdata "github.com/operator-framework/combo/test/assets/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	for _, tt := range []struct {
		name         string
		file         io.Reader
		combinations combination.Stream
		expected     []string
		err          error
	}{
		{
			name:     "can process a template",
			file:     strings.NewReader(testdata.BuildInput),
			expected: testdata.BuildOutput,
			err:      nil,
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(true),
			),
		},
		{
			name:     "processes an empty template",
			file:     strings.NewReader(``),
			expected: []string{},
			err:      nil,
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(true),
			),
		},
		{
			name:     "processes an empty template",
			file:     strings.NewReader(``),
			expected: []string{},
			err:      nil,

			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"NAMESPACE": {"foo", "bar"},
					"NAME":      {"baz"},
				}),
				combination.WithSolveAhead(true),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			templateBuilder, err := NewBuilder(tt.file, tt.combinations)
			if !errors.Is(err, tt.err) {
				t.Fatalf("received an error while building generator: %v", err)
			}

			actual, err := templateBuilder.Build(ctx)
			if !errors.Is(err, tt.err) {
				t.Fatalf("received an error during evaluation: %v", err)
			}

			assert.Equal(t, tt.err, err)

			require.ElementsMatch(t, tt.expected, actual)
		})
	}
}
