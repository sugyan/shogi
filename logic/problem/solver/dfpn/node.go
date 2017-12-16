package dfpn

import (
	"github.com/sugyan/shogi"
)

type result int

const (
	resultU = iota // Unknown
	resultT        // True
	resultF        // False
)

// Node type
type Node struct {
	result   result
	pn       uint32
	dn       uint32
	move     *shogi.Move
	state    *shogi.State
	parent   *Node
	Children []*Node
}

func (n *Node) getPhi() uint32 {
	switch n.move.Turn {
	case shogi.TurnBlack:
		return n.dn
	case shogi.TurnWhite:
		return n.pn
	}
	return 0
}

func (n *Node) setPhi(v uint32) {
	switch n.move.Turn {
	case shogi.TurnBlack:
		n.dn = v
	case shogi.TurnWhite:
		n.pn = v
	}
}

func (n *Node) getDelta() uint32 {
	switch n.move.Turn {
	case shogi.TurnBlack:
		return n.pn
	case shogi.TurnWhite:
		return n.dn
	}
	return 0
}

func (n *Node) setDelta(v uint32) {
	switch n.move.Turn {
	case shogi.TurnBlack:
		n.pn = v
	case shogi.TurnWhite:
		n.dn = v
	}
}

func (n *Node) setResult(result result) result {
	n.result = result
	if (n.move.Turn == shogi.TurnBlack && n.result == resultT) ||
		(n.move.Turn == shogi.TurnWhite && n.result == resultF) {
		n.setPhi(inf)
		n.setDelta(0)
	} else {
		n.setPhi(0)
		n.setDelta(inf)
	}
	if n.parent == nil {
		return n.result
	}
	switch result {
	case resultT:
		switch n.move.Turn {
		case shogi.TurnBlack:
			return n.parent.setResult(resultT)
		case shogi.TurnWhite:
			ok := true
			for _, sibling := range n.parent.Children {
				if sibling.result != resultT {
					ok = false
					break
				}
			}
			if ok {
				return n.parent.setResult(resultT)
			}
		}
	case resultF:
		switch n.move.Turn {
		case shogi.TurnBlack:
			ok := true
			for _, sibling := range n.parent.Children {
				if sibling.result != resultF {
					ok = false
					break
				}
			}
			if ok {
				return n.parent.setResult(resultF)
			}
		case shogi.TurnWhite:
			return n.parent.setResult(resultF)
		}
	}
	return resultU
}
