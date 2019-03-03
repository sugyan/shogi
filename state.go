package shogi

import (
	"errors"
	"fmt"
	"strings"
)

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnBlack Turn = false
	TurnWhite Turn = true
)

// Captured struct
type Captured struct {
	FU int
	KY int
	KE int
	GI int
	KI int
	KA int
	HI int
}

// Total method
func (c *Captured) Total() int {
	return c.FU + c.KY + c.KE + c.GI + c.KI + c.KA + c.HI
}

// State struct
type State struct {
	Board    [9][9]Piece
	Captured [2]Captured
	Turn     Turn
}

// InitialState variable
var InitialState = &State{
	Board: [9][9]Piece{
		{WKY, WKE, WGI, WKI, WOU, WKI, WGI, WKE, WKY},
		{EMP, WHI, EMP, EMP, EMP, EMP, EMP, WKA, EMP},
		{WFU, WFU, WFU, WFU, WFU, WFU, WFU, WFU, WFU},
		{EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP},
		{EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP},
		{EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP, EMP},
		{BFU, BFU, BFU, BFU, BFU, BFU, BFU, BFU, BFU},
		{EMP, BKA, EMP, EMP, EMP, EMP, EMP, BHI, EMP},
		{BKY, BKE, BGI, BKI, BOU, BKI, BGI, BKE, BKY},
	},
}

// ErrInvalidMove error
var ErrInvalidMove = errors.New("invalid move")

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
					s.Captured[capturedIndex].FU++
				case ky:
					s.Captured[capturedIndex].KY++
				case ke:
					s.Captured[capturedIndex].KE++
				case gi:
					s.Captured[capturedIndex].GI++
				case ki:
					s.Captured[capturedIndex].KI++
				case ka:
					s.Captured[capturedIndex].KA++
				case hi:
					s.Captured[capturedIndex].HI++
				}
			}
			s.Board[move.Src.Rank-1][9-move.Src.File] = EMP
		}
		s.Board[move.Dst.Rank-1][9-move.Dst.File] = move.Piece
	}
	return nil
}

func (s *State) String() string {
	b := &strings.Builder{}
	for i := 0; i < 9; i++ {
		b.WriteString(fmt.Sprintf("P%d", i+1))
		for j := 0; j < 9; j++ {
			b.WriteString(s.Board[i][j].String())
		}
		if i < 8 {
			b.WriteRune('\n')
		}
	}
	for i := 0; i < 2; i++ {
		if s.Captured[i].Total() > 0 {
			b.WriteRune('\n')
			switch i {
			case 0:
				b.WriteString("P+")
			case 1:
				b.WriteString("P-")
			}
			for j := 0; j < s.Captured[i].HI; j++ {
				b.WriteString("00HI")
			}
			for j := 0; j < s.Captured[i].KA; j++ {
				b.WriteString("00KA")
			}
			for j := 0; j < s.Captured[i].KI; j++ {
				b.WriteString("00KI")
			}
			for j := 0; j < s.Captured[i].GI; j++ {
				b.WriteString("00GI")
			}
			for j := 0; j < s.Captured[i].KE; j++ {
				b.WriteString("00KE")
			}
			for j := 0; j < s.Captured[i].KY; j++ {
				b.WriteString("00KY")
			}
			for j := 0; j < s.Captured[i].FU; j++ {
				b.WriteString("00FU")
			}
		}
	}
	return b.String()
}
