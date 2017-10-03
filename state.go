package shogi

// Move type
type Move bool

// MoveOrder constants
const (
	MoveFirst  Move = true
	MoveSecond Move = false
)

// State definition
type State struct {
	Board [9][9]Piece
	Hands map[Move][]Piece
}

// NewState function
func NewState() *State {
	return &State{
		Hands: make(map[Move][]Piece),
	}
}

// AddHandPieces method
func (s *State) AddHandPieces(p Piece) {
	if p.IsFirst() {
		s.Hands[MoveFirst] = append(s.Hands[MoveFirst], p)
	} else {
		s.Hands[MoveSecond] = append(s.Hands[MoveSecond], p)
	}
}
