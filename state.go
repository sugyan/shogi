package shogi

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

// SetPiece method
func (s *State) SetPiece(file, rank int, piece Piece) {
	s.Board[rank-1][9-file] = piece
}
