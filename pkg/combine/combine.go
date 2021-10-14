package combine

import (
	"github.com/jinzhu/copier"
)

// Solve takes a map[]string[]string to get the combinations
// between the values. Utilizes recursion to produce a
// []map[string]string that represents the key/value pair
// combinations.
func Solve(args map[string][]string) []map[string]string {
	// Return early if no args were sent
	if len(args) == 0 {
		return []map[string]string{}
	}

	combos := []map[string]string{}

	// Create holder arrays to process the incoming args
	var arrays [][]string
	var replacements []string
	for key, val := range args {
		arrays = append(arrays, val)
		replacements = append(replacements, key)
	}

	// Define max length of each combo
	max := len(arrays) - 1

	// Define recursive function for getting combinations
	var helper func(combo map[string]string, i int)
	helper = func(combo map[string]string, i int) {
		for _, val := range arrays[i] {
			combo[replacements[i]] = val
			if i == max {
				// Append a copy of the map to the combos
				comboCopy := map[string]string{}
				copier.Copy(&comboCopy, &combo)
				combos = append(combos, comboCopy)
			} else {
				helper(combo, i+1)
			}
		}
	}

	// Recurse to produce combinations
	helper(map[string]string{}, 0)
	return combos
}
