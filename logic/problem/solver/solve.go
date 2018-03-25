package solver

import (
	"context"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/dfpn"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

// Solver type
type Solver struct {
	state  *shogi.State
	cancel bool
}

// NewSolver function
func NewSolver(state *shogi.State) *Solver {
	return &Solver{
		state: state,
	}
}

// Solve function
func Solve(state *shogi.State) []*shogi.Move {
	root := NewSolver(state).solve()
	return SearchBestAnswer(root)
}

// SolveWithTimeout method
func (s *Solver) SolveWithTimeout(timeout time.Duration) (node.Node, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	defer cancel()

	c := make(chan node.Node)
	go func(s *Solver) {
		result := s.solve()
		if result != nil {
			c <- result
		}
	}(s)
	defer close(c)

	select {
	case <-ctx.Done():
		s.cancel = true
		return nil, ctx.Err()
	case ret := <-c:
		return ret, nil
	}
}

// solve method
func (s *Solver) solve() node.Node {
	root := dfpn.NewNode(s.state, shogi.TurnBlack)
	dfpn.Solve(root, 0)
	if root.Result() == node.ResultF {
		return root
	}

	searcher := &searcher{
		solved: map[string]node.Node{},
	}
	searcher.searchSolved(root)

	answer := searcher.searchBestAnswer(root, []string{})
	for {
		if s.cancel {
			return nil
		}
		l := len(answer)
		n := searchUnknownNode(root, l, answer)
		if n == nil {
			break
		}
		dfpn.Solve(n.(*dfpn.Node), l)
		if n.Result() == node.ResultT {
			searcher.solved[n.Hash()] = n
			answer = searcher.searchBestAnswer(root, []string{})
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
