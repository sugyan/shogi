package shogi

import "errors"

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

// State definition
type State struct {
	Board    [9][9]*BoardPiece
	Captured map[Turn]*CapturedPieces
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

// SetBoardPiece method
func (s *State) SetBoardPiece(file, rank int, turn Turn, piece *Piece) error {
	if file < 1 || file > 9 {
		return errors.New("invalid file")
	}
	if rank < 1 || rank > 9 {
		return errors.New("invalid rank")
	}
	s.Board[rank-1][9-file] = &BoardPiece{
		Turn:  turn,
		Piece: piece,
	}
	return nil
}
