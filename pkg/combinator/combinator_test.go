package combinator

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
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

func getTestData() map[string]string {
	return map[string]string{
		"TEST1": "foo,bar",
		"TEST2": "zip,zap",
		"TEST3": "bip,bap",
	}
}

func getExpectedCombos() []map[string]string {
	return []map[string]string{
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
