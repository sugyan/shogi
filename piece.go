package shogi

// Piece type
type Piece uint8

const (
	white   = 0x01 << 4
	promote = 0x01 << 5
	mask    = 0x0F
)

// Piece constants
const (
	EMP Piece          = 0
	fu  Piece          = iota // 歩
	ky                        // 香
	ke                        // 桂
	gi                        // 銀
	ki                        // 金
	ka                        // 角
	hi                        // 飛
	ou                        // 玉
	to  = fu | promote        // と
	ny  = ky | promote        // 成香
	nk  = ke | promote        // 成桂
	ng  = gi | promote        // 成銀
	um  = ka | promote        // 馬
	ry  = hi | promote        // 竜
	BFU = fu
	BKY = ky
	BKE = ke
	BGI = gi
	BKI = ki
	BKA = ka
	BHI = hi
	BOU = ou
	BTO = to
	BNY = ny
	BNK = nk
	BNG = ng
	BUM = um
	BRY = ry
	WFU = BFU | white
	WKY = BKY | white
	WKE = BKE | white
	WGI = BGI | white
	WKI = BKI | white
	WKA = BKA | white
	WHI = BHI | white
	WOU = BOU | white
	WTO = BTO | white
	WNY = BNY | white
	WNK = BNK | white
	WNG = BNG | white
	WUM = BUM | white
	WRY = BRY | white
)

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
	return p&white == 0
}
