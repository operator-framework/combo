package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatReplacements(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input map[string]string
		want  map[string][]string
	}{
		{
			name:  "formats replacements correctly",
			input: map[string]string{"TEST": "foo,bar,bap"},
			want:  map[string][]string{"TEST": {"foo", "bar", "bap"}},
		},
		{
			name:  "handles empty values",
			input: map[string]string{"TEST": ""},
			want:  map[string][]string{"TEST": {""}},
		},
		{
			name:  "handles empty input",
			input: map[string]string{},
			want:  map[string][]string{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, formatReplacements(tt.input), "replacements formatted incorrectly")
		})
	}
}
