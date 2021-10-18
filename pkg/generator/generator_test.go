package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expected struct {
	err        error
	evaluation string
}

func TestGenerate(t *testing.T) {
	for _, tt := range []struct {
		name         string
		template     string
		combinations []map[string]string
		expected     expected
	}{
		{
			name: "can process a template correctly",
			template: `---
name: test
---
name: hello
test: REPLACE_ME
---
name: world
test: REPLACE_ME
`,
			expected: expected{
				err: nil,
				evaluation: `---
name: test
---
name: hello
test: foo
---
name: hello
test: bar
---
name: world
test: foo
---
name: world
test: bar
`,
			},
			combinations: []map[string]string{
				{
					"REPLACE_ME": "foo",
				},
				{
					"REPLACE_ME": "bar",
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			evaluation, err := Generate(tt.combinations, []byte(tt.template))
			assert.Equal(t, tt.expected.err, err)
			require.ElementsMatch(
				t,
				strings.Split(string(tt.expected.evaluation), "---"),
				strings.Split(string(evaluation), "---"),
				"Document evaluations generated incorrectly")
		})
	}
}
