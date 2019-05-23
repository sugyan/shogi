package shogi

import (
	"errors"
)

// Error variables
var (
	ErrInvalidPosition = errors.New("invalid position")
	ErrInvalidMove     = errors.New("invalid move")
)
