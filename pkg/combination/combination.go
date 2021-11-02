package combination

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
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
type streamImp struct {
	combinations []map[string]string
	args         map[string][]string
	solveAhead   bool
	solved       bool
}

// Specify which errors this package can return
var (
	ErrNoArgsSet             = errors.New("args not set")
	ErrCombinationsNotSolved = errors.New("combinations not yet solved")
)

// NewStream creates a new stream and accepts stream options for it
type streamOption func(*streamImp)

func NewStream(options ...streamOption) Stream {
	cs := &streamImp{}
	for _, option := range options {
		option(cs)
	}
	return cs
}

// WithArgs specifies which args to utilize in the new stream
func WithArgs(args map[string][]string) streamOption {
	return func(cs *streamImp) {
		cs.args = args
	}
}

// WithSolveAhead specifies whether to solve before calling Next or All,
// only occurs on the first call to Next or All. By using this, the Stream
// will solve all possible combinations of its args which could take a lot
// of computation given a large enough input.
func WithSolveAhead() streamOption {
	return func(cs *streamImp) {
		cs.solveAhead = true
	}
}

// Next gets the next combination from the stream.
// TODO: Currently this does not solve each combination iteratively and will
//       need to do so in the future to ensure an optimal use of memory. This
//       is currently more of a stub to allow consumer packages to maintain its
// 		 interface.
func (cs *streamImp) Next(ctx context.Context) (map[string]string, error) {
	if cs.solveAhead && !cs.solved {
		if err := cs.solve(); err != nil {
			return nil, err
		}
	}

	if cs.solveAhead && !cs.solved {
		cs.solve()
	}

	if len(cs.combinations) == 0 {
		if cs.solved {
			return nil, nil
		}
		return nil, ErrCombinationsNotSolved
	}

	next := cs.combinations[0]
	cs.combinations = cs.combinations[1:]

	return next, nil
}

// All retrieves all of the combinations from the stream
func (cs *streamImp) All() ([]map[string]string, error) {
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
func (cs *streamImp) solve() error {
	// Return early if no args were sent
	if len(cs.args) == 0 {
		return ErrNoArgsSet
	}

	combos := []map[string]string{}

	// Create holder arrays to process the incoming args
	var arrays [][]string
	var replacements []string
	for key, val := range cs.args {
		arrays = append(arrays, val)
		replacements = append(replacements, key)
	}

	// Define max length of each combo
	max := len(arrays) - 1

	// Define recursive function for getting combinations
	var recurse func(combo map[string]string, i int)
	recurse = func(combo map[string]string, i int) {
		for _, val := range arrays[i] {
			combo[replacements[i]] = val
			if i == max {
				// Append a copy of the map to the combos
				comboCopy := map[string]string{}
				copier.Copy(&comboCopy, &combo)
				combos = append(combos, comboCopy)
			} else {
				recurse(combo, i+1)
			}
		}
	}

	// Recurse to produce combinations
	recurse(map[string]string{}, 0)

	cs.combinations = combos
	cs.solved = true

	return nil
}
