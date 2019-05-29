package shogi

import (
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

// CapturedIndex method
func (t Turn) CapturedIndex() int {
	if t == TurnBlack {
		return 0
	}
	return 1
}

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
func (c Captured) Total() int {
	return c.FU + c.KY + c.KE + c.GI + c.KI + c.KA + c.HI
}

// State struct
type State struct {
	board    [9][9]Piece
	captured [2]Captured
	Turn     Turn
	Hash     uint64
}

// NewInitialState function
func NewInitialState() *State {
	return NewState(
		[9][9]Piece{
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
		[2]Captured{},
		TurnBlack,
	)
}

// NewState function
func NewState(board [9][9]Piece, captured [2]Captured, turn Turn) *State {
	hash := hasher.turn[TurnBlack]
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			piece := board[i][j]
			if piece != EMP {
				hash += hasher.board[piece][i][j]
			}
		}
	}
	for _, turn := range []Turn{TurnBlack, TurnWhite} {
		c := captured[turn.CapturedIndex()]
		for i := 0; i < c.FU; i++ {
			hash += hasher.captured[fu][turn]
		}
		for i := 0; i < c.KY; i++ {
			hash += hasher.captured[ky][turn]
		}
		for i := 0; i < c.KE; i++ {
			hash += hasher.captured[ke][turn]
		}
		for i := 0; i < c.GI; i++ {
			hash += hasher.captured[gi][turn]
		}
		for i := 0; i < c.KI; i++ {
			hash += hasher.captured[ki][turn]
		}
		for i := 0; i < c.KA; i++ {
			hash += hasher.captured[ka][turn]
		}
		for i := 0; i < c.HI; i++ {
			hash += hasher.captured[hi][turn]
		}
	}
	return &State{
		board:    board,
		captured: captured,
		Turn:     turn,
		Hash:     hash,
	}
}

// Equals method
func (s *State) Equals(target *State) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.board[i][j] != target.board[i][j] {
				return false
			}
		}
	}
	for i := 0; i < 2; i++ {
		if false ||
			s.captured[i].FU != target.captured[i].FU ||
			s.captured[i].KY != target.captured[i].KY ||
			s.captured[i].KE != target.captured[i].KE ||
			s.captured[i].GI != target.captured[i].GI ||
			s.captured[i].KI != target.captured[i].KI ||
			s.captured[i].KA != target.captured[i].KA ||
			s.captured[i].HI != target.captured[i].HI {
			return false
		}
	}
	if s.Turn != target.Turn {
		return false
	}
	return true

}

// GetPiece method
func (s *State) GetPiece(file, rank int) (Piece, error) {
	if file < 1 || file > 9 || rank < 1 || rank > 9 {
		return ERR, ErrInvalidPosition
	}
	return s.board[rank-1][9-file], nil
}

// SetPiece method
func (s *State) SetPiece(file, rank int, piece Piece) error {
	if file < 1 || file > 9 || rank < 1 || rank > 9 {
		return ErrInvalidPosition
	}
	prev := s.board[rank-1][9-file]
	s.board[rank-1][9-file] = piece
	// update hash
	if prev != EMP {
		s.Hash -= hasher.board[prev][rank-1][9-file]
	}
	s.Hash += hasher.board[piece][rank-1][9-file]
	return nil
}

// GetCaptured method
func (s *State) GetCaptured(turn Turn) Captured {
	idx := turn.CapturedIndex()
	return s.captured[idx]
}

// UpdateCaptured method
func (s *State) UpdateCaptured(turn Turn, fu, ky, ke, gi, ki, ka, hi int) {
	idx := turn.CapturedIndex()
	s.captured[idx].FU += fu
	s.captured[idx].KY += ky
	s.captured[idx].KE += ke
	s.captured[idx].GI += gi
	s.captured[idx].KI += ki
	s.captured[idx].KA += ka
	s.captured[idx].HI += hi
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
			switch move.Piece.raw() {
			case fu:
				s.captured[capturedIndex].FU--
			case ky:
				s.captured[capturedIndex].KY--
			case ke:
				s.captured[capturedIndex].KE--
			case gi:
				s.captured[capturedIndex].GI--
			case ki:
				s.captured[capturedIndex].KI--
			case ka:
				s.captured[capturedIndex].KA--
			case hi:
				s.captured[capturedIndex].HI--
			}
		} else {
			// move piece
			src := s.board[move.Src.Rank-1][9-move.Src.File]
			dst := s.board[move.Dst.Rank-1][9-move.Dst.File]
			// TODO: check invalid move
			if src != move.Piece && src.Promote() != move.Piece {
				return ErrInvalidMove
			}
			if dst != EMP {
				if dst.Turn() == src.Turn() {
					return ErrInvalidMove
				}
				switch dst.raw() {
				case fu:
					s.captured[capturedIndex].FU++
				case ky:
					s.captured[capturedIndex].KY++
				case ke:
					s.captured[capturedIndex].KE++
				case gi:
					s.captured[capturedIndex].GI++
				case ki:
					s.captured[capturedIndex].KI++
				case ka:
					s.captured[capturedIndex].KA++
				case hi:
					s.captured[capturedIndex].HI++
				}
			}
			s.board[move.Src.Rank-1][9-move.Src.File] = EMP
		}
		s.board[move.Dst.Rank-1][9-move.Dst.File] = move.Piece
		s.Turn = !s.Turn
	}
	return nil
}

func (s *State) String() string {
	b := &strings.Builder{}
	for i := 0; i < 9; i++ {
		b.WriteString(fmt.Sprintf("P%d", i+1))
		for j := 0; j < 9; j++ {
			b.WriteString(s.board[i][j].String())
		}
		if i < 8 {
			b.WriteRune('\n')
		}
	}
	for i := 0; i < 2; i++ {
		if s.captured[i].Total() > 0 {
			b.WriteRune('\n')
			switch i {
			case 0:
				b.WriteString("P+")
			case 1:
				b.WriteString("P-")
			}
			for j := 0; j < s.captured[i].HI; j++ {
				b.WriteString("00HI")
			}
			for j := 0; j < s.captured[i].KA; j++ {
				b.WriteString("00KA")
			}
			for j := 0; j < s.captured[i].KI; j++ {
				b.WriteString("00KI")
			}
			for j := 0; j < s.captured[i].GI; j++ {
				b.WriteString("00GI")
			}
			for j := 0; j < s.captured[i].KE; j++ {
				b.WriteString("00KE")
			}
			for j := 0; j < s.captured[i].KY; j++ {
				b.WriteString("00KY")
			}
			for j := 0; j < s.captured[i].FU; j++ {
				b.WriteString("00FU")
			}
		}
	}
	return b.String()
}
