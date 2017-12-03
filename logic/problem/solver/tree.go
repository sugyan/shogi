package solver

import (
	"fmt"
	"strings"

	"github.com/sugyan/shogi"
)

type tree struct {
	root *node
}

type nodeResult rune

const (
	resultUnknown = nodeResult('?')
	resultTrue    = nodeResult('T')
	resultFalse   = nodeResult('F')
)

type moveState struct {
	move  *shogi.Move
	state *shogi.State
}

type node struct {
	result     nodeResult
	parent     *node
	moveState  *moveState
	childNodes []*node
}

func newTree(state *shogi.State) *tree {
	return &tree{
		root: &node{
			result: resultUnknown,
			moveState: &moveState{
				state: state,
			},
		},
	}
}

func (t *tree) answers() [][]*shogi.Move {
	results := t.root.collectAnswers()
	return results
}

func (t *tree) string() string {
	return t.root.archy("")[1:]
}

func (n *node) archy(prefix string) string {
	line := ""
	if n.parent != nil {
		move := n.moveState.move
		sign := '+'
		if move.Turn == shogi.TurnSecond {
			sign = '-'
		}
		line = fmt.Sprintf("%c%d%d%d%d%s (%c)", sign, move.Src.File, move.Src.Rank, move.Dst.File, move.Dst.Rank, move.Piece, n.result)
	}

	nodeLines := []string{}
	for i, node := range n.childNodes {
		_prefix := prefix + " "
		if i < len(n.childNodes)-1 {
			_prefix = prefix + "│"
		}
		_prefix += " "

		l := "└"
		if i < len(n.childNodes)-1 {
			l = "├"
		}
		m := "─"
		if len(node.childNodes) > 0 {
			m = "┬"
		}
		line := prefix +
			l + "─" +
			m + " " +
			string([]rune(node.archy(_prefix))[len([]rune(prefix))+2:])
		nodeLines = append(nodeLines, line)

	}
	return prefix +
		line + "\n" +
		strings.Join(nodeLines, "")
}

func (n *node) addChildNode(moveState *moveState) {
	n.childNodes = append(n.childNodes, &node{
		result:    resultUnknown,
		parent:    n,
		moveState: moveState,
	})
}

func (n *node) leaves() []*node {
	results := []*node{}
	if len(n.childNodes) > 0 {
		for _, child := range n.childNodes {
			results = append(results, child.leaves()...)
		}
	} else {
		return []*node{n}
	}
	return results
}

func (n *node) setResult(result nodeResult) nodeResult {
	n.result = result
	if n.parent == nil {
		return n.result
	}
	switch result {
	case resultTrue:
		switch n.moveState.move.Turn {
		case shogi.TurnFirst:
			return n.parent.setResult(resultTrue)
		case shogi.TurnSecond:
			ok := true
			for _, sibling := range n.parent.childNodes {
				if sibling.result != resultTrue {
					ok = false
					break
				}
			}
			if ok {
				return n.parent.setResult(resultTrue)
			}
		}
	case resultFalse:
		// TODO
	}
	return resultUnknown
}

func (n *node) collectAnswers() [][]*shogi.Move {
	if len(n.childNodes) == 0 {
		return [][]*shogi.Move{
			[]*shogi.Move{n.moveState.move},
		}
	}
	results := [][]*shogi.Move{}
	for _, child := range n.childNodes {
		if child.result == resultTrue {
			for _, answers := range child.collectAnswers() {
				results = append(results,
					append([]*shogi.Move{n.moveState.move}, answers...),
				)
			}
		}
	}
	return results
}
