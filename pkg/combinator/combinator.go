package combinator

import (
	"github.com/operator-framework/combo/pkg/types"
)

func Solve(args types.ComboArgs) types.Combos {
	combos := types.Combos{}

	// Create holder arrays to process the incoming args
	var arrays [][]string
	var replacements []string
	for _, arg := range args {
		arrays = append(arrays, arg.Options)
		replacements = append(replacements, arg.Name)
	}

	// Define mx length of each combo
	max := len(arrays) - 1

	// Define recrusive function for getting combinations
	var helper func(combo map[string]string, i int)
	helper = func(combo map[string]string, i int) {
		for j, l := 0, len(arrays[i]); j < l; j++ {
			combo[replacements[i]] = arrays[i][j]
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
