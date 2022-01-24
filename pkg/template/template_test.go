package template

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	comboErrors "github.com/operator-framework/combo/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewTemplate(t *testing.T) {
	// Create an invalid stream to verify that buildTemplate() bubbles up the
	// error correctly if a file is not readable
	var invalidStream, _ = os.Open("./DOES_NOT_EXIST")

	for _, tt := range []struct {
		name     string
		file     io.Reader
		expected []string
		err      error
	}{
		{
			name: "builds a template correctly",

			file: strings.NewReader(`---
testOne: 123
---
testTwo: 456
`),
			expected: []string{"testOne: 123", "testTwo: 456"},
			err:      nil,
		},
		{
			name:     "processes an empty io.Reader correctly",
			file:     strings.NewReader(""),
			expected: nil,
			err:      nil,
		},
		{
			name:     "returns ErrCouldNotReadFile if io.Reader is not readable",
			file:     invalidStream,
			expected: nil,
			err:      comboErrors.ErrCouldNotReadFile,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actualTemplate, err := newTemplate(tt.file)
			if !errors.Is(err, tt.err) {
				t.Fatal("error with test, not able to create template:", err)
			}

			actual := actualTemplate.manifests

			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestHas(t *testing.T) {
	for _, tt := range []struct {
		name     string
		template template
		find     string
		expected bool
		err      error
	}{
		{
			name:     "can find a value if present",
			expected: true,
			find:     "testOne: 123",
			template: template{
				processedManifests: []string{
					"testOne: 123",
					"testTwo: 456",
				},
			},
		},
		{
			name:     "does not find a value if not present",
			expected: false,
			find:     "testThree: 789",
			template: template{
				processedManifests: []string{
					"testOne: 123",
					"testTwo: 456",
				},
			},
		},
		{
			name:     "responds correctly given an empty template",
			expected: false,
			find:     "testThree: 789",
			template: template{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.template.has(tt.find)
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestWith(t *testing.T) {
	for _, tt := range []struct {
		name     string
		combo    map[string]string
		template template
		expected []string
	}{
		{
			name:  "processes the template given a combo",
			combo: map[string]string{"NAMESPACE": "foo", "NAME": "baz"},
			template: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
				},
			},
			expected: []string{
				"testOne: foo",
				"testTwo: baz",
			},
		},

		{
			name:  "processes the template correctly given a combo not present in template",
			combo: map[string]string{"NAMESPACE": "foo", "NOT_PRESENT": "baz"},
			template: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
				},
			},
			expected: []string{
				"testOne: foo",
				"testTwo: NAME",
			},
		},
		{
			name:  "processes manifests that have no replacements made",
			combo: map[string]string{"NAMESPACE": "foo", "NAME": "baz"},
			template: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
					"testThree: 789",
				},
			},
			expected: []string{
				"testOne: foo",
				"testTwo: baz",
				"testThree: 789",
			},
		},
		{
			name:  "does not process an empty manifest",
			combo: map[string]string{"NAMESPACE": "foo", "NAME": "baz"},
			template: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
					"",
				},
			},
			expected: []string{
				"testOne: foo",
				"testTwo: baz",
			},
		},
		{
			name:  "responds correctly given an empty combo",
			combo: map[string]string{},
			template: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
				},
			},
			expected: []string{
				"testOne: NAMESPACE",
				"testTwo: NAME",
			},
		},
		{
			name:     "responds correctly given an empty template",
			combo:    map[string]string{"NAMESPACE": "foo", "NAME": "baz"},
			template: template{},
			expected: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.template.with(tt.combo)

			require.Equal(t, tt.expected, tt.template.processedManifests)
		})
	}
}

func TestValidateFile(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input template
		err   error
	}{
		{
			name: "validates a valid file correctly",
			input: template{
				manifests: []string{
					"testOne: NAMESPACE",
					"testTwo: NAME",
				},
			},
			err: nil,
		},
		{
			name: "validates an empty file",
			input: template{
				processedManifests: []string{},
			},
			err: nil,
		},
		{
			name: "invalidates an unreadable file",
			input: template{
				manifests: []string{
					"	apiVersion: rbac.authorization.k8s.io/v1",
					"kind: ClusterRole",
				},
			},
			err: ErrInvalidYAML,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			require.ErrorIs(t, err, tt.err)
		})
	}
}
