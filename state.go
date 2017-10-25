package shogi

import (
	"errors"
)

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnFirst  Turn = true
	TurnSecond Turn = false
)

// BoardPiece type
type BoardPiece struct {
	Turn  Turn
	Piece *Piece
}

// Position type
type Position struct {
	File, Rank int
}

// Pos function
func Pos(file, rank int) *Position {
	return &Position{file, rank}
}

// State definition
type State struct {
	Board    [9][9]*BoardPiece
	Captured map[Turn]*CapturedPieces
}

// Move type
type Move struct {
	Src   *Position
	Dst   *Position
	Piece *Piece
}

// NewState function
func NewState() *State {
	return &State{
		Captured: map[Turn]*CapturedPieces{
			TurnFirst:  &CapturedPieces{},
			TurnSecond: &CapturedPieces{},
		},
	}
}

// GetBoardPiece method
func (s *State) GetBoardPiece(file, rank int) *BoardPiece {
	return s.Board[rank-1][9-file]
}

// SetBoardPiece method
func (s *State) SetBoardPiece(file, rank int, bp *BoardPiece) error {
	if file < 1 || file > 9 {
		return errors.New("invalid file")
	}
	if rank < 1 || rank > 9 {
		return errors.New("invalid rank")
	}
	if bp != nil {
		s.Board[rank-1][9-file] = &BoardPiece{
			Turn:  bp.Turn,
			Piece: bp.Piece,
		}
	} else {
		s.Board[rank-1][9-file] = nil
	}
	return nil
}
