package combinator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var combinationTests = []struct {
	name  string
	input map[string][]string
	want  []map[string]string
}{
	{"empty map input", map[string][]string{}, []map[string]string{}},
	{"standard set of args", getTestData(), getExpectedCombos()},
}

// TestCombinations tests the combinator function
func TestCombinations(t *testing.T) {
	for _, testCase := range combinationTests {
		t.Run(testCase.name, func(t *testing.T) {
			got := Solve(testCase.input)
			require.ElementsMatch(t, got, testCase.want, "Combos generated incorrectly")
		})
	}
}

func getTestData() map[string][]string {
	return map[string][]string{
		"TEST1": {"foo", "bar"},
		"TEST2": {"zip", "zap"},
		"TEST3": {"bip", "bap"},
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
