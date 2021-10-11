package generator

import (
	"os"
	"reflect"
	"testing"
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

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Document combinations generated incorrectly\n\nRecieved:\n\n%s \nbut wanted:\n\n%s", got, want)
			}
		})
	}
}
