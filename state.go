package shogi

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
func (c Captured) Total() int {
	return c.FU + c.KY + c.KE + c.GI + c.KI + c.KA + c.HI
}

// State struct
type State interface {
	GetPiece(file, rank int) (Piece, error)
	SetPiece(file, rank int, piece Piece) error
	GetCaptured(Turn) Captured
	UpdateCaptured(turn Turn, fu, ky, ke, gi, ki, ka, hi int)
	Turn() Turn
	SetTurn(Turn)

	Equals(State) bool
	Clone() State

	LegalMoves() []*Move
	Move(moves ...*Move) error

	String() string
}
