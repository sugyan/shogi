package dfpn

import (
	"github.com/sugyan/shogi"
)

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
	pn       uint32
	dn       uint32
	Move     *shogi.Move
	state    *shogi.State
	parent   *Node
	Children []*Node
}

func (n *Node) getPhi() uint32 {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		return n.dn
	case shogi.TurnWhite:
		return n.pn
	}
	return 0
}

func (n *Node) setPhi(v uint32) {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		n.dn = v
	case shogi.TurnWhite:
		n.pn = v
	}
}

func (n *Node) getDelta() uint32 {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		return n.pn
	case shogi.TurnWhite:
		return n.dn
	}
	return 0
}

func (n *Node) setDelta(v uint32) {
	switch n.Move.Turn {
	case shogi.TurnBlack:
		n.pn = v
	case shogi.TurnWhite:
		n.dn = v
	}
}

func (n *Node) setResult(result Result) Result {
	n.Result = result
	if (n.Move.Turn == shogi.TurnBlack && n.Result == ResultT) ||
		(n.Move.Turn == shogi.TurnWhite && n.Result == ResultF) {
		n.setPhi(inf)
		n.setDelta(0)
	} else {
		n.setPhi(0)
		n.setDelta(inf)
	}
	if n.parent == nil {
		return n.Result
	}

	checkSibling := false
	if (result == ResultT && n.Move.Turn == shogi.TurnWhite) ||
		(result == ResultF && n.Move.Turn == shogi.TurnBlack) {
		checkSibling = true
	}
	if checkSibling {
		ok := true
		for _, sibling := range n.parent.Children {
			if sibling.Result != result {
				ok = false
				break
			}
		}
		if ok {
			n.parent.setResult(result)
		}
	} else {
		n.parent.setResult(result)
	}
	return ResultU
}
