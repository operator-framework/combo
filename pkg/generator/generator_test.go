package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/operator-framework/combo/pkg/combination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err        error
	evaluation string
}

func TestEvaluate(t *testing.T) {
	for _, tt := range []struct {
		name         string
		template     string
		combinations combination.Stream
		expected     expected
	}{
		{
			name: "can process a template correctly",
			template: `---
name: test
---
name: replacement1
test: REPLACE_ME_1
---
name: replacement2
test: REPLACE_ME_2
`,
			expected: expected{
				err: nil,
				evaluation: `---
name: test
---
name: replacement1
test: foo
---
name: replacement2
test: zip
---
name: replacement1
test: foo
---
name: replacement2
test: zap
---
name: replacement1
test: bar
---
name: replacement2
test: zap
---
name: replacement1
test: bar
---
name: replacement2
test: zip
`,
			},
			combinations: combination.NewStream(
				combination.WithArgs(map[string][]string{
					"REPLACE_ME_1": {"foo", "bar"},
					"REPLACE_ME_2": {"zip", "zap"},
				}),
				combination.WithSolveAhead(),
			),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			evaluation, err := Evaluate(ctx, tt.template, tt.combinations)
			if err != nil {
				t.Fatal("recieved an error during evaluation:", err)
			}

			assert.Equal(t, tt.expected.err, err)
			require.ElementsMatch(
				t,
				strings.Split(string(tt.expected.evaluation), "---"),
				strings.Split(string(evaluation), "---"),
				"Document evaluations generated incorrectly")
		})
	}
}
