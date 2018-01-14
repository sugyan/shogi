package solver

import (
	"math"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

// SearchBestAnswer function
func SearchBestAnswer(n node.Node) []*shogi.Move {
	return searchBestAnswer(n, &shogi.CapturedPieces{})
}

func searchBestAnswer(n node.Node, captured *shogi.CapturedPieces) []*shogi.Move {
	if len(n.Children()) == 0 {
		return []*shogi.Move{}
	}

	answers := [][]*shogi.Move{}
	for _, c := range n.Children() {
		if c.Result() != node.ResultT {
			continue
		}
		cap := *captured
		if c.Move().Turn == shogi.TurnBlack {
			dst := c.Move().Dst
			bp := n.State().GetBoard(dst.File, dst.Rank)
			if bp != nil && bp.Turn == shogi.TurnWhite {
				cap.Add(bp.Piece)
			}
		}
		answer := append([]*shogi.Move{c.Move()}, searchBestAnswer(c, &cap)...)
		omit := false
		if len(answer) > 1 {
			if answer[0].Turn == shogi.TurnWhite && answer[0].Src.IsCaptured() && answer[1].Dst == answer[0].Dst {
				s := n.State().Clone()
				for _, m := range answer {
					s.Apply(m)
				}
				if false ||
					answer[0].Piece == shogi.FU && s.Captured[shogi.TurnBlack].FU > 0 ||
					answer[0].Piece == shogi.KY && s.Captured[shogi.TurnBlack].KY > 0 ||
					answer[0].Piece == shogi.KE && s.Captured[shogi.TurnBlack].KE > 0 ||
					answer[0].Piece == shogi.GI && s.Captured[shogi.TurnBlack].GI > 0 ||
					answer[0].Piece == shogi.KI && s.Captured[shogi.TurnBlack].KI > 0 ||
					answer[0].Piece == shogi.KA && s.Captured[shogi.TurnBlack].KA > 0 ||
					answer[0].Piece == shogi.HI && s.Captured[shogi.TurnBlack].HI > 0 {
					omit = true
				}
			}
		}
		if !omit {
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
				cap := s.Captured[shogi.TurnBlack]
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
	return candidates[0]
}
