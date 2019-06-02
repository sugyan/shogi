package logic

import (
	"fmt"
	"strings"

	"github.com/sugyan/shogi"
)

// State struct
type State struct {
	board    [9][9]shogi.Piece
	captured [2]shogi.Captured
	turn     shogi.Turn
	Hash     uint64
}

func capturedIndex(turn shogi.Turn) int {
	if turn == shogi.TurnWhite {
		return 1
	}
	return 0
}

// NewInitialState function
func NewInitialState() *State {
	return NewState(
		[9][9]shogi.Piece{
			{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
			{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKA, shogi.EMP},
			{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU},
			{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			{shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
			{shogi.EMP, shogi.BKA, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
			{shogi.BKY, shogi.BKE, shogi.BGI, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
		},
		[2]shogi.Captured{},
		shogi.TurnBlack,
	)
}

// NewState function
func NewState(board [9][9]shogi.Piece, captured [2]shogi.Captured, turn shogi.Turn) *State {
	hash := hasher.turn[turn]
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			piece := board[i][j]
			if piece != shogi.EMP {
				hash += hasher.board[piece][i][j]
			}
		}
	}
	for _, turn := range []shogi.Turn{shogi.TurnBlack, shogi.TurnWhite} {
		c := captured[capturedIndex(turn)]
		for i := 0; i < c.FU; i++ {
			hash += hasher.captured[shogi.FU][turn]
		}
		for i := 0; i < c.KY; i++ {
			hash += hasher.captured[shogi.KY][turn]
		}
		for i := 0; i < c.KE; i++ {
			hash += hasher.captured[shogi.KE][turn]
		}
		for i := 0; i < c.GI; i++ {
			hash += hasher.captured[shogi.GI][turn]
		}
		for i := 0; i < c.KI; i++ {
			hash += hasher.captured[shogi.KI][turn]
		}
		for i := 0; i < c.KA; i++ {
			hash += hasher.captured[shogi.KA][turn]
		}
		for i := 0; i < c.HI; i++ {
			hash += hasher.captured[shogi.HI][turn]
		}
	}
	return &State{
		board:    board,
		captured: captured,
		turn:     turn,
		Hash:     hash,
	}
}

// GetPiece method for shogi.State interface
func (s *State) GetPiece(file, rank int) (shogi.Piece, error) {
	if file < 1 || file > 9 || rank < 1 || rank > 9 {
		return shogi.ERR, shogi.ErrInvalidPosition
	}
	return s.board[rank-1][9-file], nil
}

// SetPiece method for shogi.State interface
func (s *State) SetPiece(file, rank int, piece shogi.Piece) error {
	if file < 1 || file > 9 || rank < 1 || rank > 9 {
		return shogi.ErrInvalidPosition
	}
	prev := s.board[rank-1][9-file]
	s.board[rank-1][9-file] = piece
	// update hash
	if prev != shogi.EMP {
		s.Hash -= hasher.board[prev][rank-1][9-file]
	}
	s.Hash += hasher.board[piece][rank-1][9-file]
	return nil
}

// GetCaptured method for shogi.State interface
func (s *State) GetCaptured(turn shogi.Turn) shogi.Captured {
	return s.captured[capturedIndex(turn)]
}

// UpdateCaptured method for shogi.State interface
func (s *State) UpdateCaptured(turn shogi.Turn, fu, ky, ke, gi, ki, ka, hi int) {
	idx := capturedIndex(turn)
	s.captured[idx].FU += fu
	s.captured[idx].KY += ky
	s.captured[idx].KE += ke
	s.captured[idx].GI += gi
	s.captured[idx].KI += ki
	s.captured[idx].KA += ka
	s.captured[idx].HI += hi
}

// Turn method for shogi.State interface
func (s *State) Turn() shogi.Turn {
	return s.turn
}

// SetTurn method for shogi.State interface
func (s *State) SetTurn(turn shogi.Turn) {
	s.turn = turn
	// TODO: update hash
}

// Equals method for shogi.State interface
func (s *State) Equals(target shogi.State) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			piece, _ := target.GetPiece(file, rank)
			if s.board[i][j] != piece {
				return false
			}
		}
	}
	for _, turn := range []shogi.Turn{shogi.TurnBlack, shogi.TurnWhite} {
		idx := capturedIndex(turn)
		if false ||
			s.captured[idx].FU != target.GetCaptured(turn).FU ||
			s.captured[idx].KY != target.GetCaptured(turn).KY ||
			s.captured[idx].KE != target.GetCaptured(turn).KE ||
			s.captured[idx].GI != target.GetCaptured(turn).GI ||
			s.captured[idx].KI != target.GetCaptured(turn).KI ||
			s.captured[idx].KA != target.GetCaptured(turn).KA ||
			s.captured[idx].HI != target.GetCaptured(turn).HI {
			return false
		}
	}
	if s.turn != target.Turn() {
		return false
	}
	return true
}

// Clone method for shogi.State interface
func (s *State) Clone() shogi.State {
	state := *s
	return &state
}

// Move method for shogi.State interface
func (s *State) Move(moves ...*shogi.Move) error {
	for _, move := range moves {
		capturedIndex := 0
		if move.Piece.Turn() == shogi.TurnWhite {
			capturedIndex = 1
		}
		if move.Src.File == 0 && move.Src.Rank == 0 {
			// use captured piece
			switch move.Piece.Raw() {
			case shogi.FU:
				s.captured[capturedIndex].FU--
			case shogi.KY:
				s.captured[capturedIndex].KY--
			case shogi.KE:
				s.captured[capturedIndex].KE--
			case shogi.GI:
				s.captured[capturedIndex].GI--
			case shogi.KI:
				s.captured[capturedIndex].KI--
			case shogi.KA:
				s.captured[capturedIndex].KA--
			case shogi.HI:
				s.captured[capturedIndex].HI--
			}
		} else {
			// move piece
			src := s.board[move.Src.Rank-1][9-move.Src.File]
			dst := s.board[move.Dst.Rank-1][9-move.Dst.File]
			// TODO: check invalid move
			if src != move.Piece && src.Promote() != move.Piece {
				return shogi.ErrInvalidMove
			}
			if dst != shogi.EMP {
				if dst.Turn() == src.Turn() {
					return shogi.ErrInvalidMove
				}
				switch dst.Raw() {
				case shogi.FU:
					s.captured[capturedIndex].FU++
				case shogi.KY:
					s.captured[capturedIndex].KY++
				case shogi.KE:
					s.captured[capturedIndex].KE++
				case shogi.GI:
					s.captured[capturedIndex].GI++
				case shogi.KI:
					s.captured[capturedIndex].KI++
				case shogi.KA:
					s.captured[capturedIndex].KA++
				case shogi.HI:
					s.captured[capturedIndex].HI++
				}
			}
			s.board[move.Src.Rank-1][9-move.Src.File] = shogi.EMP
		}
		s.board[move.Dst.Rank-1][9-move.Dst.File] = move.Piece
		s.turn = !s.turn
	}
	return nil
}

// String method for shogi.State interface
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
