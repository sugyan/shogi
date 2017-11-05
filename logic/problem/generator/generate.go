package generator

import (
	"math/rand"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type generator struct {
	solver *solver.Solver
}

// Generate function
func Generate() *shogi.State {
	// TODO: timeout?
	generator := &generator{
		solver: solver.NewSolver(),
	}
	return generator.generate()
}

func (g *generator) generate() *shogi.State {
	for {
		var state *shogi.State
		for {
			state = random()
			if g.solver.IsCheckmate(state) {
				break
			}
		}
		g.cut(state)
		result := rewind(state)
		if result != nil {
			if g.checkSolvable(result) {
				return cleanup(result)
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func random() *shogi.State {
	s := shogi.NewState()
	targetFile := rand.Intn(4) + 1
	targetRank := rand.Intn(4) + 1
	positions := make([]*shogi.Position, 0, 81)
	for _, n := range rand.Perm(81) {
		file, rank := int(n/9)+1, n%9+1
		d := abs(file-targetFile) + abs(rank-targetRank)
		if d > 0 && d < 9 {
			positions = append(positions, shogi.Pos(file, rank))
		}
	}
	// target
	{
		s.SetBoardPiece(targetFile, targetRank, &shogi.BoardPiece{
			Turn:  shogi.TurnSecond,
			Piece: shogi.OU,
		})
		positions = positions[1:]
	}
	// other pieces
	for piece, num := range map[shogi.Piece]int{
		shogi.FU: 18,
		shogi.KY: 4,
		shogi.KE: 4,
		shogi.GI: 4,
		shogi.KI: 4,
		shogi.KA: 2,
		shogi.HI: 2,
	} {
		for i := 0; i < num; i++ {
			// shift position
			p := positions[0]
			positions = positions[1:]
			// set pieces
			turn := shogi.TurnFirst
			if rand.Intn(3) > 0 {
				turn = shogi.TurnSecond
			}
			if piece == shogi.FU {
				if turn == shogi.TurnFirst && p.Rank <= 3 {
					if p.Rank == 1 || rand.Intn(2) == 0 {
						piece = shogi.TO
					}
				} else {
					exist := false
					for i := 0; i < 9; i++ {
						bp := s.GetBoardPiece(p.File, i+1)
						if bp != nil && bp.Turn == turn && bp.Piece == shogi.FU {
							exist = true
							break
						}
					}
					if exist {
						s.Captured[shogi.TurnSecond].FU++
						continue
					}
				}
			}
			if piece == shogi.KY && turn == shogi.TurnFirst {
				if p.Rank <= 3 && rand.Intn(3) == 0 {
					piece = shogi.NY
				} else if p.Rank <= 1 {
					turn = shogi.TurnSecond
				}
			}
			if piece == shogi.KE && turn == shogi.TurnFirst {
				if p.Rank <= 3 && rand.Intn(3) == 0 {
					piece = shogi.NK
				} else if p.Rank <= 2 {
					turn = shogi.TurnSecond
				}
			}
			if piece == shogi.GI && turn == shogi.TurnFirst {
				if p.Rank <= 3 && rand.Intn(3) == 0 {
					piece = shogi.NG
				}
			}
			if piece == shogi.KA && turn == shogi.TurnFirst && p.Rank <= 6 {
				if rand.Intn(2) == 0 {
					piece = shogi.UM
				}
			}
			if piece == shogi.HI && turn == shogi.TurnFirst && p.Rank <= 6 {
				if rand.Intn(2) == 0 {
					piece = shogi.RY
				}
			}
			s.SetBoardPiece(p.File, p.Rank, &shogi.BoardPiece{
				Turn:  turn,
				Piece: piece,
			})
		}
	}
	return s
}

func (g *generator) cut(state *shogi.State) {
	positions := []*shogi.Position{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-i, j+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil {
				if bp.Piece != shogi.OU {
					positions = append(positions, shogi.Pos(file, rank))
				}
			}
		}
	}
	for _, i := range rand.Perm(len(positions)) {
		if rand.Intn(10) == 0 {
			continue
		}
		p := positions[i]
		s := state.Clone()
		s.SetBoardPiece(p.File, p.Rank, nil)
		if g.solver.IsCheckmate(s) {
			bp := state.GetBoardPiece(p.File, p.Rank)
			state.SetBoardPiece(p.File, p.Rank, nil)
			state.Captured[shogi.TurnSecond].AddPieces(bp.Piece)
		}
	}
}

type posPiece struct {
	pos   *shogi.Position
	piece shogi.Piece
}

func rewind(state *shogi.State) *shogi.State {
	posPieces := []*posPiece{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil && bp.Turn == shogi.TurnFirst {
				posPieces = append(posPieces, &posPiece{
					pos:   shogi.Pos(file, rank),
					piece: bp.Piece,
				})
			}
		}
	}

	for _, i := range rand.Perm(len(posPieces)) {
		pp := posPieces[i]
		for _, s := range candidatePrevStates(state, pp) {
			check := s.Check(shogi.TurnFirst)
			if check == nil {
				return s
			}
		}
	}
	return nil
}

func candidatePrevStates(state *shogi.State, pp *posPiece) []*shogi.State {
	candidates := []*posPiece{}
	switch pp.piece {
	case shogi.FU:
		candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), pp.piece})
	case shogi.KY:
		candidates = append(candidates, &posPiece{shogi.Pos(-1, -1), pp.piece})
		for i := 1; pp.pos.Rank+i < 10; i++ {
			if state.GetBoardPiece(pp.pos.File, pp.pos.Rank+i) == nil {
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
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), shogi.FU})
			case shogi.NY:
				for i := 1; pp.pos.Rank+i < 10; i++ {
					if state.GetBoardPiece(pp.pos.File, pp.pos.Rank+i) == nil {
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
			if state.GetBoardPiece(pp.pos.File-i, pp.pos.Rank-i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank-i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank-i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File-i > 0 && pp.pos.Rank+i < 10; i++ {
			if state.GetBoardPiece(pp.pos.File-i, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank+i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank+i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File+i < 10 && pp.pos.Rank-i > 0; i++ {
			if state.GetBoardPiece(pp.pos.File+i, pp.pos.Rank-i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank-i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.UM {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank-i), shogi.KA})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File+i < 10 && pp.pos.Rank+i < 10; i++ {
			if state.GetBoardPiece(pp.pos.File+i, pp.pos.Rank+i) == nil {
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
			if state.GetBoardPiece(pp.pos.File+i, pp.pos.Rank) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File+i, pp.pos.Rank), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.File-i > 0; i++ {
			if state.GetBoardPiece(pp.pos.File-i, pp.pos.Rank) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File-i, pp.pos.Rank), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.Rank+i < 10; i++ {
			if state.GetBoardPiece(pp.pos.File, pp.pos.Rank+i) == nil {
				candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), pp.piece})
				if pp.pos.Rank <= 3 && pp.piece == shogi.RY {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+i), shogi.HI})
				}
			} else {
				break
			}
		}
		for i := 1; pp.pos.Rank-i > 0; i++ {
			if state.GetBoardPiece(pp.pos.File, pp.pos.Rank-i) == nil {
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
				s.Captured[shogi.TurnFirst].AddPieces(pp.piece)
			} else {
				bp := state.GetBoardPiece(c.pos.File, c.pos.Rank)
				if bp == nil {
					s.SetBoardPiece(c.pos.File, c.pos.Rank, &shogi.BoardPiece{
						Turn:  shogi.TurnFirst,
						Piece: c.piece,
					})
				} else {
					continue
				}
			}
			s.SetBoardPiece(pp.pos.File, pp.pos.Rank, nil)
			states = append(states, s)
		}
	}
	return states
}

func (g *generator) checkSolvable(state *shogi.State) bool {
	positions := []*shogi.Position{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-i, j+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil {
				if bp.Piece != shogi.OU {
					positions = append(positions, shogi.Pos(file, rank))
				}
			}
		}
	}
	for _, i := range rand.Perm(len(positions)) {
		p := positions[i]
		s := state.Clone()
		bp := state.GetBoardPiece(p.File, p.Rank)
		s.SetBoardPiece(p.File, p.Rank, nil)
		s.Captured[shogi.TurnSecond].AddPieces(bp.Piece)
		answers := g.solver.Solve(s, 0)
		if len(answers) >= 1 {
			bp := state.GetBoardPiece(p.File, p.Rank)
			state.SetBoardPiece(p.File, p.Rank, nil)
			state.Captured[shogi.TurnSecond].AddPieces(bp.Piece)
		}
	}
	answers := g.solver.Solve(state, 0)
	// TODO check if it is too easy
	return len(answers) == 1 && len(answers[0]) == 1
}

func cleanup(state *shogi.State) *shogi.State {
	posPieces := []*posPiece{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-i, j+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil && bp.Turn == shogi.TurnFirst {
				posPieces = append(posPieces, &posPiece{
					pos:   shogi.Pos(file, rank),
					piece: bp.Piece,
				})
			}
		}
	}
	// reaplace TO, NY, NK, NG to KI or TO
	for _, i := range rand.Perm(len(posPieces)) {
		pp := posPieces[i]
		switch pp.piece {
		case shogi.TO:
			fallthrough
		case shogi.NY:
			fallthrough
		case shogi.NK:
			fallthrough
		case shogi.NG:
			if state.Captured[shogi.TurnSecond].KI > 0 {
				state.SetBoardPiece(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnFirst,
					Piece: shogi.KI,
				})
				state.Captured[shogi.TurnSecond].KI--
				switch pp.piece {
				case shogi.TO:
					state.Captured[shogi.TurnSecond].FU++
				case shogi.NY:
					state.Captured[shogi.TurnSecond].KY++
				case shogi.NK:
					state.Captured[shogi.TurnSecond].KE++
				case shogi.NG:
					state.Captured[shogi.TurnSecond].GI++
				}
			} else if state.Captured[shogi.TurnSecond].FU > 0 {
				state.SetBoardPiece(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnFirst,
					Piece: shogi.TO,
				})
				state.Captured[shogi.TurnSecond].FU--
				switch pp.piece {
				case shogi.NY:
					state.Captured[shogi.TurnSecond].KY++
				case shogi.NK:
					state.Captured[shogi.TurnSecond].KE++
				case shogi.NG:
					state.Captured[shogi.TurnSecond].GI++
				}
			}
		}
	}
	return state
}
