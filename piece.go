package shogi

type (
	// RawPiece type
	RawPiece uint8
	// Piece types
	Piece uint8
)

const (
	white   = 0x01 << 4
	promote = 0x01 << 5
	mask    = 0x0F
)

// constant variables
const (
	FU RawPiece = iota + 1 // 歩
	KY                     // 香
	KE                     // 桂
	GI                     // 銀
	KI                     // 金
	KA                     // 角
	HI                     // 飛
	OU                     // 玉

	TO = FU | promote // と
	NY = KY | promote // 成香
	NK = KE | promote // 成桂
	NG = GI | promote // 成銀
	UM = KA | promote // 馬
	RY = HI | promote // 竜
)

// Piece constants
const (
	EMP Piece = 0
	BFU       = Piece(FU)
	BKY       = Piece(KY)
	BKE       = Piece(KE)
	BGI       = Piece(GI)
	BKI       = Piece(KI)
	BKA       = Piece(KA)
	BHI       = Piece(HI)
	BOU       = Piece(OU)
	BTO       = Piece(TO)
	BNY       = Piece(NY)
	BNK       = Piece(NK)
	BNG       = Piece(NG)
	BUM       = Piece(UM)
	BRY       = Piece(RY)
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

	ERR = 0xFF
)

// MakePiece function
func MakePiece(p RawPiece, turn Turn) Piece {
	piece := Piece(p)
	if turn == TurnWhite {
		piece |= white
	}
	return piece
}

// PieceStringMap variable
var PieceStringMap = map[Piece]string{
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

// Raw method
func (p Piece) Raw() RawPiece {
	return RawPiece(p & mask)
}

// String method
func (p Piece) String() string {
	if s, exist := PieceStringMap[p]; exist {
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
