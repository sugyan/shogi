package problem

import (
	"github.com/sugyan/shogi"
)

// MoveState type
type MoveState struct {
	Move  *shogi.Move
	State *shogi.State
}

// Candidates function
func Candidates(state *shogi.State, turn shogi.Turn) []*MoveState {
	results := []*MoveState{}
	target := searchTarget(state)
	// by moving pieces
	for _, m := range state.CandidateMoves(turn) {
		s := state.Clone()
		s.Apply(m)
		check := s.Check(shogi.TurnBlack) != nil
		if (turn == shogi.TurnBlack && check) || (turn == shogi.TurnWhite && !check) {
			results = append(results, &MoveState{m, s})
		}
	}
	// by dropping captured pieces
	switch turn {
	case shogi.TurnBlack:
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
					results = append(results, &MoveState{m, s})
				}
			}
		}
	case shogi.TurnWhite:
		if state.Captured[shogi.TurnWhite].Num() == 0 {
			break
		}
		positions := []shogi.Position{}
		for _, direction := range []shogi.Position{
			shogi.Pos(-1, -1), shogi.Pos(-1, +1), shogi.Pos(+1, -1), shogi.Pos(+1, +1),
			shogi.Pos(-1, +0), shogi.Pos(+1, +0), shogi.Pos(+0, -1), shogi.Pos(+0, +1),
		} {
			candidates := []shogi.Position{}
			for i := 1; ; i++ {
				file := target.File + i*direction.File
				rank := target.Rank + i*direction.Rank
				if !(file > 0 && file < 10 && rank > 0 && rank < 10) {
					break
				}
				b := state.GetBoard(file, rank)
				if b == nil {
					candidates = append(candidates, shogi.Pos(file, rank))
				} else {
					d := direction.File * direction.Rank
					if d == 0 && b.Turn == shogi.TurnBlack &&
						((b.Piece == shogi.HI || b.Piece == shogi.RY) ||
							(direction.Rank == +1 && b.Piece == shogi.KY)) {
						positions = append(positions, candidates...)
					}
					if d != 0 && b.Turn == shogi.TurnBlack &&
						(b.Piece == shogi.KA || b.Piece == shogi.UM) {
						positions = append(positions, candidates...)
					}
					break
				}
			}
		}
		if len(positions) > 0 {
			available := state.Captured[shogi.TurnWhite].Available()
			reachable := map[shogi.Position]struct{}{}
			for _, m := range state.CandidateMoves(shogi.TurnWhite) {
				reachable[m.Dst] = struct{}{}
			}
			for _, p := range positions {
				for _, piece := range available {
					// check duplicated FU
					if piece == shogi.FU {
						ok := true
						for rank := 1; rank < 10; rank++ {
							b := state.GetBoard(p.File, rank)
							if b != nil && b.Turn == shogi.TurnWhite && b.Piece == shogi.FU {
								ok = false
								break
							}
						}
						if !ok {
							continue
						}
					}
					m := &shogi.Move{
						Turn:  shogi.TurnWhite,
						Src:   shogi.Pos(0, 0),
						Dst:   shogi.Pos(p.File, p.Rank),
						Piece: piece,
					}
					s := state.Clone()
					s.Apply(m)
					check := s.Check(shogi.TurnBlack)
					if check == nil {
						results = append(results, &MoveState{m, s})
					}
					// FIXME: temporary improvement...
					if _, ok := reachable[p]; !ok {
						break
					}
				}
			}
		}
	}
	return results
}

func searchTarget(state *shogi.State) shogi.Position {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := state.Board[i][j]
			if bp != nil && bp.Piece == shogi.OU && bp.Turn == shogi.TurnWhite {
				return shogi.Pos(9-j, i+1)
			}
		}
	}
	return shogi.Pos(0, 0)
}
