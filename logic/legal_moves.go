package logic

import (
	"github.com/sugyan/shogi"
)

type diff struct{ i, j int }

var reachableMap = map[shogi.Piece][]diff{
	shogi.BFU: {{-1, +0}},
	shogi.BKE: {{-2, +1}, {-2, -1}},
	shogi.BGI: {{-1, -1}, {-1, +0}, {-1, +1}, {+1, -1}, {+1, +1}},
	shogi.BKI: {{-1, -1}, {-1, +0}, {-1, +1}, {+0, -1}, {+0, +1}, {+1, +0}},
	shogi.BOU: {{-1, -1}, {-1, +0}, {-1, +1}, {+0, -1}, {+0, +1}, {+1, -1}, {+1, +0}, {+1, +1}},
	shogi.BUM: {{-1, +0}, {+0, -1}, {+0, +1}, {+1, +0}},
	shogi.BRY: {{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1}},
	shogi.WFU: {{+1, +0}},
	shogi.WKE: {{+2, +1}, {+2, -1}},
	shogi.WGI: {{+1, -1}, {+1, +0}, {+1, +1}, {-1, -1}, {-1, +1}},
	shogi.WKI: {{+1, -1}, {+1, +0}, {+1, +1}, {+0, -1}, {+0, +1}, {-1, +0}},
	shogi.WOU: {{+1, -1}, {+1, +0}, {+1, +1}, {+0, -1}, {+0, +1}, {-1, -1}, {-1, +0}, {-1, +1}},
	shogi.WUM: {{+1, +0}, {+0, -1}, {+0, +1}, {-1, +0}},
	shogi.WRY: {{+1, -1}, {+1, +1}, {-1, -1}, {-1, +1}},
}

var stepMap = map[shogi.Piece][]diff{
	shogi.BKY: {{-1, +0}},
	shogi.BKA: {{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1}},
	shogi.BHI: {{-1, +0}, {+0, -1}, {+0, +1}, {+1, +0}},
	shogi.WKY: {{+1, +0}},
	shogi.WKA: {{+1, -1}, {+1, +1}, {-1, -1}, {-1, +1}},
	shogi.WHI: {{+1, +0}, {+0, -1}, {+0, +1}, {-1, +0}},
}

func init() {
	for _, p := range []shogi.Piece{shogi.BTO, shogi.BNY, shogi.BNK, shogi.BNG} {
		reachableMap[p] = reachableMap[shogi.BKI]
	}
	for _, p := range []shogi.Piece{shogi.WTO, shogi.WNY, shogi.WNK, shogi.WNG} {
		reachableMap[p] = reachableMap[shogi.WKI]
	}
	stepMap[shogi.BUM] = stepMap[shogi.BKA]
	stepMap[shogi.BRY] = stepMap[shogi.BHI]
	stepMap[shogi.WUM] = stepMap[shogi.WKA]
	stepMap[shogi.WRY] = stepMap[shogi.WHI]
}

// LegalMoves method
func (s *State) LegalMoves() []*shogi.Move {
	moves := []*shogi.Move{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			p := s.board[i][j]
			if p != shogi.EMP && p.Turn() == s.turn {
				file, rank := 9-j, i+1
				for _, d := range reachableMap[p] {
					ii, jj := i+d.i, j+d.j
					if ii >= 0 && ii < 9 && jj >= 0 && jj < 9 {
						pp := s.board[ii][jj]
						if pp == shogi.EMP || pp.Turn() != s.turn {
							moves = append(moves, &shogi.Move{
								Src:   shogi.Position{File: file, Rank: rank},
								Dst:   shogi.Position{File: 9 - jj, Rank: ii + 1},
								Piece: p,
							})
						}
					}
				}
				for _, d := range stepMap[p] {
					for ii, jj := i+d.i, j+d.j; ii >= 0 && ii < 9 && jj >= 0 && jj < 9; {
						pp := s.board[ii][jj]
						if pp == shogi.EMP || pp.Turn() != s.turn {
							moves = append(moves, &shogi.Move{
								Src:   shogi.Position{File: file, Rank: rank},
								Dst:   shogi.Position{File: 9 - jj, Rank: ii + 1},
								Piece: p,
							})
						}
						if pp != shogi.EMP {
							break
						}
						ii += d.i
						jj += d.j
					}
				}
			}
		}
	}
	// promote
	for _, move := range moves {
		switch move.Piece.Raw() {
		case shogi.FU, shogi.KY:
			if (s.turn == shogi.TurnBlack && move.Dst.Rank < 2) || (s.turn == shogi.TurnWhite && move.Dst.Rank > 8) {
				move.Piece = move.Piece.Promote()
			}
		case shogi.KE:
			if (s.turn == shogi.TurnBlack && move.Dst.Rank < 3) || (s.turn == shogi.TurnWhite && move.Dst.Rank > 7) {
				move.Piece = move.Piece.Promote()
			}
		}
	}
	promoteMoves := []*shogi.Move{}
	for _, m := range moves {
		if !m.Piece.IsPromoted() {
			switch m.Piece.Turn() {
			case shogi.TurnBlack:
				if m.Src.Rank <= 3 || m.Dst.Rank <= 3 {
					promoteMoves = append(promoteMoves, &shogi.Move{
						Src:   m.Src,
						Dst:   m.Dst,
						Piece: m.Piece.Promote(),
					})
				}
			case shogi.TurnWhite:
				if m.Src.Rank >= 7 || m.Dst.Rank >= 7 {
					promoteMoves = append(promoteMoves, &shogi.Move{
						Src:   m.Src,
						Dst:   m.Dst,
						Piece: m.Piece.Promote(),
					})
				}
			}
		}
	}
	moves = append(moves, promoteMoves...)
	// use captured pieces
	capturedMoves := []*shogi.Move{}
	captured := s.captured[capturedIndex(s.turn)]
	if captured.Total() > 0 {
		positions := []*shogi.Position{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if s.board[i][j] == shogi.EMP {
					positions = append(positions, &shogi.Position{
						File: 9 - j,
						Rank: i + 1,
					})
				}
			}
		}
		if captured.FU > 0 {
			for _, position := range positions {
				if (s.turn == shogi.TurnBlack && position.Rank > 1) || (s.turn == shogi.TurnWhite && position.Rank < 9) {
					capturedMoves = append(capturedMoves, &shogi.Move{
						Src:   shogi.Position{File: 0, Rank: 0},
						Dst:   *position,
						Piece: shogi.MakePiece(shogi.FU, s.turn),
					})
				}
			}
		}
		if captured.KY > 0 {
			for _, position := range positions {
				if (s.turn == shogi.TurnBlack && position.Rank > 1) || (s.turn == shogi.TurnWhite && position.Rank < 9) {
					capturedMoves = append(capturedMoves, &shogi.Move{
						Src:   shogi.Position{File: 0, Rank: 0},
						Dst:   *position,
						Piece: shogi.MakePiece(shogi.KY, s.turn),
					})
				}
			}
		}
		if captured.KE > 0 {
			for _, position := range positions {
				if (s.turn == shogi.TurnBlack && position.Rank > 2) || (s.turn == shogi.TurnWhite && position.Rank < 8) {
					capturedMoves = append(capturedMoves, &shogi.Move{
						Src:   shogi.Position{File: 0, Rank: 0},
						Dst:   *position,
						Piece: shogi.MakePiece(shogi.KE, s.turn),
					})
				}
			}
		}
		if captured.GI > 0 {
			for _, position := range positions {
				capturedMoves = append(capturedMoves, &shogi.Move{
					Src:   shogi.Position{File: 0, Rank: 0},
					Dst:   *position,
					Piece: shogi.MakePiece(shogi.GI, s.turn),
				})
			}
		}
		if captured.KI > 0 {
			for _, position := range positions {
				capturedMoves = append(capturedMoves, &shogi.Move{
					Src:   shogi.Position{File: 0, Rank: 0},
					Dst:   *position,
					Piece: shogi.MakePiece(shogi.KI, s.turn),
				})
			}
		}
		if captured.KA > 0 {
			for _, position := range positions {
				capturedMoves = append(capturedMoves, &shogi.Move{
					Src:   shogi.Position{File: 0, Rank: 0},
					Dst:   *position,
					Piece: shogi.MakePiece(shogi.KA, s.turn),
				})
			}
		}
		if captured.HI > 0 {
			for _, position := range positions {
				capturedMoves = append(capturedMoves, &shogi.Move{
					Src:   shogi.Position{File: 0, Rank: 0},
					Dst:   *position,
					Piece: shogi.MakePiece(shogi.HI, s.turn),
				})
			}
		}
	}
	moves = append(moves, capturedMoves...)

	return moves
}
