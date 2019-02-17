package shogi

// Player type
type Player struct {
	Name string
}

// Record type
type Record struct {
	Players [2]*Player
	State   *State
	Moves   []*Move
}
