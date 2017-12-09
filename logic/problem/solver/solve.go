package solver

import (
	"errors"
	"math"

	"github.com/sugyan/shogi"
)

// DefaultMaxDepth constant
const DefaultMaxDepth = 3

// Solver type
type Solver struct {
	maxDepth int
	solved   map[string][][]*shogi.Move
}

// NewSolver function
func NewSolver(maxDepth int) *Solver {
	return &Solver{
		maxDepth: maxDepth,
		solved:   map[string][][]*shogi.Move{},
	}
}

// Solve function
func Solve(state *shogi.State) ([]*shogi.Move, error) {
	answers, length := NewSolver(DefaultMaxDepth).ValidAnswers(state)
	var answer []*shogi.Move
	switch len(answers) {
	case 0:
		return nil, errors.New("unsolvable")
	case 1:
		answer = answers[0]
	default:
		// evaluate answers
		pointMap := map[int]float64{}
		for i, answer := range answers {
			pointMap[i] = 0.0
			if len(answer) != length {
				pointMap[i] = math.Inf(-1)
				continue
			}
			s := state.Clone()
			for j := 0; j < length; j++ {
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
		answer = answers[maxIndex]
	}
	return answer, nil
}

// ValidAnswers method
func (s *Solver) ValidAnswers(state *shogi.State) ([][]*shogi.Move, int) {
	answers := s.Solve(state, 0)
	if len(answers) == 0 {
		return [][]*shogi.Move{}, 0
	}
	if len(answers) == 1 {
		return answers, len(answers[0])
	}
	length := 0
	for _, answer := range answers {
		if len(answer) > length {
			length = len(answer)
		}
	}
	results := [][]*shogi.Move{}
	for _, answer := range answers {
		if len(answer) != length {
			continue
		}
		s := state.Clone()
		for _, move := range answer {
			s.Apply(move)
		}
		results = append(results, answer)
	}
	return results, length
}

// Solve method
func (s *Solver) Solve(state *shogi.State, n int) [][]*shogi.Move {
	answers := [][]*shogi.Move{}
	if state.Check(shogi.TurnBlack) != nil {
		return answers
	}
	hash := state.Hash()
	// TODO: find endless repetition...
	if n >= s.maxDepth {
		return answers
	}
	if result, exist := s.solved[hash]; exist {
		return result
	}

	candidates := candidates(state)
	// 1 step solving
	for _, move := range candidates {
		ss := state.Clone()
		ss.Apply(move)
		if len(counterMoves(ss)) == 0 {
			if !(move.Piece == shogi.FU && move.Src == shogi.Pos(0, 0)) {
				answers = append(answers, []*shogi.Move{move})
			}
		}
	}
	if len(answers) > 0 {
		s.solved[hash] = answers
		return answers
	}

	// recursive solving
	for _, move := range candidates {
		ss := state.Clone()
		ss.Apply(move)
		counterMoves := counterMoves(ss)

		ok := true
		// simple check
		for _, counterMove := range counterMoves {
			sss := ss.Clone()
			sss.Apply(counterMove)
			if isImpossible(sss) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		candidateAnswers := [][]*shogi.Move{}
		for _, counterMove := range counterMoves {
			nextState := ss.Clone()
			nextState.Apply(counterMove)
			solved := s.Solve(nextState, n+1)
			if len(solved) > 0 {
				for _, answer := range solved {
					candidateAnswers = append(candidateAnswers, append([]*shogi.Move{move, counterMove}, answer...))
				}
			} else {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		answers = append(answers, candidateAnswers...)
	}
	s.solved[hash] = answers
	return answers
}

// IsCheckmate method
func (s *Solver) IsCheckmate(state *shogi.State) bool {
	if state.Check(shogi.TurnBlack) != nil && len(counterMoves(state)) == 0 {
		return true
	}
	return false
}

func isImpossible(state *shogi.State) bool {
	candidates := candidates(state)
	if len(candidates) == 0 {
		return true
	}
	// number of pieces on board
	pieces := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			bp := state.GetBoard(file, rank)
			if bp != nil && bp.Turn == shogi.TurnBlack {
				pieces++
			}
		}
	}
	// OU's movable positions
	positions := map[shogi.Position]struct{}{}
	for _, m := range state.CandidateMoves(shogi.TurnWhite) {
		if m.Piece == shogi.OU {
			positions[m.Dst] = struct{}{}
		}
	}
	if len(positions) > 0 {
		if pieces == 0 || (pieces == 1 && state.Captured[shogi.TurnBlack].Num() == 0) {
			return true
		}
	}
	return false
}

func candidates(state *shogi.State) []*shogi.Move {
	results := []*shogi.Move{}

	for _, move := range state.CandidateMoves(shogi.TurnBlack) {
		s := state.Clone()
		s.Apply(move)
		if s.Check(shogi.TurnBlack) != nil {
			results = append(results, move)
		}
	}

	var targetFile, targetRank int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			bp := state.GetBoard(file, rank)
			if bp != nil && bp.Turn == shogi.TurnWhite && bp.Piece == shogi.OU {
				targetFile, targetRank = file, rank
			}
		}
	}
	for _, piece := range state.Captured[shogi.TurnBlack].Available() {
		d := []shogi.Position{}
		switch piece {
		case shogi.FU:
			d = []shogi.Position{
				shogi.Pos(0, 1),
			}
		case shogi.KY:
			for i := 1; targetRank+i < 10; i++ {
				if state.GetBoard(targetFile, targetRank+i) == nil {
					d = append(d, shogi.Pos(0, i))
				} else {
					break
				}
			}
		case shogi.KE:
			d = []shogi.Position{
				shogi.Pos(1, 2),
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
			for i := 1; targetFile+i < 10 && targetRank+i < 10; i++ {
				if state.GetBoard(targetFile+i, targetRank+i) == nil {
					d = append(d, shogi.Pos(+i, +i))
				} else {
					break
				}
			}
			for i := 1; targetFile+i < 10 && targetRank-i > 0; i++ {
				if state.GetBoard(targetFile+i, targetRank-i) == nil {
					d = append(d, shogi.Pos(+i, -i))
				} else {
					break
				}
			}
			for i := 1; targetFile-i > 0 && targetRank+i < 10; i++ {
				if state.GetBoard(targetFile-i, targetRank+i) == nil {
					d = append(d, shogi.Pos(-i, +i))
				} else {
					break
				}
			}
			for i := 1; targetFile-i > 0 && targetRank-i > 0; i++ {
				if state.GetBoard(targetFile-i, targetRank-i) == nil {
					d = append(d, shogi.Pos(-i, -i))
				} else {
					break
				}
			}
		case shogi.HI:
			for i := 1; targetFile+i < 10; i++ {
				if state.GetBoard(targetFile+i, targetRank) == nil {
					d = append(d, shogi.Pos(+i, 0))
				} else {
					break
				}
			}
			for i := 1; targetRank+i < 10; i++ {
				if state.GetBoard(targetFile, targetRank+i) == nil {
					d = append(d, shogi.Pos(0, +i))
				} else {
					break
				}
			}
			for i := 1; targetFile-i > 0; i++ {
				if state.GetBoard(targetFile-i, targetRank) == nil {
					d = append(d, shogi.Pos(-i, 0))
				} else {
					break
				}
			}
			for i := 1; targetRank-i > 0; i++ {
				if state.GetBoard(targetFile, targetRank-i) == nil {
					d = append(d, shogi.Pos(0, -i))
				} else {
					break
				}
			}
		}
		for _, p := range d {
			file, rank := targetFile+p.File, targetRank+p.Rank
			if (file > 0 && file < 10 && rank > 0 && rank < 10) && state.GetBoard(file, rank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnBlack,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(file, rank),
					Piece: piece,
				})
			}
		}
	}
	return results
}

func counterMoves(state *shogi.State) []*shogi.Move {
	results := []*shogi.Move{}

	for _, move := range state.CandidateMoves(shogi.TurnWhite) {
		s := state.Clone()
		s.Apply(move)
		check := s.Check(shogi.TurnBlack)
		if check == nil {
			results = append(results, move)
		}
	}

	var targetFile, targetRank int
searchTarget:
	for i := 0; i < 9; i++ {
		rank := i + 1
		for j := 0; j < 9; j++ {
			file := 9 - j
			bp := state.GetBoard(file, rank)
			if bp != nil && bp.Turn == shogi.TurnWhite && bp.Piece == shogi.OU {
				targetFile, targetRank = file, rank
				break searchTarget
			}
		}
	}
	if state.Captured[shogi.TurnWhite].Num() > 0 {
		available := state.Captured[shogi.TurnWhite].Available()
		positions := []shogi.Position{}
		for _, direction := range []shogi.Position{
			shogi.Pos(-1, -1), shogi.Pos(-1, +1), shogi.Pos(+1, -1), shogi.Pos(+1, +1),
			shogi.Pos(-1, +0), shogi.Pos(+1, +0), shogi.Pos(+0, -1), shogi.Pos(+0, +1),
		} {
			candidates := []shogi.Position{}
			for i := 1; ; i++ {
				file := targetFile + i*direction.File
				rank := targetRank + i*direction.Rank
				if !(file > 0 && file < 10 && rank > 0 && rank < 10) {
					break
				}
				bp := state.GetBoard(file, rank)
				if bp == nil {
					candidates = append(candidates, shogi.Pos(file, rank))
				} else {
					d := direction.File * direction.Rank
					if d == 0 && bp.Turn == shogi.TurnBlack &&
						((bp.Piece == shogi.HI || bp.Piece == shogi.RY) ||
							(direction.Rank == 1 && bp.Piece == shogi.KY)) {
						positions = append(positions, candidates...)
					}
					if d != 0 && bp.Turn == shogi.TurnBlack &&
						(bp.Piece == shogi.KA || bp.Piece == shogi.UM) {
						positions = append(positions, candidates...)
					}
					break
				}
			}
			if len(positions) > 0 {
				break
			}
		}
		if len(positions) > 0 {
			movableF := map[shogi.Position][]*shogi.Move{}
			movableS := map[shogi.Position][]*shogi.Move{}
			for _, m := range state.CandidateMoves(shogi.TurnBlack) {
				movableF[m.Dst] = append(movableF[m.Dst], m)
			}
			for _, m := range state.CandidateMoves(shogi.TurnWhite) {
				movableS[m.Dst] = append(movableS[m.Dst], m)
			}
			for _, p := range positions {
				// check wasted placing
				if moves, exist := movableS[p]; exist && len(moves) == 1 && moves[0].Piece == shogi.OU {
					src := map[shogi.Position]struct{}{}
					for _, m := range movableF[p] {
						src[m.Src] = struct{}{}
					}
					if len(src) > 1 {
						continue
					}
				}
				for _, piece := range available {
					// check duplicated FU
					if piece == shogi.FU {
						ok := true
						for rank := 1; rank < 10; rank++ {
							bp := state.GetBoard(p.File, rank)
							if bp != nil && bp.Turn == shogi.TurnWhite && bp.Piece == shogi.FU {
								ok = false
								break
							}
						}
						if !ok {
							continue
						}
					}
					move := &shogi.Move{
						Turn:  shogi.TurnWhite,
						Src:   shogi.Pos(0, 0),
						Dst:   shogi.Pos(p.File, p.Rank),
						Piece: piece,
					}
					s := state.Clone()
					s.Apply(move)
					check := s.Check(shogi.TurnBlack)
					if check == nil {
						results = append(results, move)
					}
				}
			}
		}
	}
	return results
}
