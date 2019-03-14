package shogi

type rawPiece uint8

// Piece type
type Piece uint8

const (
	white   = 0x01 << 4
	promote = 0x01 << 5
	mask    = 0x0F
)

const (
	fu rawPiece       = iota + 1 // 歩
	ky                           // 香
	ke                           // 桂
	gi                           // 銀
	ki                           // 金
	ka                           // 角
	hi                           // 飛
	ou                           // 玉
	to = fu | promote            // と
	ny = ky | promote            // 成香
	nk = ke | promote            // 成桂
	ng = gi | promote            // 成銀
	um = ka | promote            // 馬
	ry = hi | promote            // 竜
)

// Piece constants
const (
	EMP Piece = 0
	BFU       = Piece(fu)
	BKY       = Piece(ky)
	BKE       = Piece(ke)
	BGI       = Piece(gi)
	BKI       = Piece(ki)
	BKA       = Piece(ka)
	BHI       = Piece(hi)
	BOU       = Piece(ou)
	BTO       = Piece(to)
	BNY       = Piece(ny)
	BNK       = Piece(nk)
	BNG       = Piece(ng)
	BUM       = Piece(um)
	BRY       = Piece(ry)
	WFU       = BFU | white
	WKY       = BKY | white
	WKE       = BKE | white
	WGI       = BGI | white
	WKI       = BKI | white
	WKA       = BKA | white
	WHI       = BHI | white
	WOU       = BOU | white
	WTO       = BTO | white
	WNY       = BNY | white
	WNK       = BNK | white
	WNG       = BNG | white
	WUM       = BUM | white
	WRY       = BRY | white
)

func makePiece(p rawPiece, turn Turn) Piece {
	piece := Piece(p)
	if turn == TurnWhite {
		piece |= white
	}
	return piece
}

var pieceStringMap = map[Piece]string{
	EMP: " * ",
	BFU: "+FU",
	BKY: "+KY",
	BKE: "+KE",
	BGI: "+GI",
	BKI: "+KI",
	BKA: "+KA",
	BHI: "+HI",
	BOU: "+OU",
	BTO: "+TO",
	BNY: "+NY",
	BNK: "+NK",
	BNG: "+NG",
	BUM: "+UM",
	BRY: "+RY",
	WFU: "-FU",
	WKY: "-KY",
	WKE: "-KE",
	WGI: "-GI",
	WKI: "-KI",
	WKA: "-KA",
	WHI: "-HI",
	WOU: "-OU",
	WTO: "-TO",
	WNY: "-NY",
	WNK: "-NK",
	WNG: "-NG",
	WUM: "-UM",
	WRY: "-RY",
}

func (p Piece) raw() rawPiece {
	return rawPiece(p & mask)
}

// String method
func (p Piece) String() string {
	if s, exist := pieceStringMap[p]; exist {
		return s
	}
	return ""
}

// IsPromoted method
func (p Piece) IsPromoted() bool {
	return p&promote != 0
}

// Promote method
func (p Piece) Promote() Piece {
	return p | promote
}

// Turn method
func (p Piece) Turn() Turn {
	return p&white != 0
}
