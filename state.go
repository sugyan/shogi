package shogi

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnBlack = Turn(true)
	TurnWhite = Turn(false)
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

// NewState function
func NewState() *State {
	return &State{
		Captured: map[Turn]*CapturedPieces{
			TurnBlack: &CapturedPieces{},
			TurnWhite: &CapturedPieces{},
		},
	}
}

var pieceMap = map[string]byte{
	"FU": 0x01, "TO": 0x02,
	"KY": 0x03, "NY": 0x04,
	"KE": 0x05, "NK": 0x06,
	"GI": 0x07, "NG": 0x08,
	"KI": 0x09, "OU": 0x0A,
	"KA": 0x0B, "UM": 0x0C,
	"HI": 0x0D, "RY": 0x0E,
}

// Hash method
func (s *State) Hash() string {
	// TODO
	bytes := make([]byte, 9*9+2*7)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := s.Board[i][j]
			if b != nil {
				v := pieceMap[b.Piece.String()]
				if b.Turn == TurnWhite {
					v += 0x10
				}
				bytes[i*9+j] = v
			}
		}
	}
	for i, turn := range []Turn{TurnBlack, TurnWhite} {
		j := 9*9 + i*7
		bytes[j+0] = byte(s.Captured[turn].FU)
		bytes[j+1] = byte(s.Captured[turn].KY)
		bytes[j+2] = byte(s.Captured[turn].KE)
		bytes[j+3] = byte(s.Captured[turn].GI)
		bytes[j+4] = byte(s.Captured[turn].KI)
		bytes[j+5] = byte(s.Captured[turn].KA)
		bytes[j+6] = byte(s.Captured[turn].HI)
	}
	hash := sha1.New()
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// GetBoard method
func (s *State) GetBoard(file, rank int) *BoardPiece {
	return s.Board[rank-1][9-file]
}

// SetBoard method
func (s *State) SetBoard(file, rank int, b *BoardPiece) error {
	if file < 1 || file > 9 {
		return errors.New("invalid file")
	}
	if rank < 1 || rank > 9 {
		return errors.New("invalid rank")
	}
	if b != nil {
		s.Board[rank-1][9-file] = &BoardPiece{
			Turn:  b.Turn,
			Piece: b.Piece,
		}
	} else {
		s.Board[rank-1][9-file] = nil
	}
	return nil
}
