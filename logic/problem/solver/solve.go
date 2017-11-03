package solver

import (
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
func Solve(state *shogi.State) []*shogi.Move {
	return NewSolver(state).Solve()
}

// Solve method
func (s *Solver) Solve() []*shogi.Move {
	if s.state.Check(shogi.TurnFirst) != nil {
		return nil
	}
	answers := []*shogi.Move{}
	for _, m := range candidates(s.state) {
		// log.Println(s.state.MoveString(m))
		state := simulate(s.state, m, shogi.TurnFirst)
		if IsCheckmate(state, shogi.TurnFirst) {
			answers = append(answers, m)
		}
	}
	return answers
}

// IsCheckmate function
func IsCheckmate(state *shogi.State, turn shogi.Turn) bool {
	move := state.Check(turn)
	if move != nil {
		checkMate := true
		// move?
		for _, m := range state.CandidateMoves(!turn) {
			check := simulate(state, m, !turn).Check(turn)
			if check == nil {
				checkMate = false
				break
			}
		}
		// use captured pieces?
		if checkMate && state.Captured[!turn].Num() > 0 {
			available := []shogi.Piece{}
			if state.Captured[!turn].FU > 0 {
				available = append(available, shogi.FU)
			}
			if state.Captured[!turn].KY > 0 {
				available = append(available, shogi.KY)
			}
			if state.Captured[!turn].KE > 0 {
				available = append(available, shogi.KE)
			}
			if state.Captured[!turn].GI > 0 {
				available = append(available, shogi.GI)
			}
			if state.Captured[!turn].KI > 0 {
				available = append(available, shogi.KI)
			}
			if state.Captured[!turn].KA > 0 {
				available = append(available, shogi.KA)
			}
			if state.Captured[!turn].HI > 0 {
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
				if piece == shogi.FU {
					ok := true
					for rank := 1; rank < 10; rank++ {
						bp := state.GetBoardPiece(p.File, rank)
						if bp != nil && bp.Turn != turn && bp.Piece == shogi.FU {
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
					Src:   shogi.Pos(0, 0),
					Dst:   shogi.Pos(p.File, p.Rank),
					Piece: piece,
				}
				check := simulate(state, move, !turn).Check(turn)
				// TODO isn't it wasted?
				if check == nil {
					checkMate = false
					break
				}
			}
		}
		if checkMate {
			return true
		}
	}
	return false
}

func simulate(state *shogi.State, move *shogi.Move, turn shogi.Turn) *shogi.State {
	// copy board state and captured
	s := state.Clone()
	// move, or use captured piece
	if move.Src.File > 0 && move.Src.Rank > 0 {
		bp := s.GetBoardPiece(move.Dst.File, move.Dst.Rank)
		if bp != nil {
			s.Captured[turn].AddPieces(bp.Piece)
		}
		s.SetBoardPiece(move.Src.File, move.Src.Rank, nil)
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &shogi.BoardPiece{
			Turn:  turn,
			Piece: move.Piece,
		})
	} else {
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &shogi.BoardPiece{
			Turn:  turn,
			Piece: move.Piece,
		})
		s.Captured[turn].SubPieces(move.Piece)
	}
	return s
}

func candidates(state *shogi.State) []*shogi.Move {
	results := []*shogi.Move{}
	for _, move := range state.CandidateMoves(shogi.TurnFirst) {
		if simulate(state, move, shogi.TurnFirst).Check(shogi.TurnFirst) != nil {
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
		if state.GetBoardPiece(targetFile, targetRank+1) == nil {
			results = append(results, &shogi.Move{
				Turn:  shogi.TurnFirst,
				Src:   shogi.Pos(0, 0),
				Dst:   shogi.Pos(targetFile, targetRank+1),
				Piece: shogi.FU,
			})
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
					Dst:   shogi.Pos(targetFile+i, targetRank+i),
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
