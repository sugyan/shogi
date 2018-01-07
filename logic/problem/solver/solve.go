package solver

import (
	"math"

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
	root := dfpn.NewNode(state, shogi.TurnBlack)
	dfpn.NewSearcher().Search(root)
	for {
		answer := SearchBestAnswer(root)
		l := len(answer)
		n := searchUnknownNode(root, l, answer)
		if n == nil {
			break
		}
		searcher := dfpn.NewSearcher()
		searcher.SetMaxDepth(l)
		searcher.Search(n.(*dfpn.Node))
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
			if answer[0].Turn == shogi.TurnWhite && answer[0].Src.IsCaptured() && answer[1].Dst == answer[0].Dst {
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
	if len(candidates) > 1 {
		best := 0
		switch n.Move().Turn {
		case shogi.TurnBlack:
			points := map[int]int{}
			for i, answer := range candidates {
				points[i] = 0
				s := n.State().Clone()
				for _, move := range answer {
					s.Apply(move)
				}
				captured := s.Captured[shogi.TurnBlack].Num()
				points[i] -= captured * 10
			}
			max := math.MinInt32
			for k, v := range points {
				if v > max {
					max = v
					best = k
				}
			}
		case shogi.TurnWhite:
			max := 0
			for i, answer := range candidates {
				s := n.State().Clone()
				for _, move := range answer {
					s.Apply(move)
				}
				captured := s.Captured[shogi.TurnBlack].Num()
				if captured > max {
					best = i
					max = captured
				}
			}
		}
		return candidates[best]
	}
	return candidates[0]
}
