package shogi

// PieceCode type
type PieceCode interface {
	String() string
}

type pieceCode struct {
	code string
}

// variables
var (
	FU = &pieceCode{"FU"} // 歩
	KY = &pieceCode{"KY"} // 香
	KE = &pieceCode{"KE"} // 桂
	GI = &pieceCode{"GI"} // 銀
	KI = &pieceCode{"KI"} // 金
	KA = &pieceCode{"KA"} // 角
	HI = &pieceCode{"HI"} // 飛
	OU = &pieceCode{"OU"} // 王, 玉
	TO = &pieceCode{"TO"} // と
	NY = &pieceCode{"NY"} // 成香
	NK = &pieceCode{"NK"} // 成桂
	NG = &pieceCode{"NG"} // 成銀
	UM = &pieceCode{"UM"} // 馬
	RY = &pieceCode{"RY"} // 龍
)

func (p *pieceCode) String() string {
	return p.code
}

// Piece type
type Piece struct {
	code PieceCode
}

// NewPiece function
func NewPiece(code PieceCode) *Piece {
	return &Piece{
		code: code,
	}
}

// Code methoda
func (p *Piece) Code() PieceCode {
	return p.code
}

// CapturedPieces type
type CapturedPieces struct {
	FU int
	KY int
	KE int
	GI int
	KI int
	KA int
	HI int
}

// Num method
func (cp *CapturedPieces) Num() int {
	return cp.FU + cp.KY + cp.KE + cp.GI + cp.KI + cp.KA + cp.HI
}

// AddPieces method
func (cp *CapturedPieces) AddPieces(p *Piece) {
	switch PieceCode(p.Code()) {
	case FU:
		fallthrough
	case TO:
		cp.FU++
	case KY:
		fallthrough
	case NY:
		cp.KY++
	case KE:
		fallthrough
	case NK:
		cp.KE++
	case GI:
		fallthrough
	case NG:
		cp.GI++
	case KI:
		cp.KI++
	case KA:
		fallthrough
	case UM:
		cp.KA++
	case HI:
		fallthrough
	case RY:
		cp.HI++
	}
}
