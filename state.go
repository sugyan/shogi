package shogi

import (
	"errors"
)

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnBlack Turn = true
	TurnWhite Turn = false
)

// Captured struct
type Captured struct {
	FU uint8
	KY uint8
	KE uint8
	GI uint8
	KI uint8
	KA uint8
	HI uint8
}

// State struct
type State struct {
	Board    [9][9]Piece
	Captured [2]Captured
}

// ErrInvalidMove error
var ErrInvalidMove = errors.New("invalid move")

// SetPiece method
func (s *State) SetPiece(file, rank int, piece Piece) {
	s.Board[rank-1][9-file] = piece
}

// Move method
func (s *State) Move(moves ...*Move) error {
	for _, move := range moves {
		capturedIndex := 0
		if move.Piece.Turn() == TurnWhite {
			capturedIndex = 1
		}
		if move.Src.File == 0 && move.Src.Rank == 0 {
			// use captured piece
			switch move.Piece & mask {
			case fu:
				s.Captured[capturedIndex].FU--
			case ky:
				s.Captured[capturedIndex].KY--
			case ke:
				s.Captured[capturedIndex].KE--
			case gi:
				s.Captured[capturedIndex].GI--
			case ki:
				s.Captured[capturedIndex].KI--
			case ka:
				s.Captured[capturedIndex].KA--
			case hi:
				s.Captured[capturedIndex].HI--
			}
		} else {
			// move piece
			src := s.Board[move.Src.Rank-1][9-move.Src.File]
			dst := s.Board[move.Dst.Rank-1][9-move.Dst.File]
			// TODO: check invalid move
			if src != move.Piece && src.Promote() != move.Piece {
				return ErrInvalidMove
			}
			if dst != EMP {
				if dst.Turn() == src.Turn() {
					return ErrInvalidMove
				}
				switch dst & mask {
				case fu:
					s.Captured[1-capturedIndex].FU++
				case ky:
					s.Captured[1-capturedIndex].KY++
				case ke:
					s.Captured[1-capturedIndex].KE++
				case gi:
					s.Captured[1-capturedIndex].GI++
				case ki:
					s.Captured[1-capturedIndex].KI++
				case ka:
					s.Captured[1-capturedIndex].KA++
				case hi:
					s.Captured[1-capturedIndex].HI++
				}
			}
			s.Board[move.Src.Rank-1][9-move.Src.File] = EMP
		}
		s.Board[move.Dst.Rank-1][9-move.Dst.File] = move.Piece
	}
	return nil
}
