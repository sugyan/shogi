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
	root := NewSolver().Search(state, 0)
	return SearchBestAnswer(root)
}

// Search method
func (s *Solver) Search(state *shogi.State, maxDepth int) node.Node {
	root := dfpn.NewNode(state, shogi.TurnBlack)
	dfpn.Solve(root, maxDepth)
	searcher := &searcher{
		solved: map[string]node.Node{},
	}
	searcher.searchSolved(root)
	for {
		answer := searcher.searchAnswers(root, &shogi.CapturedPieces{}, []string{})
		l := len(answer)
		n := searchUnknownNode(root, l, answer)
		if n == nil {
			break
		}
		dfpn.Solve(n.(*dfpn.Node), l)
		if n.Result() == node.ResultT {
			searcher.solved[n.Hash()] = n
		}
	}
	return root
}

func searchUnknownNode(n node.Node, maxDepth int, answer []*shogi.Move) node.Node {
	// search around the provisional best answer.
	if answer != nil && len(answer) > 0 {
		for _, c := range n.Children() {
			if *c.Move() == *answer[0] {
				result := searchUnknownNode(c, maxDepth-1, answer[1:])
				if result != nil {
					return result
				}
			}
		}
	}
	// depth-first search
	if maxDepth > 0 {
		for _, c := range n.Children() {
			if c.Result() == node.ResultT {
				depth := maxDepth - 1
				if n.Move().Turn == shogi.TurnWhite && n.Move().Src.IsCaptured() && c.Move().Dst == n.Move().Dst {
					depth += 2
				}
				result := searchUnknownNode(c, depth, nil)
				if result != nil {
					return result
				}
			}
		}
	}
	for _, c := range n.Children() {
		if c.Result() == node.ResultU {
			return c
		}
	}
	return nil
}
