package shogi

// Player type
type Player struct {
	Namee string
}

// Record type
type Record struct {
	Players [2]*Player
	State   *State
	Moves   []*Move
}
