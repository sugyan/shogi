package shogi

// Position struct
type Position struct {
	file, rank uint8
}

// Move struct
type Move struct {
	Turn  Turn
	Src   Position
	Dst   Position
	Piece Piece
}
