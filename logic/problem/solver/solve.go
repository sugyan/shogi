package solver

import (
	"math"

	"github.com/sugyan/shogi"
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
	return selectBestAnswer(state, answers)[1:]
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
	tree := newTree(state)

	solved := s.do(state, tree.root)
	if solved {
		return tree.answers()
	}
	for _, node := range tree.root.leaves() {
		solved = solved || s.do(node.moveState.state, node)
	}
	if solved {
		return tree.answers()
	}
	return [][]*shogi.Move{}
}

func (s *Solver) do(state *shogi.State, node *node) bool {
	for _, ms := range candidates(state, shogi.TurnBlack) {
		node.addChildNode(ms)
	}

	solved := false
	for _, child := range node.childNodes {
		candidates := candidates(child.moveState.state, shogi.TurnWhite)
		if len(candidates) == 0 {
			result := child.setResult(resultTrue)
			if result == resultTrue {
				solved = true
			}
		}
		for _, ms := range candidates {
			child.addChildNode(ms)
		}
	}
	return solved
}

func searchTarget(state *shogi.State) *shogi.Position {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := state.Board[i][j]
			if bp != nil && bp.Piece == shogi.OU && bp.Turn == shogi.TurnWhite {
				return &shogi.Position{File: 9 - j, Rank: i + 1}
			}
		}
	}
	return nil
}

func candidates(state *shogi.State, turn shogi.Turn) []*moveState {
	results := []*moveState{}
	target := *searchTarget(state)
	// by moving pieces
	for _, m := range state.CandidateMoves(turn) {
		s := state.Clone()
		s.Apply(m)
		check := s.Check(shogi.TurnBlack) != nil
		if (turn == shogi.TurnBlack && check) || (turn == shogi.TurnWhite && !check) {
			results = append(results, &moveState{m, s})
		}
	}
	switch turn {
	case shogi.TurnBlack:
		// check by placing captured pieces
		for _, piece := range state.Captured[shogi.TurnBlack].Available() {
			d := []shogi.Position{}
			switch piece {
			case shogi.FU:
				d = []shogi.Position{
					shogi.Pos(0, 1),
				}
			case shogi.KY:
				for i := 1; target.Rank+i < 10; i++ {
					if state.GetBoard(target.File, target.Rank+i) == nil {
						d = append(d, shogi.Pos(0, i))
					} else {
						break
					}
				}
			case shogi.KE:
				d = []shogi.Position{
					shogi.Pos(+1, 2),
					shogi.Pos(-1, 2),
				}
			case shogi.GI:
				d = []shogi.Position{
					shogi.Pos(-1, -1),
					shogi.Pos(+1, -1),
					shogi.Pos(+0, +1),
					shogi.Pos(-1, +1),
					shogi.Pos(+1, +1),
				}
			case shogi.KI:
				d = []shogi.Position{
					shogi.Pos(+0, -1),
					shogi.Pos(-1, +0),
					shogi.Pos(+1, +0),
					shogi.Pos(-1, +1),
					shogi.Pos(+0, +1),
					shogi.Pos(+1, +1),
				}
			case shogi.KA:
				for _, direction := range []shogi.Position{
					shogi.Pos(-1, -1),
					shogi.Pos(-1, +1),
					shogi.Pos(+1, -1),
					shogi.Pos(+1, +1),
				} {
					for i := 1; ; i++ {
						file, rank := target.File+direction.File*i, target.Rank+direction.Rank*i
						if file > 0 && file < 10 && rank > 0 && rank < 10 &&
							state.GetBoard(file, rank) == nil {
							d = append(d, shogi.Pos(direction.File*i, direction.Rank*i))
						} else {
							break
						}
					}
				}
			case shogi.HI:
				for _, direction := range []shogi.Position{
					shogi.Pos(-1, +0),
					shogi.Pos(+1, +0),
					shogi.Pos(+0, -1),
					shogi.Pos(+0, +1),
				} {
					for i := 1; ; i++ {
						file, rank := target.File+direction.File*i, target.Rank+direction.Rank*i
						if file > 0 && file < 10 && rank > 0 && rank < 10 &&
							state.GetBoard(file, rank) == nil {
							d = append(d, shogi.Pos(direction.File*i, direction.Rank*i))
						} else {
							break
						}
					}
				}
			}
			for _, pos := range d {
				file, rank := target.File+pos.File, target.Rank+pos.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 &&
					state.GetBoard(file, rank) == nil {
					m := &shogi.Move{
						Turn:  shogi.TurnBlack,
						Src:   shogi.Pos(0, 0),
						Dst:   shogi.Pos(file, rank),
						Piece: piece,
					}
					s := state.Clone()
					s.Apply(m)
					results = append(results, &moveState{m, s})
				}
			}
		}
	case shogi.TurnWhite:
		// TODO
	}
	return results
}
