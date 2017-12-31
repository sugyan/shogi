package dfpn

import "github.com/sugyan/shogi"

// Result type
type Result int

// Resulst constants
const (
	ResultU = iota // Unknown
	ResultT        // True
	ResultF        // False
)

// Node type
type Node struct {
	Result   Result
	Move     *shogi.Move
	State    *shogi.State
	Children []*Node
	expanded bool
	pn, dn   uint32
}

func (n *Node) getP() uint32 {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		return n.dn
	case shogi.TurnWhite:
		return n.pn
	}
	return 0
}

func (n *Node) setP(v uint32) {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		n.dn = v
	case shogi.TurnWhite:
		n.pn = v
	}
}

func (n *Node) getD() uint32 {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		return n.pn
	case shogi.TurnWhite:
		return n.dn
	}
	return 0
}

func (n *Node) setD(v uint32) {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		n.pn = v
	case shogi.TurnWhite:
		n.dn = v
	}
}
