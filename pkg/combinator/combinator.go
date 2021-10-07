package combinator

import (
	"strings"
)

func Solve(args map[string]string) []map[string]string {
	// Return early if no args were sent
	if len(args) == 0 {
		return []map[string]string{}
	}

	combos := []map[string]string{}

	// Create holder arrays to process the incoming args
	var arrays [][]string
	var replacements []string
	for key, val := range args {
		arrays = append(arrays, strings.Split(val, ","))
		replacements = append(replacements, key)
	}

	// Define mx length of each combo
	max := len(arrays) - 1

	// Define recursive function for getting combinations
	var helper func(combo map[string]string, i int)
	helper = func(combo map[string]string, i int) {
		for _, val := range arrays[i] {
			combo[replacements[i]] = val
			if i == max {
				// Append a copy of the map to the combos
				combos = append(combos, copyMap(combo))
			} else {
				helper(combo, i+1)
			}
		}
	}

	// Recurse for combinations
	helper(map[string]string{}, 0)
	return combos
}

// copyMap simple takes in a map and returns a copy of it
func copyMap(original map[string]string) map[string]string {
	copy := map[string]string{}
	for key, val := range original {
		copy[key] = val
	}
	return copy
}
