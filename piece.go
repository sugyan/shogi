package shogi

// PieceType definition
type PieceType string

// PieceTypes
const (
	FU PieceType = "FU"
	KY PieceType = "KY"
	KE PieceType = "KE"
	GI PieceType = "GI"
	KI PieceType = "KI"
	KA PieceType = "KA"
	HI PieceType = "HI"
	OU PieceType = "OU"
	TO PieceType = "TO"
	NY PieceType = "NY"
	NK PieceType = "NK"
	UM PieceType = "UM"
	RY PieceType = "RY"
)

// Piece type
type Piece struct {
	First bool
	Type  PieceType
}
