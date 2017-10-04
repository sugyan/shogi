package shogi

type pieceCode string

// PieceCodes
const (
	FU = pieceCode("FU") // 歩
	KY = pieceCode("KY") // 香
	KE = pieceCode("KE") // 桂
	GI = pieceCode("GI") // 銀
	KI = pieceCode("KI") // 金
	KA = pieceCode("KA") // 角
	HI = pieceCode("HI") // 飛
	OU = pieceCode("OU") // 王, 玉
	TO = pieceCode("TO") // と
	NY = pieceCode("NY") // 成香
	NK = pieceCode("NK") // 成桂
	NG = pieceCode("NG") // 成銀
	UM = pieceCode("UM") // 馬
	RY = pieceCode("RY") // 龍
)

// Piece interface
type Piece interface {
	Code() string
	Turn() Turn
	SetTurn(Turn)
}

type piece struct {
	turn Turn
	code pieceCode
}

// NewPiece function
func NewPiece(turn Turn, code pieceCode) Piece {
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
