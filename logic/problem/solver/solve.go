package solver

import (
	"log"
	"math"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/dfpn"
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
	answers := NewSolver().Answers(state)
	for _, answer := range answers {
		ms, _ := state.MoveStrings(answer[1:])
		log.Printf("%v", ms)
	}

	// TODO
	// return selectBestAnswer(state, answers)[1:]
	return []*shogi.Move{}
}

func selectBestAnswer(state *shogi.State, answers [][]*shogi.Move) []*shogi.Move {
	pointMap := map[int]float64{}
	for i, answer := range answers {
		pointMap[i] = 0.0
		s := state.Clone()
		for j := 1; j < len(answer); j++ {
			move := answer[j]
			s.Apply(move)
			if j > 0 {
				prev := answer[j-1]
				if move.Turn == shogi.TurnWhite && move.Dst == prev.Dst {
					pointMap[i] += 1.0
				}
			}
			if move.Turn == shogi.TurnWhite && move.Src == shogi.Pos(0, 0) {
				switch move.Piece {
				case shogi.FU:
					pointMap[i] -= 0.1
				case shogi.KY:
					pointMap[i] -= 0.2
				case shogi.KE:
					pointMap[i] -= 0.3
				case shogi.GI:
					pointMap[i] -= 0.4
				case shogi.KI:
					pointMap[i] -= 0.5
				case shogi.KA:
					pointMap[i] -= 0.6
				case shogi.HI:
					pointMap[i] -= 0.7
				}
			}
		}
		if s.Captured[shogi.TurnBlack].Num() > 0 {
			pointMap[i] -= 10
		}
	}
	maxIndex, point := 0, math.Inf(-1)
	for k, v := range pointMap {
		if v > point {
			point = v
			maxIndex = k
		}
	}
	return answers[maxIndex]
}

// Answers method
func (s *Solver) Answers(state *shogi.State) [][]*shogi.Move {
	root := &dfpn.Node{
		Move: &shogi.Move{
			Turn: shogi.TurnWhite,
		},
		State: state,
	}
	dfpn.NewSolver().Solve(root)
	searchMoreAnswers(root)

	return collectAnswers(root)
}

func searchMoreAnswers(n *dfpn.Node) {
	for _, child := range n.Children {
		switch child.Result {
		case dfpn.ResultU:
			dfpn.NewSolver().Solve(child)
		case dfpn.ResultT:
			searchMoreAnswers(child)
		}
	}
}

func collectAnswers(n *dfpn.Node) [][]*shogi.Move {
	if len(n.Children) == 0 {
		return [][]*shogi.Move{[]*shogi.Move{n.Move}}
	}

	results := [][]*shogi.Move{}
	for _, child := range n.Children {
		if child.Result == dfpn.ResultT {
			for _, answer := range collectAnswers(child) {
				result := append([]*shogi.Move{n.Move}, answer...)
				results = append(results, result)
			}
		}
	}
	return results
}
