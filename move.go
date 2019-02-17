package shogi

// Position struct
type Position struct {
	File, Rank uint8
}

// Move struct
type Move struct {
	Src   Position
	Dst   Position
	Piece Piece
}
