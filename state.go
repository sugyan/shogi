package shogi

import (
	"errors"
	"fmt"
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
	Piece Piece
}

// Position type
type Position struct {
	File, Rank int
}

// Pos function
func Pos(file, rank int) *Position {
	return &Position{file, rank}
}

// IsCaptured method
func (p *Position) IsCaptured() bool {
	return p.File == 0 && p.Rank == 0
}

// State definition
type State struct {
	Board    [9][9]*BoardPiece
	Captured map[Turn]*CapturedPieces
}

// Move type
type Move struct {
	Turn  Turn
	Src   *Position
	Dst   *Position
	Piece Piece
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

// MoveString method
func (s *State) MoveString(move *Move) (string, error) {
	nameMap := map[Piece]string{
		FU: "歩",
		TO: "と",
		KY: "香",
		NY: "成香",
		KE: "桂",
		NK: "成桂",
		GI: "銀",
		NG: "成銀",
		KI: "金",
		KA: "角",
		UM: "馬",
		HI: "飛",
		RY: "竜",
		OU: "玉",
	}
	result := "▲"
	if move.Turn == TurnSecond {
		result = "△"
	}
	result += fmt.Sprintf("%c%c",
		[]rune("123456789")[move.Dst.File-1],
		[]rune("一二三四五六七八九")[move.Dst.Rank-1],
	)
	if move.Src.IsCaptured() {
		result += nameMap[move.Piece]
	} else {
		bp := s.GetBoardPiece(move.Src.File, move.Src.Rank)
		if bp == nil {
			return "", fmt.Errorf("piece does not exist in (%d, %d)", move.Src.File, move.Src.Rank)
		}
		if bp.Piece != move.Piece {
			result += nameMap[bp.Piece] + "成"
		} else {
			result += nameMap[move.Piece]
		}
	}
	// TODO special case
	return result, nil
}
