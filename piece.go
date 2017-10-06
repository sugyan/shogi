package shogi

// PieceCode type
type PieceCode string

// PieceCodes
const (
	FU = PieceCode("FU") // 歩
	KY = PieceCode("KY") // 香
	KE = PieceCode("KE") // 桂
	GI = PieceCode("GI") // 銀
	KI = PieceCode("KI") // 金
	KA = PieceCode("KA") // 角
	HI = PieceCode("HI") // 飛
	OU = PieceCode("OU") // 王, 玉
	TO = PieceCode("TO") // と
	NY = PieceCode("NY") // 成香
	NK = PieceCode("NK") // 成桂
	NG = PieceCode("NG") // 成銀
	UM = PieceCode("UM") // 馬
	RY = PieceCode("RY") // 龍
)

// Piece interface
type Piece interface {
	Code() string
	Turn() Turn
	SetTurn(Turn)
}

type piece struct {
	turn Turn
	code PieceCode
}

// NewPiece function
func NewPiece(turn Turn, code PieceCode) Piece {
	return &piece{
		turn: turn,
		code: code,
	}
}

func (p *piece) Code() string {
	return string(p.code)
}

func (p *piece) Turn() Turn {
	return p.turn
}

func (p *piece) SetTurn(turn Turn) {
	p.turn = turn
}
