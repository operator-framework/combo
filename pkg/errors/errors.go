package errors

import "errors"

// Define errors shared across Combo
var (
	ErrCouldNotReadFile = errors.New("could not read file")
)
