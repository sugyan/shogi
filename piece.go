package shogi

type pieceCode string

// PieceCodes
const (
	FU pieceCode = "FU"
	KY pieceCode = "KY"
	KE pieceCode = "KE"
	GI pieceCode = "GI"
	KI pieceCode = "KI"
	KA pieceCode = "KA"
	HI pieceCode = "HI"
	OU pieceCode = "OU"
	TO pieceCode = "TO"
	NY pieceCode = "NY"
	NK pieceCode = "NK"
	UM pieceCode = "UM"
	RY pieceCode = "RY"
)

// Piece interface
type Piece interface {
	IsFirst() bool
	SetMove(Move)
	Code() string
}

type piece struct {
	first bool
	code  pieceCode
}

// NewPiece function
func NewPiece(move Move, code pieceCode) Piece {
	return &piece{
		first: bool(move),
		code:  code,
	}
}

func (p *piece) IsFirst() bool {
	return p.first
}

func (p *piece) SetMove(move Move) {
	p.first = bool(move)
}

func (p *piece) Code() string {
	return string(p.code)
}
