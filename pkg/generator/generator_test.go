package generator

import (
	"os"
	"reflect"
	"testing"

	"github.com/operator-framework/combo/pkg/types"
)

func TestGenerate(t *testing.T) {
	input, err := os.ReadFile("./testData/input.yaml")
	if err != nil {
		t.Fatal("Error with test, could not process input test file: ", err.Error())
	}

	desiredResult, err := os.ReadFile("./testData/output.yaml")
	if err != nil {
		t.Fatal("Error with test, could not process output test file: ", err.Error())
	}

	result, err := Generate(types.Combos{
		{
			"REPLACE_ME": "foo",
		},
		{
			"REPLACE_ME": "bar",
		},
	}, input)

	if err != nil {
		t.Fatalf("Recieved an error while running Generate(): %v", err)
	}

	if !reflect.DeepEqual(result, desiredResult) {
		t.Fatalf("Document combinations generated incorrectly\n\nRecieved:\n\n%s \nbut wanted:\n\n%s", result, desiredResult)
	}

}
