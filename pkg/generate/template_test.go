package generate

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildTemplate(t *testing.T) {
	for _, tt := range []struct {
		name     string
		file     io.Reader
		expected []string
		err      error
	}{
		{
			name: "can find a value if present",

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
	} {
		t.Run(tt.name, func(t *testing.T) {
			actualTemplate, err := buildTemplate(tt.file)
			if !errors.Is(err, tt.err) {
				t.Fatal("error with test, not able to create template:", err)
			}

			actual := actualTemplate.documents

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
				processedDocuments: []string{
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
				processedDocuments: []string{
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
			name:  "can create a template correctly with a combo",
			combo: map[string]string{"NAMESPACE": "foo", "NAME": "baz"},
			template: template{
				documents: []string{
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
			name:  "responds correctly given an empty combo",
			combo: map[string]string{},
			template: template{
				documents: []string{
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

			require.Equal(t, tt.expected, tt.template.processedDocuments)
		})
	}
}
