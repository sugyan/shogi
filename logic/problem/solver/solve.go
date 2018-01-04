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
		n, depth := searchUnknownNode(root, 0)
		if n == nil || depth >= l {
			break
		}
		searcher.SetMaxDepth(l - depth)
		searcher.Search(n.(*dfpn.Node))
	}
	return root
}

func searchUnknownNode(n node.Node, d int) (node.Node, int) {
	for _, c := range n.Children() {
		if c.Result() == node.ResultU {
			return c, d + 1
		}
	}
	for _, c := range n.Children() {
		if c.Result() == node.ResultT {
			return searchUnknownNode(c, d+1)
		}
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
