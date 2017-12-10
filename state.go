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
		TurnBlack: 1,
		TurnWhite: 0,
	}
	hash := sha1.New()
	for i := 0; i < 9; i++ {
		bytes := make([]byte, 0, 18)
		for j := 0; j < 9; j++ {
			b := s.Board[i][j]
			if b != nil {
				bytes = append(bytes, []byte{pieceByteMap[b.Piece], turnByteMap[b.Turn]}...)
			} else {
				bytes = append(bytes, []byte{0, 0}...)
			}
		}
		hash.Write(bytes)
	}
	for _, turn := range []Turn{TurnBlack, TurnWhite} {
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
