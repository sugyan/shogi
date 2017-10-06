package shogi

import "errors"

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnFirst  Turn = true
	TurnSecond Turn = false
)

// State definition
type State struct {
	Board    [9][9]Piece
	Captured map[Turn]*CapturedPieces
}

// CapturedPieces type
type CapturedPieces struct {
	Fu int
	Ky int
	Ke int
	Gi int
	Ki int
	Ka int
	Hi int
}

// Num method
func (cp *CapturedPieces) Num() int {
	return cp.Fu + cp.Ky + cp.Ke + cp.Gi + cp.Ki + cp.Ka + cp.Hi
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

// AddCapturedPieces method
func (s *State) AddCapturedPieces(p Piece) {
	cp := s.Captured[p.Turn()]
	switch PieceCode(p.Code()) {
	case FU:
		fallthrough
	case TO:
		cp.Fu++
	case KY:
		fallthrough
	case NY:
		cp.Ky++
	case KE:
		fallthrough
	case NK:
		cp.Ke++
	case GI:
		fallthrough
	case NG:
		cp.Gi++
	case KI:
		cp.Ki++
	case KA:
		fallthrough
	case UM:
		cp.Ka++
	case HI:
		fallthrough
	case RY:
		cp.Hi++
	}
}

// SetPiece method
func (s *State) SetPiece(file, rank int, piece Piece) error {
	if file < 1 || file > 9 {
		return errors.New("invalid file")
	}
	if rank < 1 || rank > 9 {
		return errors.New("invalid rank")
	}
	s.Board[rank-1][9-file] = piece
	return nil
}
