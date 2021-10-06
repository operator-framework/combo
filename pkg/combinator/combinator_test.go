package combinator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCombinations tests the combinator function
func TestCombinations(t *testing.T) {
	// Get needed data for comparison
	testData := getTestData()
	expectedCombos := getExpectedCombos()
	generatedCombos := Solve(testData)

	require.ElementsMatch(t, generatedCombos, expectedCombos, "Combos generated incorrectly")
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
