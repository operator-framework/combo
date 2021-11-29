package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	comboErrors "github.com/operator-framework/combo/pkg/errors"
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

func TestValidateFile(t *testing.T) {
	var invalidStream, _ = os.Open("./DOES_NOT_EXIST")
	for _, tt := range []struct {
		name  string
		input io.Reader
		err   error
	}{
		{
			name:  "validates a valid file correctly",
			input: strings.NewReader(`foo: bar`),
			err:   nil,
		},
		{
			name:  "invalidates an empty file",
			input: strings.NewReader(""),
			err:   ErrEmptyFile,
		},
		{
			name:  "invalidates an unreadable file",
			input: invalidStream,
			err:   comboErrors.ErrCouldNotReadFile,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFile(tt.input)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
