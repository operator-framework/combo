package combination

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
)

// Specify which errors this package can return
var (
	ErrNoArgsSet             = errors.New("args not set")
	ErrCombinationsNotSolved = errors.New("combinations not yet solved")
)

// Stream is a representation of all possible combinations
// its args. Can use either the Next() function to get each
// combination one at a time or All() to get them all.
// WithSolveAhead() ensures that the combinations are generated
// before each call to Next() or All() but is only run once.
type Stream interface {
	Next(ctx context.Context) (map[string]string, error)
	All() ([]map[string]string, error)
}

// args: the raw data from the stream
// parameterListFromArgs: a list of the names of the parameters
//                        taken from the stream(args).
// positionMapInArgs: array of integers that tracks the position
//                    of the next combination within args.
// lastUpdatedParameter: an integer that represents the parameter
//                       we are currently looking at. Intializes as
//                       second to last value in parameters.
type stream struct {
	combinations          []map[string]string
	args                  map[string][]string
	solveAhead            bool
	solved                bool
	parameterListFromArgs []string
	positionsMapInArgs    []int
	lastUpdatedParameter  int
}

type StreamOption func(*stream)

// NewStream creates a new stream and accepts stream options for it
func NewStream(options ...StreamOption) Stream {
	cs := &stream{}
	for _, option := range options {
		option(cs)
	}
	for key := range cs.args {
		cs.parameterListFromArgs = append(cs.parameterListFromArgs, key)
		cs.positionsMapInArgs = append(cs.positionsMapInArgs, 0)
	}
	cs.lastUpdatedParameter = len(cs.parameterListFromArgs) - 2
	return cs
}

// WithArgs specifies which args to utilize in the new stream
func WithArgs(args map[string][]string) StreamOption {
	return func(cs *stream) {
		cs.args = args
	}
}

// WithSolveAhead specifies whether to solve before calling Next or All,
// only occurs on the first call to Next or All. By using this, the Stream
// will solve all possible combinations of its args which could take a lot
// of computation given a large enough input.
func WithSolveAhead() StreamOption {
	return func(cs *stream) {
		cs.solveAhead = true
	}
}

// Next gets the next combination from the stream.
// By using this, the stream will solve each combination
// iteratively.
func (cs *stream) Next(ctx context.Context) (map[string]string, error) {
	// Check to see if anymore combinations exist
	if !cs.solved {
		// Edge case: No parameterListFromArgs
		if len(cs.parameterListFromArgs) == 0 {
			cs.solved = true
			return nil, ErrNoArgsSet
		}

		// comboList is a variable that holds a list of combinations to be returned.
		comboList := map[string]string{}

		// Generate the list of combinations based off current positions
		for x := 0; x < len(cs.parameterListFromArgs); x++ {
			var combo string = cs.args[cs.parameterListFromArgs[x]][cs.positionsMapInArgs[x]]
			var key string = cs.parameterListFromArgs[x]
			comboList[key] = combo
		}

		// Iterates through position map in reverse
		// looking for the first updatable value then breaks loop.
		// Otherwise, resets values to zero, we know to update last parameter based off i.
		var i int
		for i = len(cs.positionsMapInArgs) - 1; i > cs.lastUpdatedParameter; i-- {
			if cs.positionsMapInArgs[i]+1 < len(cs.args[cs.parameterListFromArgs[i]]) {
				cs.positionsMapInArgs[i]++
				break
			} else {
				cs.positionsMapInArgs[i] = 0
			}
		}

		// Checks to see if we need to update lastParameter based
		// off value of i.
		if i == cs.lastUpdatedParameter {
			// Checks to see if this is the last argument of the parameter.
			// Then updates parameter, and checks to see if the combination is solved.
			if cs.positionsMapInArgs[cs.lastUpdatedParameter]+1 == len(cs.args[cs.parameterListFromArgs[cs.lastUpdatedParameter]]) {
				cs.positionsMapInArgs[cs.lastUpdatedParameter] = 0
				cs.lastUpdatedParameter--
				if cs.lastUpdatedParameter == -1 {
					cs.solved = true
				}
			}
			// If combination is not solved, we find the next lastUpdatedParameter,
			// if lastUpdatedParameter has only 1 argument. Also, will mark as solved
			// if reach end up parameters.
			if !cs.solved {
				runner := true
				for runner {
					length := len(cs.args[cs.parameterListFromArgs[cs.lastUpdatedParameter]])
					if length == 1 {
						cs.lastUpdatedParameter--
					} else {
						runner = false
						cs.positionsMapInArgs[cs.lastUpdatedParameter]++
					}
					if cs.lastUpdatedParameter == -1 {
						cs.solved = true
						runner = false
					}
				}
			}
		}

		return comboList, nil
	}
	return nil, nil
}

// All retrieves all of the combinations from the stream
func (cs *stream) All() ([]map[string]string, error) {
	if cs.solveAhead && !cs.solved {
		if err := cs.solve(); err != nil {
			return nil, err
		}
	}

	if len(cs.combinations) == 0 {
		if cs.solved {
			return nil, nil
		}
		return nil, ErrCombinationsNotSolved
	}

	return cs.combinations, nil
}

// solve takes the current stream and its args to solve their combinations
func (cs *stream) solve() error {
	// Return early if no args were sent
	if len(cs.args) == 0 {
		return ErrNoArgsSet
	}

	combos := []map[string]string{}

	// Create holder arrays to process the incoming args
	arrays := [][]string{}
	replacements := []string{}
	for key, val := range cs.args {
		arrays = append(arrays, val)
		replacements = append(replacements, key)
	}

	// Define max length of each combo
	max := len(arrays) - 1

	// Define recursive function for getting combinations
	var err error
	var recurse func(combo map[string]string, i int)
	recurse = func(combo map[string]string, i int) {
		for _, val := range arrays[i] {
			combo[replacements[i]] = val
			if i == max {
				// Append a copy of the map to the combos
				comboCopy := map[string]string{}
				err = copier.Copy(&comboCopy, &combo)
				combos = append(combos, comboCopy)
			} else {
				recurse(combo, i+1)
			}
		}
	}

	// Recurse to produce combinations
	recurse(map[string]string{}, 0)
	if err != nil {
		return err
	}

	cs.combinations = combos
	cs.solved = true

	return nil
}
