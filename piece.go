package shogi

// Piece type
type Piece interface {
	String() string
	Promoted() bool
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

func (p *piece) Promoted() bool {
	switch p {
	case TO, NY, NK, NG, UM, RY:
		return true
	}
	return false
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
	case FU, TO:
		cp.FU++
	case KY, NY:
		cp.KY++
	case KE, NK:
		cp.KE++
	case GI, NG:
		cp.GI++
	case KI:
		cp.KI++
	case KA, UM:
		cp.KA++
	case HI, RY:
		cp.HI++
	}
}

// SubPieces method
func (cp *CapturedPieces) SubPieces(p Piece) {
	switch p {
	case FU, TO:
		cp.FU--
	case KY, NY:
		cp.KY--
	case KE, NK:
		cp.KE--
	case GI, NG:
		cp.GI--
	case KI:
		cp.KI--
	case KA, UM:
		cp.KA--
	case HI, RY:
		cp.HI--
	}
}

// Available method
func (cp *CapturedPieces) Available() []Piece {
	results := []Piece{}
	if cp.FU > 0 {
		results = append(results, FU)
	}
	if cp.KY > 0 {
		results = append(results, KY)
	}
	if cp.KE > 0 {
		results = append(results, KE)
	}
	if cp.GI > 0 {
		results = append(results, GI)
	}
	if cp.KI > 0 {
		results = append(results, KI)
	}
	if cp.KA > 0 {
		results = append(results, KA)
	}
	if cp.HI > 0 {
		results = append(results, HI)
	}
	return results
}
