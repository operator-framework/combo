package generator

import (
	"strings"
	"testing"

	"testing/fstest"

	"github.com/stretchr/testify/require"
)

const (
	INPUT_PATH  = "input.yaml"
	OUTPUT_PATH = "output.yaml"
)

var generateTests = []struct {
	name         string
	fileSystem   fstest.MapFS
	combinations []map[string]string
}{
	{
		"simple input file",
		fstest.MapFS{
			INPUT_PATH: &fstest.MapFile{
				Data: []byte(`---
name: test
---
name: hello
test: REPLACE_ME
---
name: world
test: REPLACE_ME`,
				),
			},
			OUTPUT_PATH: &fstest.MapFile{
				Data: []byte(`---
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
				),
			},
		},
		[]map[string]string{
			{
				"REPLACE_ME": "foo",
			},
			{
				"REPLACE_ME": "bar",
			},
		},
	},
}

func TestGenerate(t *testing.T) {
	for _, testCase := range generateTests {
		t.Run(testCase.name, func(t *testing.T) {
			input, err := testCase.fileSystem.ReadFile(INPUT_PATH)
			if err != nil {
				t.Fatal("Error with test, could not process input test file: ", err.Error())
			}

			want, err := testCase.fileSystem.ReadFile(OUTPUT_PATH)
			if err != nil {
				t.Fatal("Error with test, could not process output test file: ", err.Error())
			}

			got, err := Generate(testCase.combinations, input)
			if err != nil {
				t.Fatalf("Recieved an error while running Generate(): %v", err)
			}

			require.ElementsMatch(t, strings.Split(string(got), "---"), strings.Split(string(want), "---"), "Document combinations generated incorrectly")
		})
	}
}
