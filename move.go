package shogi

// Position struct
type Position struct {
	File, Rank int
}

// Move struct
type Move struct {
	Src   Position
	Dst   Position
	Piece Piece
}

// NewMove function
func NewMove(srcFile, srcRank, dstFile, dstRank int, piece Piece) *Move {
	return &Move{
		Position{srcFile, srcRank},
		Position{dstFile, dstRank},
		piece,
	}
}
