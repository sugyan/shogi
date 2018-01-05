package solver

import (
	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/dfpn"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

// Solver type
type Solver struct {
}

// NewSolver function
func NewSolver() *Solver {
	return &Solver{}
}

// Solve function
func Solve(state *shogi.State) []*shogi.Move {
	root := NewSolver().Search(state)
	return SearchBestAnswer(root)
}

// Search method
func (s *Solver) Search(state *shogi.State) node.Node {
	searcher := dfpn.NewSearcher()
	root := dfpn.NewNode(state, shogi.TurnBlack)
	searcher.Search(root)
	for {
		l := len(SearchBestAnswer(root))
		n, depth := searchUnknownNode(root, l)
		if n == nil {
			break
		}
		searcher.SetMaxDepth(depth)
		searcher.Search(n.(*dfpn.Node))
	}
	return root
}

func searchUnknownNode(n node.Node, maxDepth int) (node.Node, int) {
	type entry struct {
		node  node.Node
		depth int
	}
	q := []*entry{&entry{node: n, depth: 0}}
	for {
		if len(q) == 0 {
			break
		}
		e := q[0]
		switch e.node.Result() {
		case node.ResultU:
			return e.node, e.depth
		case node.ResultT:
			if e.depth < maxDepth {
				for _, c := range e.node.Children() {
					q = append(q, &entry{
						node:  c,
						depth: e.depth + 1,
					})
				}
			}
		}
		q = q[1:]
	}
	return nil, 0
}

// SearchBestAnswer function
func SearchBestAnswer(n node.Node) []*shogi.Move {
	if len(n.Children()) == 0 {
		return []*shogi.Move{}
	}
	answers := [][]*shogi.Move{}
	for _, c := range n.Children() {
		if c.Result() != node.ResultT {
			continue
		}
		answer := append([]*shogi.Move{c.Move()}, SearchBestAnswer(c)...)
		ok := true
		if len(answer) > 1 {
			if answer[0].Src.IsCaptured() && answer[1].Dst == answer[0].Dst {
				s := n.State().Clone()
				for _, m := range answer {
					s.Apply(m)
				}
				for _, piece := range s.Captured[shogi.TurnBlack].Available() {
					if piece == answer[0].Piece {
						ok = false
					}
				}
			}
		}
		if ok {
			answers = append(answers, answer)
		}
	}
	if len(answers) == 0 {
		return []*shogi.Move{}
	}
	min, max := int(^uint(0)>>1), 0
	for _, answer := range answers {
		l := len(answer)
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
	}
	candidates := [][]*shogi.Move{}
	for _, answer := range answers {
		switch n.Move().Turn {
		case shogi.TurnBlack:
			if len(answer) == max {
				candidates = append(candidates, answer)
			}
		case shogi.TurnWhite:
			if len(answer) == min {
				candidates = append(candidates, answer)
			}
		}
	}
	// TODO: select from candidates
	return candidates[0]
}
