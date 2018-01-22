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
	dfpn.NewSolver().Solve(root, maxDepth)
	for {
		answer := SearchBestAnswer(root)
		l := len(answer)
		n := searchUnknownNode(root, l, answer)
		if n == nil {
			break
		}
		dfpn.NewSolver().Solve(n.(*dfpn.Node), l)
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
				result := searchUnknownNode(c, maxDepth-1, nil)
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
