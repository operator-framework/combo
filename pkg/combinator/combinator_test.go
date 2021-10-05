package combinator

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/operator-framework/combo/pkg/types"
)

// TestCombinations tests the combinator function
func TestCombinations(t *testing.T) {
	// Get needed data for comparison
	testData := getTestData()
	expectedCombos := getExpectedCombos()
	generatedCombos := Solve(testData)

	// Sort the data to prep for a deep equal
	sort.Slice(generatedCombos, func(i, j int) bool {
		return fmt.Sprint(generatedCombos[i]) < fmt.Sprint(generatedCombos[j])
	})

	sort.Slice(expectedCombos, func(i, j int) bool {
		return fmt.Sprint(expectedCombos[i]) < fmt.Sprint(expectedCombos[j])
	})

	// If the outcome and result are not equal, fail
	if !reflect.DeepEqual(generatedCombos, expectedCombos) {
		t.Fatalf(
			"\nCombos generated incorrectly.\n\nHave:\n%v\nWant:\n%v",
			generatedCombos,
			expectedCombos,
		)
	}

}

func getTestData() types.ComboArgs {
	return types.ComboArgs{
		{
			Name: "TEST1",
			Options: []string{
				"foo",
				"bar",
			},
		},
		{
			Name: "TEST2",
			Options: []string{
				"zip",
				"zap",
			},
		},
		{
			Name: "TEST3",
			Options: []string{
				"bip",
				"bap",
			},
		},
	}
}

func getExpectedCombos() types.Combos {
	return types.Combos{
		{
			"TEST1": "foo",
			"TEST2": "zip",
			"TEST3": "bip",
		},
		{
			"TEST1": "foo",
			"TEST2": "zap",
			"TEST3": "bip",
		},
		{
			"TEST1": "bar",
			"TEST2": "zip",
			"TEST3": "bip",
		},
		{
			"TEST1": "bar",
			"TEST2": "zap",
			"TEST3": "bip",
		},
		{
			"TEST1": "foo",
			"TEST2": "zip",
			"TEST3": "bap",
		},
		{
			"TEST1": "foo",
			"TEST2": "zap",
			"TEST3": "bap",
		},
		{
			"TEST1": "bar",
			"TEST2": "zip",
			"TEST3": "bap",
		},
		{
			"TEST1": "bar",
			"TEST2": "zap",
			"TEST3": "bap",
		},
	}
}
