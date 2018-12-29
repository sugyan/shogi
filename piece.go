package shogi

// Piece type
type Piece uint8

const w = 0x01 << 4
const p = 0x01 << 5

// Piece constants
const (
	BLANK Piece    = 0
	fu    Piece    = iota // 歩
	ky                    // 香
	ke                    // 桂
	gi                    // 銀
	ki                    // 金
	ka                    // 角
	hi                    // 飛
	ou                    // 玉
	to    = fu | p        // と
	ny    = ky | p        // 成香
	nk    = ke | p        // 成桂
	ng    = gi | p        // 成銀
	um    = ka | p        // 馬
	ry    = hi | p        // 竜
	BFU   = fu
	BKY   = ky
	BKE   = ke
	BGI   = gi
	BKI   = ki
	BKA   = ka
	BHI   = hi
	BOU   = ou
	BTO   = to
	BNY   = ny
	BNK   = nk
	BNG   = ng
	BUM   = um
	BRY   = ry
	WFU   = BFU | w
	WKY   = BKY | w
	WKE   = BKE | w
	WGI   = BGI | w
	WKI   = BKI | w
	WKA   = BKA | w
	WHI   = BHI | w
	WOU   = BOU | w
	WTO   = BTO | w
	WNY   = BNY | w
	WNK   = BNK | w
	WNG   = BNG | w
	WUM   = BUM | w
	WRY   = BRY | w
)

var pieceStringMap = map[Piece]string{
	BLANK: " * ",
	BFU:   "+FU",
	BKY:   "+KY",
	BKE:   "+KE",
	BGI:   "+GI",
	BKI:   "+KI",
	BKA:   "+KA",
	BHI:   "+HI",
	BOU:   "+OU",
	BTO:   "+TO",
	BNY:   "+NY",
	BNK:   "+NK",
	BNG:   "+NG",
	BUM:   "+UM",
	BRY:   "+RY",
	WFU:   "-FU",
	WKY:   "-KY",
	WKE:   "-KE",
	WGI:   "-GI",
	WKI:   "-KI",
	WKA:   "-KA",
	WHI:   "-HI",
	WOU:   "-OU",
	WTO:   "-TO",
	WNY:   "-NY",
	WNK:   "-NK",
	WNG:   "-NG",
	WUM:   "-UM",
	WRY:   "-RY",
}

// String method
func (piece Piece) String() string {
	if s, exist := pieceStringMap[piece]; exist {
		return s
	}
	return ""
}

// IsPromoted method
func (piece Piece) IsPromoted() bool {
	return piece&p != 0
}
