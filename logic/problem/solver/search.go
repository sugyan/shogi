package solver

import (
	"math"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

type searcher struct {
	solved map[string]node.Node
}

// SearchBestAnswer function
func SearchBestAnswer(n node.Node) []*shogi.Move {
	s := &searcher{
		solved: map[string]node.Node{},
	}
	s.searchSolved(n)
	return s.searchBestAnswer(n, []string{})
}

func (s *searcher) searchSolved(n node.Node) {
	if n.Result() == node.ResultT {
		s.solved[n.Hash()] = n
	}
	for _, c := range n.Children() {
		s.searchSolved(c)
	}
}

func (s *searcher) searchBestAnswer(n node.Node, ancestors []string) []*shogi.Move {
	if len(n.Children()) == 0 {
		return []*shogi.Move{}
	}
	ancestors = append(ancestors, n.Hash())

	answers := [][]*shogi.Move{}
	omitted2 := false
	for _, c := range n.Children() {
		move := c.Move()
		if solved, exist := s.solved[c.Hash()]; exist {
			c = solved
		}
		if c.Result() != node.ResultT {
			continue
		}
		// prevent cycle
		ok := true
		for _, ancestor := range ancestors {
			if ancestor == c.Hash() {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		answer := append([]*shogi.Move{move}, s.searchBestAnswer(c, ancestors)...)

		omit := false
		if len(answer) > 1 {
			if move.Turn == shogi.TurnWhite && move.Src.IsCaptured() && answer[1].Dst == move.Dst {
				isLive := true
				pos := move.Dst
				state := n.State().Clone()
				for i, m := range answer {
					if m.Turn == shogi.TurnBlack && m.Src == pos {
						pos = m.Dst
					}
					if m.Turn == shogi.TurnWhite && m.Dst == pos && i > 0 {
						isLive = false
					}
					state.Apply(m)
				}
				isLeft := false ||
					(move.Piece == shogi.FU && state.Captured[shogi.TurnBlack].FU > 0) ||
					(move.Piece == shogi.KY && state.Captured[shogi.TurnBlack].KY > 0) ||
					(move.Piece == shogi.KE && state.Captured[shogi.TurnBlack].KE > 0) ||
					(move.Piece == shogi.GI && state.Captured[shogi.TurnBlack].GI > 0) ||
					(move.Piece == shogi.KI && state.Captured[shogi.TurnBlack].KI > 0) ||
					(move.Piece == shogi.KA && state.Captured[shogi.TurnBlack].KA > 0) ||
					(move.Piece == shogi.HI && state.Captured[shogi.TurnBlack].HI > 0)
				if isLive && isLeft {
					omit = true
					if len(answer) == 2 {
						omitted2 = true
					}
				}
			}
		}
		if omit {
			continue
		}
		if answer[len(answer)-1].Turn == shogi.TurnBlack {
			answers = append(answers, answer)
		}
	}
	if len(answers) == 0 {
		return []*shogi.Move{}
	}
	// check if already checkmated
	if omitted2 {
		dst := map[shogi.Position]struct{}{}
		for _, c := range n.Children() {
			dst[c.Move().Dst] = struct{}{}
		}
		if len(dst) == 1 {
			return []*shogi.Move{}
		}
	}
	min, max := math.MaxInt32, 0
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

	if len(candidates) == 1 {
		return candidates[0]
	}
	best := 0
	switch n.Move().Turn {
	case shogi.TurnBlack:
		points := map[int]int{}
		for i, answer := range candidates {
			points[i] = 0
			state := n.State().Clone()
			for _, move := range answer {
				state.Apply(move)
			}
			points[i] -= state.Captured[shogi.TurnBlack].Num()
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
