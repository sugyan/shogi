package generator

import (
	"math/rand"

	"github.com/sugyan/shogi"
)

// previous states of first turn's player
func candidatePrevStatesF(state *shogi.State, pp *posPiece) []*shogi.State {
	candidates := []*posPiece{}
	switch pp.piece {
	case shogi.FU:
		ok := true
		for i := 0; i < 9; i++ {
			b := state.GetBoard(pp.pos.File, i+1)
			if b != nil && b.Turn == shogi.TurnBlack && b.Piece == shogi.FU {
				ok = false
				break
			}
		}
		if ok {
			candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), pp.piece})
		}
	case shogi.KY:
		candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		for i := 1; pp.pos.Rank+i < 10; i++ {
			if state.GetBoard(pp.pos.File, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), pp.piece})
			} else {
				break
			}
		}
	case shogi.KE:
		candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+2), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+2), pp.piece})
	case shogi.GI:
		candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank-1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank-1), pp.piece})
	case shogi.TO, shogi.NY, shogi.NK, shogi.NG, shogi.KI:
		if pp.piece == shogi.KI {
			candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		}
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank-1), pp.piece})
		if pp.pos.Rank <= 3 {
			switch pp.piece {
			case shogi.TO:
				ok := true
				for i := 0; i < 9; i++ {
					b := state.GetBoard(pp.pos.File, i+1)
					if b != nil && b.Turn == shogi.TurnBlack && b.Piece == shogi.FU {
						ok = false
						break
					}
				}
				if ok {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), shogi.FU})
				}
			case shogi.NY:
				for i := 1; pp.pos.Rank+i < 10; i++ {
					if state.GetBoard(pp.pos.File, pp.pos.Rank+i) == nil {
						candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), shogi.KY})
					} else {
						break
					}
				}
			case shogi.NK:
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+2), shogi.KE})
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+2), shogi.KE})
			case shogi.NG:
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+1), shogi.GI})
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+1), shogi.GI})
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), shogi.GI})
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank-1), shogi.GI})
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank-1), shogi.GI})
			}
		}
	case shogi.UM:
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank-1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank), pp.piece})
		fallthrough
	case shogi.KA:
		if pp.piece == shogi.KA {
			candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		}
		for i := 1; pp.pos.File-i > 0 && pp.pos.Rank-i > 0; i++ {
			if state.GetBoard(pp.pos.File-i, pp.pos.Rank-i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank-i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank-i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File-i > 0 && pp.pos.Rank+i < 10; i++ {
			if state.GetBoard(pp.pos.File-i, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank+i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank+i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File+i < 10 && pp.pos.Rank-i > 0; i++ {
			if state.GetBoard(pp.pos.File+i, pp.pos.Rank-i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank-i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank-i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File+i < 10 && pp.pos.Rank+i < 10; i++ {
			if state.GetBoard(pp.pos.File+i, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank+i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank+i), shogi.KA})
				}
			} else {
				break
			}
		}
	case shogi.RY:
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank-1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-1, pp.pos.Rank+1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank-1), pp.piece})
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+1, pp.pos.Rank+1), pp.piece})
		fallthrough
	case shogi.HI:
		if pp.piece == shogi.HI {
			candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		}
		for i := 1; pp.pos.File+i < 10; i++ {
			if state.GetBoard(pp.pos.File+i, pp.pos.Rank) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File-i > 0; i++ {
			if state.GetBoard(pp.pos.File-i, pp.pos.Rank) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.Rank+i < 10; i++ {
			if state.GetBoard(pp.pos.File, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.Rank-i > 0; i++ {
			if state.GetBoard(pp.pos.File, pp.pos.Rank-i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank-i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank-i), shogi.HI})
				}
			} else {
				break
			}
		}
	}

	states := []*shogi.State{}
	for _, i := range rand.Perm(len(candidates)) {
		c := candidates[i]
		if (c.pos.File == -1 && c.pos.Rank == -1) ||
			c.pos.File > 0 && c.pos.File < 10 && c.pos.Rank > 0 && c.pos.Rank < 10 {
			s := state.Clone()
			if c.pos.File == -1 && c.pos.Rank == -1 {
				s.Captured[shogi.TurnBlack].Add(pp.piece)
			} else {
				if state.GetBoard(c.pos.File, c.pos.Rank) == nil {
					s.SetBoard(c.pos.File, c.pos.Rank, &shogi.BoardPiece{
						Turn:  shogi.TurnBlack,
						Piece: c.piece,
					})
				} else {
					continue
				}
			}
			s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
			states = append(states, s)
		}
	}
	return states
}

// previous states of second turn's player
func candidatePrevStatesS(state *shogi.State, pp *posPiece, targetPos shogi.Position) []*shogi.State {
	states := []*shogi.State{}
	if pp.piece != shogi.OU {
		dFile, dRank := pp.pos.File-targetPos.File, pp.pos.Rank-targetPos.Rank
		switch abs(dFile) {
		case 0:
		case 1:
			if dRank > 2 || dRank < -1 {
				return states
			}
		default:
			if dRank != 0 {
				return states
			}
		}
	}

	available := []shogi.Piece{nil}
	if state.Captured[shogi.TurnWhite].FU > 0 {
		available = append(available, shogi.FU)
		available = append(available, shogi.TO)
	}
	if state.Captured[shogi.TurnWhite].KY > 0 {
		available = append(available, shogi.KY)
		available = append(available, shogi.NY)
	}
	if state.Captured[shogi.TurnWhite].KE > 0 {
		available = append(available, shogi.KE)
		available = append(available, shogi.NK)
	}
	if state.Captured[shogi.TurnWhite].GI > 0 {
		available = append(available, shogi.GI)
		available = append(available, shogi.NG)
	}
	if state.Captured[shogi.TurnWhite].KI > 0 {
		available = append(available, shogi.KI)
	}
	if state.Captured[shogi.TurnWhite].KA > 0 {
		available = append(available, shogi.KA)
		available = append(available, shogi.UM)
	}
	if state.Captured[shogi.TurnWhite].HI > 0 {
		available = append(available, shogi.HI)
		available = append(available, shogi.RY)
	}

	if pp.piece == shogi.OU {
		for _, p := range []shogi.Position{
			shogi.Pos(-1, -1), shogi.Pos(-1, +0), shogi.Pos(-1, +1),
			shogi.Pos(+0, -1), shogi.Pos(+0, +1),
			shogi.Pos(+1, -1), shogi.Pos(+1, +0), shogi.Pos(+1, +1),
		} {
			file, rank := pp.pos.File+p.File, pp.pos.Rank+p.Rank
			if !(file > 0 && file < 10 && rank > 0 && rank < 10) {
				continue
			}
			b := state.GetBoard(file, rank)
			if b != nil {
				continue
			}
			for _, piece := range available {
				s := state.Clone()
				s.SetBoard(file, rank, &shogi.BoardPiece{
					Turn:  shogi.TurnWhite,
					Piece: shogi.OU,
				})
				if piece == nil {
					s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
				} else {
					ok := true
					switch piece {
					case shogi.FU:
						if pp.pos.Rank <= 1 {
							ok = false
							break
						}
						if piece == shogi.FU {
							for i := 0; i < 9; i++ {
								bpf := state.GetBoard(pp.pos.File, i+1)
								if bpf != nil && bpf.Turn == shogi.TurnBlack && bpf.Piece == shogi.FU {
									ok = false
									break
								}
							}
						}
					case shogi.KY:
						if pp.pos.Rank <= 1 {
							ok = false
							break
						}
					case shogi.KE:
						if pp.pos.Rank <= 2 {
							ok = false
							break
						}
					}
					if !ok {
						continue
					}
					s.SetBoard(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
						Turn:  shogi.TurnBlack,
						Piece: piece,
					})
					s.Captured[shogi.TurnWhite].Sub(piece)
				}
				states = append(states, s)
			}
		}
	} else {
		if pp.pos.Rank == 1 && (pp.piece == shogi.FU || pp.piece == shogi.KY) {
			return states
		}
		if pp.pos.Rank <= 2 && pp.piece == shogi.KE {
			return states
		}
		prevPositions := []shogi.Position{}
		// TODO: promotion...?
		switch pp.piece {
		case shogi.FU:
			if pp.pos.Rank > 2 && state.GetBoard(pp.pos.File, pp.pos.Rank-1) == nil {
				prevPositions = append(prevPositions, shogi.Pos(pp.pos.File, pp.pos.Rank-1))
			}
		case shogi.KY:
			for i := 1; pp.pos.Rank-i > 0; i++ {
				if state.GetBoard(pp.pos.File, pp.pos.Rank-i) == nil {
					prevPositions = append(prevPositions, shogi.Pos(pp.pos.File, pp.pos.Rank-i))
				} else {
					break
				}
			}
		case shogi.KE:
			for _, d := range []shogi.Position{shogi.Pos(-1, -2), shogi.Pos(+1, -2)} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.GI:
			for _, d := range []shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(+0, -1),
				shogi.Pos(+1, -1),
				shogi.Pos(-1, +1),
				shogi.Pos(+1, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.TO, shogi.NY, shogi.NK, shogi.NG, shogi.KI:
			for _, d := range []shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(+0, -1),
				shogi.Pos(+1, -1),
				shogi.Pos(-1, +0),
				shogi.Pos(+1, +0),
				shogi.Pos(+0, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.UM:
			for _, d := range []shogi.Position{
				shogi.Pos(+0, -1),
				shogi.Pos(+0, +1),
				shogi.Pos(-1, +0),
				shogi.Pos(+1, +0),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
			fallthrough
		case shogi.KA:
			for i := 1; pp.pos.File-i > 0 && pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank-i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File-i > 0 && pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank+i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10 && pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank-i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10 && pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank+i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
		case shogi.RY:
			for _, d := range []shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(-1, +1),
				shogi.Pos(+1, -1),
				shogi.Pos(+1, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
			fallthrough
		case shogi.HI:
			for i := 1; pp.pos.File-i > 0; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File, pp.pos.Rank-i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File, pp.pos.Rank+i
				if state.GetBoard(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
		}
		for _, prevPos := range prevPositions {
			for _, piece := range available {
				s := state.Clone()
				s.SetBoard(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnBlack,
					Piece: piece,
				})
				s.SetBoard(prevPos.File, prevPos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnWhite,
					Piece: pp.piece,
				})
				s.Captured[shogi.TurnWhite].Sub(piece)
				states = append(states, s)
			}
		}
		// put a captured piece
		{
			s := state.Clone()
			s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
			s.Captured[shogi.TurnWhite].Add(pp.piece)
			states = append(states, s)
		}
	}
	return states
}
