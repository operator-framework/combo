package combination

import (
	"context"
	"errors"
	"sync"

	"github.com/jinzhu/copier"
)

// Set is a representation of one possible combination of
// key/value pairs.
type Set map[string]string

// Stream is a representation of all possible combinations
// its args. Can use either the Next() function to get each
// combination one at a time or All() to get them all.
// WithSolveAhead() ensures that the combinations are generated
// before each call to Next() or All() but is only run once.
type Stream struct {
	combinations []Set
	args         map[string][]string
	current      int
	solveAhead   bool
	solveOnce    sync.Once
}

// streamOption is an option function that can be used within
// a stream.
type streamOption func(*Stream)

// Specify which errors this package can return
var (
	ErrNoArgsSet             = errors.New("args not set")
	ErrCombinationsNotSolved = errors.New("combinations not yet solved")
)

// NewStream creates a new stream and accepts stream options for it
func NewStream(options ...streamOption) Stream {
	cs := &Stream{current: -1}
	for _, option := range options {
		option(cs)
	}
	return *cs
}

// WithArgs specifies which args to utilize in the new stream
func WithArgs(args map[string][]string) streamOption {
	return func(cs *Stream) {
		cs.args = args
	}
}

// WithSolveAhead specifies whether to solve before calling Next or All,
// only occurs on the first call to Next or All. By using this, the Stream
// will solve all possible combinations of its args which could take a lot
// of computation given a large enough input.
func WithSolveAhead() streamOption {
	return func(cs *Stream) {
		cs.solveAhead = true
	}
}

// Next gets the next combination from the stream.
// TODO: Currently this does not solve each combination iteratively and will
//       need to do so in the future to ensure an optimal use of memory. This
//       is currently more of a stub to allow consumer packages to maintain its
// 		 interface.
func (cs *Stream) Next(ctx context.Context) (Set, error) {
	if cs.solveAhead {
		cs.solveOnce.Do(cs.solve)
	}

	if len(cs.combinations) == 0 {
		ctx.Err()
		err := ErrCombinationsNotSolved
		if len(cs.args) == 0 {
			err = ErrNoArgsSet
		}
		return nil, err
	}

	if cs.current == len(cs.combinations)-1 {
		return nil, nil
	}

	cs.current++

	return cs.combinations[cs.current], nil
}

// All retrieves all of the combinations from the stream
func (cs *Stream) All() ([]Set, error) {
	if cs.solveAhead {
		cs.solveOnce.Do(cs.solve)
	}

	if len(cs.combinations) == 0 {
		err := ErrCombinationsNotSolved
		if len(cs.args) == 0 {
			err = ErrNoArgsSet
		}
		return nil, err
	}

	return cs.combinations, nil
}

// solve takes the current stream and its args to solve their combinations
func (cs *Stream) solve() {
	// Return early if no args were sent
	if len(cs.args) == 0 {
		return
	}

	combos := []Set{}

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
	cs.combinations = combos
}
