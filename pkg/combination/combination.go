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
// its args. That uses the Next() function to get each
// combination. WithSolveAhead() ensures combinations are generated
// all at once using. If not, it will solve iterativey with nextIterativeCombination()
type Stream interface {
	Next(ctx context.Context) (map[string]string, error)
}

type stream struct {
	combinations          []map[string]string
	args                  map[string][]string // the raw data from the stream
	solveAhead            bool                // if true the Next() function will solve combinations all at once using nextPreSolvedCombination()
	solved                bool
	combinationParameters []string // a list of the names of the parameters taken from the stream(args).
	positionTree          []int    // array of integers that tracks the position of the next combination within args.
	nextParameterToUpdate int      // an integer that represents the parameter we are currently looking at. Intializes as second to last value in parameters.
}

type StreamOption func(*stream)

// NewStream creates a new stream and accepts stream options for it
func NewStream(options ...StreamOption) Stream {
	cs := &stream{}
	for _, option := range options {
		option(cs)
	}
	for key := range cs.args {
		cs.combinationParameters = append(cs.combinationParameters, key)
		cs.positionTree = append(cs.positionTree, 0)
	}
	cs.nextParameterToUpdate = len(cs.combinationParameters) - 2 // Intializes nextParameterToUpdate to point to the second to last element in parameterList...
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
func WithSolveAhead(boolVar ...bool) StreamOption {
	boolVarValue := false
	if len(boolVar) > 0 {
		boolVarValue = boolVar[0]
	}
	return func(cs *stream) {
		cs.solveAhead = boolVarValue
	}
}

// updateParameterPositionsList() is a method that iterates cs.positionTree
// as well as cs.nextParameterToUpdate to be updated to the next parameter/combination in the tree
func (cs *stream) updateParameterPositionsList() {
	// Iterates through position map in reverse
	// looking for the first updatable value then breaks loop.
	// Otherwise, resets values to zero, we know to update last parameter based off i.
	var i int
	for i = len(cs.positionTree) - 1; i > cs.nextParameterToUpdate; i-- {
		if cs.positionTree[i]+1 < len(cs.args[cs.combinationParameters[i]]) {
			cs.positionTree[i]++
			break
		}
		cs.positionTree[i] = 0
	}

	// Checks to see if we need to update lastParameter based
	// off value of i.
	if i == cs.nextParameterToUpdate {
		// Checks to see if this is the last argument of the parameter.
		// Then updates parameter, and checks to see if the combination is solved.
		if cs.positionTree[cs.nextParameterToUpdate]+1 == len(cs.args[cs.combinationParameters[cs.nextParameterToUpdate]]) {
			cs.positionTree[cs.nextParameterToUpdate] = 0
			cs.nextParameterToUpdate--
		}
		// If combination is not solved, we find the next nextParameterToUpdate,
		// if nextParameterToUpdate has only 1 argument. Also, will mark as solved
		// if reach end up parameters.
		continueIterating := true
		for !cs.solved && continueIterating && cs.nextParameterToUpdate != -1 {
			if len(cs.args[cs.combinationParameters[cs.nextParameterToUpdate]]) == 1 {
				cs.nextParameterToUpdate--
				continue
			}
			continueIterating = false
			cs.positionTree[cs.nextParameterToUpdate]++
		}
	}
}

// checkSolved() accesses whether or not a combination stream is solved,
// if so it updates cs.solved variable to be true
func (cs *stream) checkSolved() {
	if cs.nextParameterToUpdate == -1 {
		cs.solved = true
	}
}

// nextIterativeCombination() gets the next combination from the stream.
// By using this, the stream will solve each combination
// iteratively.
func (cs *stream) nextIterativeCombination() (map[string]string, error) {
	if cs.solved {
		return nil, nil
	}

	// Edge case: 0 parameters
	if len(cs.combinationParameters) == 0 {
		cs.solved = true
		return nil, ErrNoArgsSet
	}

	// Edge case: 1 parameter
	if len(cs.combinationParameters) == 1 {
		comboList := map[string]string{}
		key := cs.combinationParameters[0]
		comboList[key] = cs.args[cs.combinationParameters[0]][cs.positionTree[0]]
		cs.positionTree[0]++
		if cs.positionTree[0] == len(cs.args[cs.combinationParameters[0]]) {
			cs.solved = true
		}
		return comboList, nil
	}

	// Generate a list of combinations from the positionTree
	combinationList := map[string]string{}
	for x := 0; x < len(cs.combinationParameters); x++ {
		combination := cs.args[cs.combinationParameters[x]][cs.positionTree[x]]
		parameter := cs.combinationParameters[x]
		combinationList[parameter] = combination
	}

	cs.updateParameterPositionsList()
	cs.checkSolved()
	return combinationList, nil
}

// nextPreSolvedCombination() generates all combinations at once.
// Once combinations are generated it will return the next combination.
func (cs *stream) nextPreSolvedCombination() (map[string]string, error) {
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
	combo := cs.combinations[0]
	cs.combinations = cs.combinations[1:]
	return combo, nil
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

// Next generates the combinations iteartively. If the solveAhead = true,
// the combinations are generated all at once then the next one will be returned
// within the list of cs.combinations.
func (cs *stream) Next(ctx context.Context) (map[string]string, error) {
	if cs.solveAhead {
		return cs.nextPreSolvedCombination()
	}
	return cs.nextIterativeCombination()
}
