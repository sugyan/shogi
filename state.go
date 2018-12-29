package shogi

// Turn type
type Turn bool

// MoveOrder constants
const (
	TurnBlack Turn = true
	TurnWhite Turn = false
)

type captured struct {
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
	board    [9][9]Piece
	captured [2]captured
}

// SetPiece method
func (s *State) SetPiece(file, rank uint8, piece Piece) {
	s.board[rank-1][9-file] = piece
}
