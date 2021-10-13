package generator

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var generateTests = []struct {
	name         string
	inputFile    string
	outputFile   string
	combinations []map[string]string
}{
	{"simple input file", "./testdata/input.yaml", "./testData/output.yaml", []map[string]string{
		{
			"REPLACE_ME": "foo",
		},
		{
			"REPLACE_ME": "bar",
		},
	}},
	{"complex input file", "./testdata/complexInput.yaml", "./testData/complexOutput.yaml", []map[string]string{
		{
			"NAMESPACE": "foo",
			"NAME":      "baz",
		},
		{
			"NAMESPACE": "bar",
			"NAME":      "baz",
		},
	}},
}

func TestGenerate(t *testing.T) {
	for _, testCase := range generateTests {
		t.Run(testCase.name, func(t *testing.T) {
			input, err := os.ReadFile(testCase.inputFile)
			if err != nil {
				t.Fatal("Error with test, could not process input test file: ", err.Error())
			}

			want, err := os.ReadFile(testCase.outputFile)
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
