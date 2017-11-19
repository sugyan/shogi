package generator

import (
	"math/rand"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver"
)

// Problem interface
type Problem interface {
	Steps() int
}

type problemType struct {
	n int
}

func (p *problemType) Steps() int {
	return p.n
}

// ProblemType variables
var (
	ProblemType1 = &problemType{1}
	ProblemType3 = &problemType{3}
)

type posPiece struct {
	pos   *shogi.Position
	piece shogi.Piece
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type generator struct {
	steps  int
	solver *solver.Solver
}

// Generate function
func Generate(pType Problem) *shogi.State {
	// TODO: timeout?
	generator := &generator{
		steps:  pType.Steps(),
		solver: solver.NewSolver(3),
	}
	return generator.generate()
}

func (g *generator) generate() *shogi.State {
	for {
		var state *shogi.State
		// random generate
		for {
			state = random()
			if g.solver.IsCheckmate(state) {
				break
			}
		}
		// reduce pieces
		g.cut(state)
		// rewind and check
		for _, s := range g.rewind(state, shogi.TurnFirst) {
			switch g.steps {
			case 1:
				if g.checkSolvable(s) {
					return g.cleanup(s)
				}
			case 3:
				states := g.rewind(s, shogi.TurnSecond)
				for _, i := range rand.Perm(len(states)) {
					if i > 5 {
						break
					}
					ss := states[i]
					if g.checkSolvable(ss) {
						return g.cleanup(ss)
					}
				}
			}
		}
	}
}

func (g *generator) rewind(state *shogi.State, turn shogi.Turn) []*shogi.State {
	var (
		targetPos *shogi.Position
		posPieces []*posPiece
	)
	// search pieces
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil && bp.Turn == turn {
				posPieces = append(posPieces, &posPiece{
					pos:   shogi.Pos(file, rank),
					piece: bp.Piece,
				})
				if bp.Piece == shogi.OU {
					targetPos = shogi.Pos(file, rank)
				}
			}
		}
	}

	results := []*shogi.State{}
	switch turn {
	case shogi.TurnFirst:
		for _, i := range rand.Perm(len(posPieces)) {
			pp := posPieces[i]
			candidates := candidatePrevStatesF(state, pp)
			for _, j := range rand.Perm(len(candidates)) {
				s := candidates[j]
				if s.Check(shogi.TurnFirst) == nil {
					results = append(results, s)
				}
			}
		}
	case shogi.TurnSecond:
		for _, pp := range posPieces {
			candidates := candidatePrevStatesS(state, pp, targetPos)
			for _, i := range rand.Perm(len(candidates)) {
				s := candidates[i]
				if s.Check(shogi.TurnFirst) != nil {
					results = append(results, g.rewind(s, shogi.TurnFirst)...)
					break
				}
			}
		}
	}
	return results
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
	for originalPiece, num := range map[shogi.Piece]int{
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
			var (
				piece = originalPiece
				turn  = shogi.TurnFirst
			)
			if rand.Intn(3) > 0 {
				turn = shogi.TurnSecond
			}
			switch originalPiece {
			case shogi.FU:
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
			case shogi.KY:
				if turn == shogi.TurnFirst {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NY
					} else if p.Rank <= 1 {
						turn = shogi.TurnSecond
					}
				}
			case shogi.KE:
				if turn == shogi.TurnFirst {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NK
					} else if p.Rank <= 2 {
						turn = shogi.TurnSecond
					}
				}
			case shogi.GI:
				if turn == shogi.TurnFirst {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NG
					}
				}
			case shogi.KI:
			case shogi.KA:
				if turn == shogi.TurnFirst && p.Rank <= 6 {
					if rand.Intn(2) == 0 {
						piece = shogi.UM
					}
				}
			case shogi.HI:
				if turn == shogi.TurnFirst && p.Rank <= 6 {
					if rand.Intn(2) == 0 {
						piece = shogi.RY
					}
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

// TODO: fix bug...
func (g *generator) cut(state *shogi.State) {
	positions := []*shogi.Position{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-i, j+1
			bp := state.GetBoardPiece(file, rank)
			if bp != nil {
				if bp.Turn == shogi.TurnFirst {
					positions = append(positions, shogi.Pos(file, rank))
				}
			}
		}
	}
	for _, i := range rand.Perm(len(positions)) {
		if rand.Intn(5) == 0 {
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

// previous states of first turn's player
func candidatePrevStatesF(state *shogi.State, pp *posPiece) []*shogi.State {
	candidates := []*posPiece{}
	switch pp.piece {
	case shogi.FU:
		ok := true
		for i := 0; i < 9; i++ {
			bp := state.GetBoardPiece(pp.pos.File, i+1)
			if bp != nil && bp.Turn == shogi.TurnFirst && bp.Piece == shogi.FU {
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
				ok := true
				for i := 0; i < 9; i++ {
					bp := state.GetBoardPiece(pp.pos.File, i+1)
					if bp != nil && bp.Turn == shogi.TurnFirst && bp.Piece == shogi.FU {
						ok = false
						break
					}
				}
				if ok {
					candidates = append(candidates, &posPiece{shogi.Pos(pp.pos.File, pp.pos.Rank+1), shogi.FU})
				}
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

// previous states of second turn's player
func candidatePrevStatesS(state *shogi.State, pp *posPiece, targetPos *shogi.Position) []*shogi.State {
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
	if state.Captured[shogi.TurnSecond].FU > 0 {
		available = append(available, shogi.FU)
		available = append(available, shogi.TO)
	}
	if state.Captured[shogi.TurnSecond].KY > 0 {
		available = append(available, shogi.KY)
		available = append(available, shogi.NY)
	}
	if state.Captured[shogi.TurnSecond].KE > 0 {
		available = append(available, shogi.KE)
		available = append(available, shogi.NK)
	}
	if state.Captured[shogi.TurnSecond].GI > 0 {
		available = append(available, shogi.GI)
		available = append(available, shogi.NG)
	}
	if state.Captured[shogi.TurnSecond].KI > 0 {
		available = append(available, shogi.KI)
	}
	if state.Captured[shogi.TurnSecond].KA > 0 {
		available = append(available, shogi.KA)
		available = append(available, shogi.UM)
	}
	if state.Captured[shogi.TurnSecond].HI > 0 {
		available = append(available, shogi.HI)
		available = append(available, shogi.RY)
	}

	if pp.piece == shogi.OU {
		for _, p := range []*shogi.Position{
			shogi.Pos(-1, -1), shogi.Pos(-1, +0), shogi.Pos(-1, +1),
			shogi.Pos(+0, -1), shogi.Pos(+0, +1),
			shogi.Pos(+1, -1), shogi.Pos(+1, +0), shogi.Pos(+1, +1),
		} {
			file, rank := pp.pos.File+p.File, pp.pos.Rank+p.Rank
			if !(file > 0 && file < 10 && rank > 0 && rank < 10) {
				continue
			}
			bp := state.GetBoardPiece(file, rank)
			if bp != nil {
				continue
			}
			for _, piece := range available {
				s := state.Clone()
				s.SetBoardPiece(file, rank, &shogi.BoardPiece{
					Turn:  shogi.TurnSecond,
					Piece: shogi.OU,
				})
				switch piece {
				case nil:
					s.SetBoardPiece(pp.pos.File, pp.pos.Rank, nil)
				default:
					ok := true
					if piece == shogi.FU {
						for i := 0; i < 9; i++ {
							bpf := state.GetBoardPiece(i+1, pp.pos.Rank)
							if bpf != nil && bpf.Turn == shogi.TurnSecond && bpf.Piece == shogi.FU {
								ok = false
								break
							}
						}
					}
					if !ok {
						continue
					}
					s.SetBoardPiece(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
						Turn:  shogi.TurnFirst,
						Piece: piece,
					})
					s.Captured[shogi.TurnSecond].SubPieces(piece)
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
		prevPositions := []*shogi.Position{}
		// TODO: promotion...?
		switch pp.piece {
		case shogi.FU:
			if pp.pos.Rank > 2 && state.GetBoardPiece(pp.pos.File, pp.pos.Rank-1) == nil {
				prevPositions = append(prevPositions, shogi.Pos(pp.pos.File, pp.pos.Rank-1))
			}
		case shogi.KY:
			for i := 1; pp.pos.Rank-i > 0; i++ {
				if state.GetBoardPiece(pp.pos.File, pp.pos.Rank-i) == nil {
					prevPositions = append(prevPositions, shogi.Pos(pp.pos.File, pp.pos.Rank-i))
				} else {
					break
				}
			}
		case shogi.KE:
			for _, d := range []*shogi.Position{shogi.Pos(-1, -2), shogi.Pos(+1, -2)} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.GI:
			for _, d := range []*shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(+0, -1),
				shogi.Pos(+1, -1),
				shogi.Pos(-1, +1),
				shogi.Pos(+1, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.TO, shogi.NY, shogi.NK, shogi.NG, shogi.KI:
			for _, d := range []*shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(+0, -1),
				shogi.Pos(+1, -1),
				shogi.Pos(-1, +0),
				shogi.Pos(+1, +0),
				shogi.Pos(+0, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
		case shogi.UM:
			for _, d := range []*shogi.Position{
				shogi.Pos(+0, -1),
				shogi.Pos(+0, +1),
				shogi.Pos(-1, +0),
				shogi.Pos(+1, +0),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
			fallthrough
		case shogi.KA:
			for i := 1; pp.pos.File-i > 0 && pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank-i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File-i > 0 && pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank+i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10 && pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank-i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10 && pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank+i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
		case shogi.RY:
			for _, d := range []*shogi.Position{
				shogi.Pos(-1, -1),
				shogi.Pos(-1, +1),
				shogi.Pos(+1, -1),
				shogi.Pos(+1, +1),
			} {
				file, rank := pp.pos.File+d.File, pp.pos.Rank+d.Rank
				if file > 0 && file < 10 && rank > 0 && rank < 10 && state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				}
			}
			fallthrough
		case shogi.HI:
			for i := 1; pp.pos.File-i > 0; i++ {
				file, rank := pp.pos.File-i, pp.pos.Rank
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.Rank-i > 0; i++ {
				file, rank := pp.pos.File, pp.pos.Rank-i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.File+i < 10; i++ {
				file, rank := pp.pos.File+i, pp.pos.Rank
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
			for i := 1; pp.pos.Rank+i < 10; i++ {
				file, rank := pp.pos.File, pp.pos.Rank+i
				if state.GetBoardPiece(file, rank) == nil {
					prevPositions = append(prevPositions, shogi.Pos(file, rank))
				} else {
					break
				}
			}
		}
		for _, prevPos := range prevPositions {
			for _, piece := range available {
				s := state.Clone()
				s.SetBoardPiece(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnFirst,
					Piece: piece,
				})
				s.SetBoardPiece(prevPos.File, prevPos.Rank, &shogi.BoardPiece{
					Turn:  shogi.TurnSecond,
					Piece: pp.piece,
				})
				s.Captured[shogi.TurnSecond].SubPieces(piece)
				states = append(states, s)
			}
		}
		// put a captured piece
		{
			s := state.Clone()
			s.SetBoardPiece(pp.pos.File, pp.pos.Rank, nil)
			s.Captured[shogi.TurnSecond].AddPieces(pp.piece)
			states = append(states, s)
		}
	}
	return states
}

func (g *generator) checkSolvable(state *shogi.State) bool {
	answers, length := g.solver.ValidAnswers(state)
	if len(answers) == 0 {
		return false
	}
	if length != g.steps {
		return false
	}
	// check catured pieces
	ok := false
	for _, answer := range answers {
		s := state.Clone()
		for _, move := range answer {
			s.Apply(move)
		}
		if s.Captured[shogi.TurnFirst].Num() == 0 {
			ok = true
		}
	}
	if !ok {
		return false
	}
	// check uniqueness
	switch g.steps {
	case 1:
		return len(answers) == 1
	case 3:
		a := answers[0][0]
		for _, answer := range answers {
			if *a != *answer[0] {
				return false
			}
		}
		return true
	}
	return false
}

func (g *generator) cleanup(state *shogi.State) *shogi.State {
	// remove unnecessary pieces
	{
		posPieces := []*posPiece{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				file, rank := 9-i, j+1
				bp := state.GetBoardPiece(file, rank)
				if bp != nil {
					if bp.Piece != shogi.OU {
						posPieces = append(posPieces, &posPiece{
							pos:   shogi.Pos(file, rank),
							piece: bp.Piece,
						})
					}
				}
			}
		}
		for _, i := range rand.Perm(len(posPieces)) {
			pp := posPieces[i]
			s := state.Clone()
			s.SetBoardPiece(pp.pos.File, pp.pos.Rank, nil)
			s.Captured[shogi.TurnSecond].AddPieces(pp.piece)
			if s.Check(shogi.TurnFirst) == nil && g.checkSolvable(s) {
				state.SetBoardPiece(pp.pos.File, pp.pos.Rank, nil)
				state.Captured[shogi.TurnSecond].AddPieces(pp.piece)
			}
		}
	}
	// reaplace TO, NY, NK, NG to KI or TO
	{
		posPieces := map[shogi.Turn][]*posPiece{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				file, rank := 9-i, j+1
				bp := state.GetBoardPiece(file, rank)
				if bp != nil {
					posPieces[bp.Turn] = append(posPieces[bp.Turn], &posPiece{
						pos:   shogi.Pos(file, rank),
						piece: bp.Piece,
					})
				}
			}
		}
		for _, turn := range []shogi.Turn{shogi.TurnSecond, shogi.TurnFirst} {
			for _, i := range rand.Perm(len(posPieces[turn])) {
				pp := posPieces[turn][i]
				switch pp.piece {
				case shogi.TO, shogi.NY, shogi.NK, shogi.NG:
					if state.Captured[shogi.TurnSecond].KI > 0 {
						state.SetBoardPiece(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
							Turn:  turn,
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
						if pp.piece == shogi.TO {
							continue
						}
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
		}
	}
	return state
}
