package dfpn

import (
	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

// Node type
type Node struct {
	result   node.Result
	move     *shogi.Move
	state    *shogi.State
	children []*Node
	expanded bool
	pn, dn   uint32
	depth    int
	hash     string
}

// NewNode function
func NewNode(state *shogi.State, turn shogi.Turn) *Node {
	return &Node{
		state: state,
		move: &shogi.Move{
			Turn: !turn,
		},
		hash: state.Hash(),
	}
}

// Children method
func (n *Node) Children() []node.Node {
	result := make([]node.Node, 0, len(n.children))
	for _, c := range n.children {
		result = append(result, c)
	}
	return result
}

// Result method
func (n *Node) Result() node.Result {
	return n.result
}

// Move method
func (n *Node) Move() *shogi.Move {
	return n.move
}

// State method
func (n *Node) State() *shogi.State {
	return n.state
}

// Hash method
func (n *Node) Hash() string {
	return n.hash
}

func (n *Node) getP() uint32 {
	switch n.move.Turn {
	case shogi.TurnBlack:
		return n.dn
	case shogi.TurnWhite:
		return n.pn
	}
	return 0
}

func (n *Node) setP(v uint32) {
	switch n.move.Turn {
	case shogi.TurnBlack:
		n.dn = v
	case shogi.TurnWhite:
		n.pn = v
	}
}

func (n *Node) getD() uint32 {
	switch n.move.Turn {
	case shogi.TurnBlack:
		return n.pn
	case shogi.TurnWhite:
		return n.dn
	}
	return 0
}

func (n *Node) setD(v uint32) {
	switch n.move.Turn {
	case shogi.TurnBlack:
		n.pn = v
	case shogi.TurnWhite:
		n.dn = v
	}
}
