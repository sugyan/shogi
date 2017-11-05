package solver

import (
	"errors"

	"github.com/sugyan/shogi"
)

// Solver type
type Solver struct {
	state *shogi.State
}

// NewSolver function
func NewSolver(state *shogi.State) *Solver {
	return &Solver{state: state}
}

// Solve function
func Solve(state *shogi.State) ([]string, error) {
	answers := NewSolver(state).Solve([]string{}, 0)
	var answer []*shogi.Move
	switch len(answers) {
	case 0:
		return nil, errors.New("unsolvable")
	case 1:
		answer = answers[0]
	default:
		// mutlple answers
		// TODO choose shortest answers
		length := len(answers[0])
		pointMap := map[int]int{}
		for i := 0; i < length; i++ {
			for j := 0; j < len(answers); j++ {
				a := answers[j]
				if i > 0 {
					move := a[i]
					prev := a[i-1]
					if move.Turn == shogi.TurnSecond && *move.Dst == *prev.Dst {
						pointMap[j]++
					}
				}
			}
		}
		maxIndex, point := 0, 0
		for k, v := range pointMap {
			if v > point {
				point = v
				maxIndex = k
			}
		}
		answer = answers[maxIndex]
	}
	var (
		prev    *shogi.Move
		results []string
	)
	for _, move := range answer {
		ms, err := state.MoveString(move, prev)
		if err != nil {
			return nil, err
		}
		results = append(results, ms)
		state = simulate(state, move)
		prev = move
	}
	return results, nil
}

// Solve method
func (s *Solver) Solve(history []string, n int) [][]*shogi.Move {
	answers := [][]*shogi.Move{}
	if s.state.Check(shogi.TurnFirst) != nil {
		return answers
	}
	hash := s.state.Hash()
	// TODO: find endless repitition...
	if n > 3 {
		return nil
	}
	for i := len(history) - 1; i >= 0; i-- {
		if history[i] == hash {
			return nil
		}
	}

	candidates := candidates(s.state)
	// 1 step solving
	for _, move := range candidates {
		state := simulate(s.state, move)
		if len(counterMoves(state)) == 0 {
			answers = append(answers, []*shogi.Move{move})
		}
	}
	if len(answers) > 0 {
		return answers
	}

	// recursive solving
	for _, move := range candidates {
		state := simulate(s.state, move)
		counterMoves := counterMoves(state)

		ok := true
		// simple check
		for _, counterMove := range counterMoves {
			if isImpossible(simulate(state, counterMove)) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		candidateAnswers := [][]*shogi.Move{}
		for _, counterMove := range counterMoves {
			nextState := simulate(state, counterMove)
			solved := NewSolver(nextState).Solve(append(history, hash), n+1)
			if solved == nil {
				ok = false
				break
			}
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
	return answers
}

// IsCheckmate method
func (s *Solver) IsCheckmate() bool {
	if s.state.Check(shogi.TurnFirst) != nil && len(counterMoves(s.state)) == 0 {
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
			bp := state.GetBoardPiece(file, rank)
			if bp != nil && bp.Turn == shogi.TurnFirst {
				pieces++
			}
		}
	}
	// OU's movable positions
	positions := map[shogi.Position]struct{}{}
	for _, m := range state.CandidateMoves(shogi.TurnSecond) {
		if m.Piece == shogi.OU {
			positions[*m.Dst] = struct{}{}
		}
	}
	if len(positions) > 0 {
		if pieces == 0 || (pieces == 1 && state.Captured[shogi.TurnFirst].Num() == 0) {
			return true
		}
		for _, move := range state.CandidateMoves(shogi.TurnFirst) {
			if _, exist := positions[*move.Dst]; exist {
				return false
			}
		}
		return true
	}
	return false
}

func candidates(state *shogi.State) []*shogi.Move {
	results := []*shogi.Move{}
	for _, move := range state.CandidateMoves(shogi.TurnFirst) {
		if simulate(state, move).Check(shogi.TurnFirst) != nil {
			results = append(results, move)
		}
	}
	var targetFile, targetRank int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil && bp.Turn == shogi.TurnSecond && bp.Piece == shogi.OU {
				targetFile, targetRank = file, rank
			}
		}
	}
	if state.Captured[shogi.TurnFirst].FU > 0 {
		// TODO check mating with a FU drop
		if targetRank+1 < 10 {
			if state.GetBoardPiece(targetFile, targetRank+1) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile, targetRank+1),
					Piece: shogi.FU,
				})
			}
		}
	}
	if state.Captured[shogi.TurnFirst].KY > 0 {
		for i := 1; targetRank+i < 10; i++ {
			if state.GetBoardPiece(targetFile, targetRank+i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile, targetRank+i),
					Piece: shogi.KY,
				})
			} else {
				break
			}
		}
	}
	if state.Captured[shogi.TurnFirst].KE > 0 {
		for _, d := range []*shogi.Position{
			shogi.Pos(1, 2),
			shogi.Pos(-1, 2),
		} {
			if (targetFile+d.File > 0 && targetFile+d.File < 10 && targetRank+d.Rank > 0 && targetRank+d.Rank < 10) &&
				state.GetBoardPiece(targetFile+d.File, targetRank+d.Rank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+d.File, targetRank+d.Rank),
					Piece: shogi.KE,
				})
			}
		}
	}
	if state.Captured[shogi.TurnFirst].GI > 0 {
		for _, d := range []*shogi.Position{
			shogi.Pos(-1, -1),
			shogi.Pos(+1, -1),
			shogi.Pos(+0, +1),
			shogi.Pos(-1, +1),
			shogi.Pos(+1, +1),
		} {
			if (targetFile+d.File > 0 && targetFile+d.File < 10 && targetRank+d.Rank > 0 && targetRank+d.Rank < 10) &&
				state.GetBoardPiece(targetFile+d.File, targetRank+d.Rank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+d.File, targetRank+d.Rank),
					Piece: shogi.GI,
				})
			}
		}
	}
	if state.Captured[shogi.TurnFirst].KI > 0 {
		for _, d := range []*shogi.Position{
			shogi.Pos(+0, -1),
			shogi.Pos(-1, +0),
			shogi.Pos(+1, +0),
			shogi.Pos(-1, +1),
			shogi.Pos(+0, +1),
			shogi.Pos(+1, +1),
		} {
			if (targetFile+d.File > 0 && targetFile+d.File < 10 && targetRank+d.Rank > 0 && targetRank+d.Rank < 10) &&
				state.GetBoardPiece(targetFile+d.File, targetRank+d.Rank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+d.File, targetRank+d.Rank),
					Piece: shogi.KI,
				})
			}
		}
	}
	if state.Captured[shogi.TurnFirst].KA > 0 {
		for i := 1; targetFile+i < 10 && targetRank+i < 10; i++ {
			if state.GetBoardPiece(targetFile+i, targetRank+i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+i, targetRank+i),
					Piece: shogi.KA,
				})
			} else {
				break
			}
		}
		for i := 1; targetFile+i < 10 && targetRank-i > 0; i++ {
			if state.GetBoardPiece(targetFile+i, targetRank-i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+i, targetRank-i),
					Piece: shogi.KA,
				})
			} else {
				break
			}
		}
		for i := 1; targetFile-i > 0 && targetRank+i < 10; i++ {
			if state.GetBoardPiece(targetFile-i, targetRank+i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile-i, targetRank+i),
					Piece: shogi.KA,
				})
			} else {
				break
			}
		}
		for i := 1; targetFile-i > 0 && targetRank-i > 0; i++ {
			if state.GetBoardPiece(targetFile-i, targetRank-i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile-i, targetRank-i),
					Piece: shogi.KA,
				})
			} else {
				break
			}
		}
	}
	if state.Captured[shogi.TurnFirst].HI > 0 {
		for i := 1; targetFile+i < 10; i++ {
			if state.GetBoardPiece(targetFile+i, targetRank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile+i, targetRank),
					Piece: shogi.HI,
				})
			} else {
				break
			}
		}
		for i := 1; targetRank+i < 10; i++ {
			if state.GetBoardPiece(targetFile, targetRank+i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile, targetRank+i),
					Piece: shogi.HI,
				})
			} else {
				break
			}
		}
		for i := 1; targetFile-i > 0; i++ {
			if state.GetBoardPiece(targetFile-i, targetRank) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile-i, targetRank),
					Piece: shogi.HI,
				})
			} else {
				break
			}
		}
		for i := 1; targetRank-i > 0; i++ {
			if state.GetBoardPiece(targetFile, targetRank-i) == nil {
				results = append(results, &shogi.Move{
					Turn:  shogi.TurnFirst,
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(targetFile, targetRank-i),
					Piece: shogi.HI,
				})
			} else {
				break
			}
		}
	}
	return results
}

func simulate(state *shogi.State, move *shogi.Move) *shogi.State {
	// copy board state and captured
	s := state.Clone()
	// move, or use captured piece
	if move.Src.File > 0 && move.Src.Rank > 0 {
		bp := s.GetBoardPiece(move.Dst.File, move.Dst.Rank)
		if bp != nil {
			s.Captured[move.Turn].AddPieces(bp.Piece)
		}
		s.SetBoardPiece(move.Src.File, move.Src.Rank, nil)
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &shogi.BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
	} else {
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &shogi.BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
		s.Captured[move.Turn].SubPieces(move.Piece)
	}
	return s
}

func counterMoves(state *shogi.State) []*shogi.Move {
	results := []*shogi.Move{}
	move := state.Check(shogi.TurnFirst)
	if move == nil {
		return results
	}
	// move?
	for _, m := range state.CandidateMoves(shogi.TurnSecond) {
		check := simulate(state, m).Check(shogi.TurnFirst)
		if check == nil {
			results = append(results, m)
		}
	}
	// use captured pieces?
	if state.Captured[shogi.TurnSecond].Num() > 0 {
		available := []shogi.Piece{}
		if state.Captured[shogi.TurnSecond].FU > 0 {
			available = append(available, shogi.FU)
		}
		if state.Captured[shogi.TurnSecond].KY > 0 {
			available = append(available, shogi.KY)
		}
		if state.Captured[shogi.TurnSecond].KE > 0 {
			available = append(available, shogi.KE)
		}
		if state.Captured[shogi.TurnSecond].GI > 0 {
			available = append(available, shogi.GI)
		}
		if state.Captured[shogi.TurnSecond].KI > 0 {
			available = append(available, shogi.KI)
		}
		if state.Captured[shogi.TurnSecond].KA > 0 {
			available = append(available, shogi.KA)
		}
		if state.Captured[shogi.TurnSecond].HI > 0 {
			available = append(available, shogi.HI)
		}
		positions := []*shogi.Position{}
		target := *move.Dst
		for i := 1; target.File+i < 10; i++ {
			if state.GetBoardPiece(target.File+i, target.Rank) == nil {
				positions = append(positions, shogi.Pos(target.File+i, target.Rank))
			} else {
				break
			}
		}
		for i := 1; target.Rank+i < 10; i++ {
			if state.GetBoardPiece(target.File, target.Rank+i) == nil {
				positions = append(positions, shogi.Pos(target.File, target.Rank+i))
			} else {
				break
			}
		}
		for i := 1; target.File-i > 0; i++ {
			if state.GetBoardPiece(target.File-i, target.Rank) == nil {
				positions = append(positions, shogi.Pos(target.File-i, target.Rank))
			} else {
				break
			}
		}
		for i := 1; target.Rank-i > 0; i++ {
			if state.GetBoardPiece(target.File, target.Rank-i) == nil {
				positions = append(positions, shogi.Pos(target.File, target.Rank-i))
			} else {
				break
			}
		}
		for i := 1; target.File-i > 0 && target.Rank-i > 0; i++ {
			if state.GetBoardPiece(target.File-i, target.Rank-i) == nil {
				positions = append(positions, shogi.Pos(target.File-i, target.Rank-i))
			} else {
				break
			}
		}
		for i := 1; target.File-i > 0 && target.Rank+i < 10; i++ {
			if state.GetBoardPiece(target.File-i, target.Rank+i) == nil {
				positions = append(positions, shogi.Pos(target.File-i, target.Rank+i))
			} else {
				break
			}
		}
		for i := 1; target.File+i < 10 && target.Rank-i > 0; i++ {
			if state.GetBoardPiece(target.File+i, target.Rank-i) == nil {
				positions = append(positions, shogi.Pos(target.File+i, target.Rank-i))
			} else {
				break
			}
		}
		for i := 1; target.File+i < 10 && target.Rank+i < 10; i++ {
			if state.GetBoardPiece(target.File+i, target.Rank+i) == nil {
				positions = append(positions, shogi.Pos(target.File+i, target.Rank+i))
			} else {
				break
			}
		}

		for _, p := range positions {
			piece := available[0]
			// check duplicated FU
			if piece == shogi.FU {
				ok := true
				for rank := 1; rank < 10; rank++ {
					bp := state.GetBoardPiece(p.File, rank)
					if bp != nil && bp.Turn == shogi.TurnSecond && bp.Piece == shogi.FU {
						ok = false
						break
					}
				}
				if !ok {
					if len(available) > 1 {
						piece = available[1]
					} else {
						continue
					}
				}
			}
			move := &shogi.Move{
				Turn:  shogi.TurnSecond,
				Src:   shogi.Pos(0, 0),
				Dst:   shogi.Pos(p.File, p.Rank),
				Piece: piece,
			}
			check := simulate(state, move).Check(shogi.TurnFirst)
			// TODO isn't it wasted?
			if check == nil {
				results = append(results, move)
			}
		}
	}
	return results
}
