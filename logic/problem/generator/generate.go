package generator

import (
	"math/rand"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem"
	"github.com/sugyan/shogi/logic/problem/solver"
	"github.com/sugyan/shogi/logic/problem/solver/node"
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

// Type variables
var (
	Type1 = &problemType{1}
	Type3 = &problemType{3}
)

type posPiece struct {
	pos   shogi.Position
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
		steps: pType.Steps(),
	}
	return generator.generate()
}

func (g *generator) generate() *shogi.State {
	for {
		var state *shogi.State
		// random generate
		for {
			state = random()
			if isCheckmate(state) {
				break
			}
		}
		// reduce pieces
		g.cut(state)
		// rewind and check
		for _, s := range g.rewind(state, shogi.TurnBlack) {
			switch g.steps {
			case 1:
				if isValidProblem(s, g.steps) {
					// TODO: evaluate
					g.cleanup(s)
					return s
				}
			case 3:
				states := g.rewind(s, shogi.TurnWhite)
				for _, i := range rand.Perm(len(states)) {
					if i > 5 {
						break
					}
					s := states[i]
					if isValidProblem(s, g.steps) {
						g.cleanup(s)
						return s
					}
				}
			}
		}
	}
}

func isCheckmate(state *shogi.State) bool {
	if state.Check(shogi.TurnBlack) == nil {
		return false
	}
	candidates := problem.Candidates(state, shogi.TurnWhite)
	if len(candidates) == 0 {
		return true
	}
	// check wasted
	isCaptured := true
	for _, ms := range candidates {
		if !ms.Move.Src.IsCaptured() {
			isCaptured = false
			break
		}
	}
	if isCaptured {
		dst := map[shogi.Position]*shogi.Move{}
		for _, ms := range candidates {
			if _, exist := dst[ms.Move.Dst]; !exist {
				dst[ms.Move.Dst] = ms.Move
			}
		}
		result := true
		for _, move := range dst {
			s := state.Clone()
			s.Apply(move)
			root := solver.NewSolver().Search(s, 0)
			answer := solver.SearchBestAnswer(root)
			ok := false
			if len(answer) == 1 {
				for _, c := range root.Children() {
					if c.Move().Dst == move.Dst && c.Result() == node.ResultT {
						ok = true
						break
					}
				}
			}
			if !ok {
				result = false
				break
			}
		}
		return result
	}
	return false
}

func (g *generator) rewind(state *shogi.State, turn shogi.Turn) []*shogi.State {
	var (
		targetPos shogi.Position
		posPieces []*posPiece
	)
	// search pieces
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-j, i+1
			b := state.GetBoard(file, rank)
			if b != nil && b.Turn == turn {
				posPieces = append(posPieces, &posPiece{
					pos:   shogi.Pos(file, rank),
					piece: b.Piece,
				})
				if b.Piece == shogi.OU {
					targetPos = shogi.Pos(file, rank)
				}
			}
		}
	}

	results := []*shogi.State{}
	switch turn {
	case shogi.TurnBlack:
		for _, i := range rand.Perm(len(posPieces)) {
			pp := posPieces[i]
			candidates := candidatePrevStatesF(state, pp)
			for _, j := range rand.Perm(len(candidates)) {
				s := candidates[j]
				if s.Check(shogi.TurnBlack) == nil {
					results = append(results, s)
				}
			}
		}
	case shogi.TurnWhite:
		for _, pp := range posPieces {
			candidates := candidatePrevStatesS(state, pp, targetPos)
			for _, i := range rand.Perm(len(candidates)) {
				s := candidates[i]
				if s.Check(shogi.TurnBlack) != nil {
					results = append(results, g.rewind(s, shogi.TurnBlack)...)
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
	positions := make([]shogi.Position, 0, 81)
	for _, n := range rand.Perm(81) {
		file, rank := int(n/9)+1, n%9+1
		d := abs(file-targetFile) + abs(rank-targetRank)
		if d > 0 && d < 9 {
			positions = append(positions, shogi.Pos(file, rank))
		}
	}
	// target
	{
		s.SetBoard(targetFile, targetRank, &shogi.BoardPiece{
			Turn:  shogi.TurnWhite,
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
				turn  = shogi.TurnBlack
			)
			if rand.Intn(3) > 0 {
				turn = shogi.TurnWhite
			}
			switch originalPiece {
			case shogi.FU:
				if turn == shogi.TurnBlack && p.Rank <= 3 {
					if p.Rank == 1 || rand.Intn(2) == 0 {
						piece = shogi.TO
					}
				} else {
					exist := false
					for i := 0; i < 9; i++ {
						b := s.GetBoard(p.File, i+1)
						if b != nil && b.Turn == turn && b.Piece == shogi.FU {
							exist = true
							break
						}
					}
					if exist {
						s.Captured[shogi.TurnWhite].FU++
						continue
					}
				}
			case shogi.KY:
				if turn == shogi.TurnBlack {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NY
					} else if p.Rank <= 1 {
						turn = shogi.TurnWhite
					}
				}
			case shogi.KE:
				if turn == shogi.TurnBlack {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NK
					} else if p.Rank <= 2 {
						turn = shogi.TurnWhite
					}
				}
			case shogi.GI:
				if turn == shogi.TurnBlack {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NG
					}
				}
			case shogi.KI:
			case shogi.KA:
				if turn == shogi.TurnBlack && p.Rank <= 6 {
					if rand.Intn(2) == 0 {
						piece = shogi.UM
					}
				}
			case shogi.HI:
				if turn == shogi.TurnBlack && p.Rank <= 6 {
					if rand.Intn(2) == 0 {
						piece = shogi.RY
					}
				}
			}
			s.SetBoard(p.File, p.Rank, &shogi.BoardPiece{
				Turn:  turn,
				Piece: piece,
			})
		}
	}
	return s
}

func (g *generator) cut(state *shogi.State) {
	positions := []shogi.Position{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			file, rank := 9-i, j+1
			b := state.GetBoard(file, rank)
			if b != nil {
				if b.Piece != shogi.OU {
					positions = append(positions, shogi.Pos(file, rank))
				}
			}
		}
	}
	for _, i := range rand.Perm(len(positions)) {
		if rand.Intn(3) == 0 {
			continue
		}
		p := positions[i]
		s := state.Clone()
		s.SetBoard(p.File, p.Rank, nil)
		if isCheckmate(s) {
			b := state.GetBoard(p.File, p.Rank)
			state.SetBoard(p.File, p.Rank, nil)
			state.Captured[shogi.TurnWhite].Add(b.Piece)
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
				switch piece {
				case nil:
					s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
				default:
					ok := true
					if piece == shogi.FU {
						for i := 0; i < 9; i++ {
							bpf := state.GetBoard(i+1, pp.pos.Rank)
							if bpf != nil && bpf.Turn == shogi.TurnWhite && bpf.Piece == shogi.FU {
								ok = false
								break
							}
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

func isValidProblem(state *shogi.State, steps int) bool {
	// TODO: check if there are multiple answers
	root := solver.NewSolver().Search(state, steps+1)
	bestAnswer := solver.SearchBestAnswer(root)

	if len(bestAnswer) != steps {
		return false
	}
	switch steps {
	case 1:
		num := 0
		for _, c := range root.Children() {
			if c.Result() == node.ResultT {
				num++
			}
		}
		if num == 1 {
			return true
		}
	default:
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
				b := state.GetBoard(file, rank)
				if b != nil {
					if b.Piece != shogi.OU {
						posPieces = append(posPieces, &posPiece{
							pos:   shogi.Pos(file, rank),
							piece: b.Piece,
						})
					}
				}
			}
		}
		for _, i := range rand.Perm(len(posPieces)) {
			pp := posPieces[i]
			s := state.Clone()
			s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
			s.Captured[shogi.TurnWhite].Add(pp.piece)
			if s.Check(shogi.TurnBlack) == nil && isValidProblem(s, g.steps) {
				state.SetBoard(pp.pos.File, pp.pos.Rank, nil)
				state.Captured[shogi.TurnWhite].Add(pp.piece)
			}
		}
	}
	// reaplace TO, NY, NK, NG to KI or TO
	{
		posPieces := map[shogi.Turn][]*posPiece{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				file, rank := 9-i, j+1
				b := state.GetBoard(file, rank)
				if b != nil {
					posPieces[b.Turn] = append(posPieces[b.Turn], &posPiece{
						pos:   shogi.Pos(file, rank),
						piece: b.Piece,
					})
				}
			}
		}
		for _, turn := range []shogi.Turn{shogi.TurnWhite, shogi.TurnBlack} {
			for _, i := range rand.Perm(len(posPieces[turn])) {
				pp := posPieces[turn][i]
				switch pp.piece {
				case shogi.TO, shogi.NY, shogi.NK, shogi.NG:
					if state.Captured[shogi.TurnWhite].KI > 0 {
						state.SetBoard(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
							Turn:  turn,
							Piece: shogi.KI,
						})
						state.Captured[shogi.TurnWhite].KI--
						switch pp.piece {
						case shogi.TO:
							state.Captured[shogi.TurnWhite].FU++
						case shogi.NY:
							state.Captured[shogi.TurnWhite].KY++
						case shogi.NK:
							state.Captured[shogi.TurnWhite].KE++
						case shogi.NG:
							state.Captured[shogi.TurnWhite].GI++
						}
					} else if state.Captured[shogi.TurnWhite].FU > 0 {
						if pp.piece == shogi.TO {
							continue
						}
						state.SetBoard(pp.pos.File, pp.pos.Rank, &shogi.BoardPiece{
							Turn:  shogi.TurnBlack,
							Piece: shogi.TO,
						})
						state.Captured[shogi.TurnWhite].FU--
						switch pp.piece {
						case shogi.NY:
							state.Captured[shogi.TurnWhite].KY++
						case shogi.NK:
							state.Captured[shogi.TurnWhite].KE++
						case shogi.NG:
							state.Captured[shogi.TurnWhite].GI++
						}
					}
				}
			}
		}
	}
	return state
}
