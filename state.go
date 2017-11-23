package shogi

import (
	"crypto/sha1"
	"encoding/hex"
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
func Pos(file, rank int) Position {
	return Position{file, rank}
}

// IsCaptured method
func (p *Position) IsCaptured() bool {
	return p.File == 0 && p.Rank == 0
}

// State definition
type State struct {
	Board      [9][9]*BoardPiece
	Captured   map[Turn]*CapturedPieces
	latestMove *Move
}

// Move type
type Move struct {
	Turn  Turn
	Src   Position
	Dst   Position
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

// Hash method
func (s *State) Hash() string {
	// TODO
	pieceByteMap := map[Piece]byte{
		FU: 1, TO: 2,
		KY: 3, NY: 4,
		KE: 5, NK: 6,
		GI: 7, NG: 8,
		KI: 9, OU: 10,
		KA: 11, UM: 12,
		HI: 13, RY: 14,
	}
	turnByteMap := map[Turn]byte{
		TurnFirst:  1,
		TurnSecond: 0,
	}
	hash := sha1.New()
	for i := 0; i < 9; i++ {
		bytes := make([]byte, 0, 18)
		for j := 0; j < 9; j++ {
			bp := s.Board[i][j]
			if bp != nil {
				bytes = append(bytes, []byte{pieceByteMap[bp.Piece], turnByteMap[bp.Turn]}...)
			} else {
				bytes = append(bytes, []byte{0, 0}...)
			}
		}
		hash.Write(bytes)
	}
	for _, turn := range []Turn{TurnFirst, TurnSecond} {
		bytes := make([]byte, 7)
		bytes[0] = byte(s.Captured[turn].FU)
		bytes[1] = byte(s.Captured[turn].KY)
		bytes[2] = byte(s.Captured[turn].KE)
		bytes[3] = byte(s.Captured[turn].GI)
		bytes[4] = byte(s.Captured[turn].KI)
		bytes[5] = byte(s.Captured[turn].KA)
		bytes[6] = byte(s.Captured[turn].HI)
		hash.Write(bytes)
	}
	return hex.EncodeToString(hash.Sum(nil))
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

// Apply method
func (s *State) Apply(move *Move) {
	// update state
	s.latestMove = move
	if move.Src.File > 0 && move.Src.Rank > 0 {
		bp := s.GetBoardPiece(move.Dst.File, move.Dst.Rank)
		if bp != nil {
			s.Captured[move.Turn].AddPieces(bp.Piece)
		}
		s.SetBoardPiece(move.Src.File, move.Src.Rank, nil)
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
	} else {
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
		s.Captured[move.Turn].SubPieces(move.Piece)
	}
}

// MoveStrings method
func (s *State) MoveStrings(moves []*Move) ([]string, error) {
	results := []string{}
	state := s.Clone()
	for _, move := range moves {
		ms, err := state.MoveString(move)
		if err != nil {
			return nil, err
		}
		state.Apply(move)
		results = append(results, ms)
	}
	return results, nil
}

// MoveString method
func (s *State) MoveString(move *Move) (string, error) {
	// move string
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
	if s.latestMove != nil && move.Dst == s.latestMove.Dst {
		result += "同"
	} else {
		result += fmt.Sprintf("%c%c",
			[]rune("123456789")[move.Dst.File-1],
			[]rune("一二三四五六七八九")[move.Dst.Rank-1],
		)
	}
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
