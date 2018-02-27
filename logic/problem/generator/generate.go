package generator

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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
	Type5 = &problemType{5}
)

type posPiece struct {
	pos   shogi.Position
	piece shogi.Piece
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type generator struct {
	steps   int
	timeout time.Duration
}

// Generate function
func Generate(pType Problem) (*shogi.State, int) {
	g := &generator{
		steps:   pType.Steps(),
		timeout: time.Second,
	}
	state := g.generate()
	score := g.calculateScore(state)
	return state, score
}

func (g *generator) generate() *shogi.State {
	for {
		var state *shogi.State
		// random generate
		for {
			state = random()
			if g.isCheckmate(state) {
				break
			}
		}
		// reduce pieces
		g.cut(state)
		// rewind and check
		for _, s := range g.rewind(state, shogi.TurnBlack) {
			switch g.steps {
			case 1:
				if g.isValidProblem(s, 1) {
					// TODO: evaluate
					g.cleanup(s)
					return s
				}
			case 3, 5:
				states := g.rewind(s, shogi.TurnWhite)
				if g.steps == 5 && len(states) > 0 {
					states = append(states, g.rewind(states[rand.Intn(len(states))], shogi.TurnWhite)...)
				}
				for _, i := range rand.Perm(len(states)) {
					if i > 5*(g.steps-1)/2 {
						break
					}
					result := states[i]
					if g.isValidProblem(result, g.steps) {
						g.cleanup(result)
						return result
					}
				}
			}
		}
	}
}

func (g *generator) isCheckmate(state *shogi.State) bool {
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
			root, err := solver.NewSolver(s).SolveWithTimeout(g.timeout)
			if err != nil {
				// timed out
				return false
			}
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
		if g.isCheckmate(s) {
			b := state.GetBoard(p.File, p.Rank)
			state.SetBoard(p.File, p.Rank, nil)
			state.Captured[shogi.TurnWhite].Add(b.Piece)
		}
	}
}

func hasMultipleAnswers(n node.Node, depth int) bool {
	if depth <= 2 {
		return false
	}
	switch n.Move().Turn {
	case shogi.TurnBlack:
		num := 0
		for _, c := range n.Children() {
			if c.Result() == node.ResultT {
				answer := solver.SearchBestAnswer(c)
				if len(answer) != depth-1 {
					continue
				}
				s := n.State().Clone()
				s.Apply(c.Move())
				for _, m := range answer {
					s.Apply(m)
				}
				if s.Captured[shogi.TurnBlack].Num() == 0 {
					num++
				}
			}
		}
		if num > 1 {
			return true
		}

		answer := solver.SearchBestAnswer(n)[0]
		for _, c := range n.Children() {
			m := c.Move()
			if m.Src == answer.Src && m.Dst == answer.Dst && m.Piece == answer.Piece {
				return hasMultipleAnswers(c, depth-1)
			}
		}
	case shogi.TurnWhite:
		num := 0
		for _, c := range n.Children() {
			if c.Result() == node.ResultT {
				answer := solver.SearchBestAnswer(c)
				if len(answer) != depth-1 {
					continue
				}
				num++
				if hasMultipleAnswers(c, depth-1) {
					return true
				}
			}
			if num > 1 {
				return true
			}
		}
	}
	return false
}

func (g *generator) isValidProblem(state *shogi.State, steps int) bool {
	// must consider wasted placing pieces
	root, err := solver.NewSolver(state).SolveWithTimeout(g.timeout)
	if err != nil {
		// timed out
		return false
	}
	bestAnswer := solver.SearchBestAnswer(root)

	// check answer length
	if len(bestAnswer) != steps {
		return false
	}
	// check captured pieces
	s := state.Clone()
	for _, m := range bestAnswer {
		s.Apply(m)
	}
	if s.Captured[shogi.TurnBlack].Num() > 0 {
		return false
	}
	// check if there are multiple answers
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
		if !hasMultipleAnswers(root, steps) {
			return true
		}
	}
	return false
}

func (g *generator) cleanup(state *shogi.State) *shogi.State {
	// remove unnecessary pieces
	for {
		removed := false
		posPieces := []*posPiece{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				file, rank := 9-i, j+1
				b := state.GetBoard(file, rank)
				if b != nil && b.Piece != shogi.OU {
					posPieces = append(posPieces, &posPiece{
						pos:   shogi.Pos(file, rank),
						piece: b.Piece,
					})
				}
			}
		}
		for _, i := range rand.Perm(len(posPieces)) {
			pp := posPieces[i]
			s := state.Clone()
			s.SetBoard(pp.pos.File, pp.pos.Rank, nil)
			s.Captured[shogi.TurnWhite].Add(pp.piece)
			if s.Check(shogi.TurnBlack) == nil && g.isValidProblem(s, g.steps) {
				state.SetBoard(pp.pos.File, pp.pos.Rank, nil)
				state.Captured[shogi.TurnWhite].Add(pp.piece)
				removed = true
			}
		}
		if !removed {
			break
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

func countScore(n node.Node, depth int) int {
	if depth == 0 {
		return 0
	}

	sum := 1
	scoreMap := map[node.Result][]int{}
	dup := map[string]struct{}{}
	for _, c := range n.Children() {
		if c.Result() == node.ResultU {
			continue
		}
		dupKeys := []string{}
		for _, cc := range c.Children() {
			dupKeys = append(dupKeys, fmt.Sprintf("%v", cc.Move()))
		}
		dupKey := strings.Join(dupKeys, ":")
		if _, exist := dup[dupKey]; exist {
			continue
		} else {
			dup[dupKey] = struct{}{}
		}
		score := countScore(c, depth-1)
		scoreMap[c.Result()] = append(scoreMap[c.Result()], score)
	}
	type minsum struct {
		min, sum int
	}
	ms := map[node.Result]*minsum{
		node.ResultT: &minsum{math.MaxInt32, 0},
		node.ResultF: &minsum{math.MaxInt32, 0},
	}
	for _, result := range []node.Result{node.ResultT, node.ResultF} {
		for _, c := range scoreMap[result] {
			if c < ms[result].min {
				ms[result].min = c
			}
			ms[result].sum += c
		}
		if len(scoreMap[result]) == 0 {
			ms[result].min = 0
		}
	}
	switch n.Move().Turn {
	case shogi.TurnBlack:
		switch n.Result() {
		case node.ResultT:
			sum += ms[node.ResultT].sum
		case node.ResultF:
			sum += ms[node.ResultF].min + ms[node.ResultT].sum
		}
	case shogi.TurnWhite:
		switch n.Result() {
		case node.ResultT:
			sum += ms[node.ResultT].min + ms[node.ResultF].sum
		case node.ResultF:
			sum += ms[node.ResultF].min
		}
	}
	return sum
}

func (g *generator) calculateScore(state *shogi.State) int {
	root, _ := solver.NewSolver(state).SolveWithTimeout(0)
	return countScore(root, g.steps+1)
}
