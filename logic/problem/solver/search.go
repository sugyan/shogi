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
	return s.searchAnswers(n, &shogi.CapturedPieces{}, []string{})
}

func (s *searcher) searchSolved(n node.Node) {
	if n.Result() == node.ResultT {
		s.solved[n.Hash()] = n
	}
	for _, c := range n.Children() {
		s.searchSolved(c)
	}
}

func (s *searcher) searchAnswers(n node.Node, captured *shogi.CapturedPieces, ancestors []string) []*shogi.Move {
	if len(n.Children()) == 0 {
		return []*shogi.Move{}
	}
	hash := n.Hash()
	for _, ancestor := range ancestors {
		if ancestor == hash {
			return []*shogi.Move{}
		}
	}
	ancestors = append(ancestors, hash)

	answers := [][]*shogi.Move{}
	for _, c := range n.Children() {
		move := c.Move()
		if solved, exist := s.solved[c.Hash()]; exist {
			c = solved
		}
		if c.Result() != node.ResultT {
			continue
		}
		cap := *captured
		if move.Turn == shogi.TurnBlack {
			dst := move.Dst
			bp := n.State().GetBoard(dst.File, dst.Rank)
			if bp != nil && bp.Turn == shogi.TurnWhite {
				cap.Add(bp.Piece)
			}
		}
		answer := append([]*shogi.Move{move}, s.searchAnswers(c, &cap, ancestors)...)

		omit := false
		if len(answer) > 1 {
			if answer[0].Turn == shogi.TurnWhite && answer[0].Src.IsCaptured() && answer[1].Dst == answer[0].Dst {
				state := n.State().Clone()
				for _, m := range answer {
					state.Apply(m)
				}
				if false ||
					answer[0].Piece == shogi.FU && state.Captured[shogi.TurnBlack].FU > 0 ||
					answer[0].Piece == shogi.KY && state.Captured[shogi.TurnBlack].KY > 0 ||
					answer[0].Piece == shogi.KE && state.Captured[shogi.TurnBlack].KE > 0 ||
					answer[0].Piece == shogi.GI && state.Captured[shogi.TurnBlack].GI > 0 ||
					answer[0].Piece == shogi.KI && state.Captured[shogi.TurnBlack].KI > 0 ||
					answer[0].Piece == shogi.KA && state.Captured[shogi.TurnBlack].KA > 0 ||
					answer[0].Piece == shogi.HI && state.Captured[shogi.TurnBlack].HI > 0 {
					omit = true
				}
			}
		}
		if !omit {
			if answer[len(answer)-1].Turn == shogi.TurnBlack {
				answers = append(answers, answer)
			}
		}
	}
	if len(answers) == 0 {
		return []*shogi.Move{}
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
			cap := state.Captured[shogi.TurnBlack]
			points[i] -= cap.Num()
			if captured.FU > 0 && cap.FU > 0 {
				points[i]++
			}
			if captured.KY > 0 && cap.KY > 0 {
				points[i]++
			}
			if captured.KE > 0 && cap.KE > 0 {
				points[i]++
			}
			if captured.GI > 0 && cap.GI > 0 {
				points[i]++
			}
			if captured.KI > 0 && cap.KI > 0 {
				points[i]++
			}
			if captured.KA > 0 && cap.KA > 0 {
				points[i]++
			}
			if captured.HI > 0 && cap.HI > 0 {
				points[i]++
			}
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
