package shogi

// Piece type
type Piece interface {
	String() string
}

type piece struct {
	code string
}

// variables
var (
	FU = &piece{"FU"} // 歩
	KY = &piece{"KY"} // 香
	KE = &piece{"KE"} // 桂
	GI = &piece{"GI"} // 銀
	KI = &piece{"KI"} // 金
	KA = &piece{"KA"} // 角
	HI = &piece{"HI"} // 飛
	OU = &piece{"OU"} // 王, 玉
	TO = &piece{"TO"} // と
	NY = &piece{"NY"} // 成香
	NK = &piece{"NK"} // 成桂
	NG = &piece{"NG"} // 成銀
	UM = &piece{"UM"} // 馬
	RY = &piece{"RY"} // 龍
)

var promoteMap = map[Piece]Piece{
	FU: TO,
	KY: NY,
	KE: NK,
	GI: NG,
	KA: UM,
	HI: RY,
}

func (p *piece) String() string {
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
func (cp *CapturedPieces) AddPieces(p Piece) {
	switch p {
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

// SubPieces method
func (cp *CapturedPieces) SubPieces(p Piece) {
	switch p {
	case FU:
		cp.FU--
	case KY:
		cp.KY--
	case KE:
		cp.KE--
	case GI:
		cp.GI--
	case KI:
		cp.KI--
	case KA:
		cp.KA--
	case HI:
		cp.HI--
	}
}
