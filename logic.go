package shogi

type diff struct{ i, j int }

var reachableMap = map[Piece][]diff{
	BFU: {{-1, +0}},
	BKE: {{-2, +1}, {-2, -1}},
	BGI: {{-1, -1}, {-1, +0}, {-1, +1}, {+1, -1}, {+1, +1}},
	BKI: {{-1, -1}, {-1, +0}, {-1, +1}, {+0, -1}, {+0, +1}, {+1, +0}},
	BOU: {{-1, -1}, {-1, +0}, {-1, +1}, {+0, -1}, {+0, +1}, {+1, -1}, {+1, +0}, {+1, +1}},
	BUM: {{-1, +0}, {+0, -1}, {+0, +1}, {+1, +0}},
	BRY: {{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1}},
	WFU: {{+1, +0}},
	WKE: {{+2, +1}, {+2, -1}},
	WGI: {{+1, -1}, {+1, +0}, {+1, +1}, {-1, -1}, {-1, +1}},
	WKI: {{+1, -1}, {+1, +0}, {+1, +1}, {+0, -1}, {+0, +1}, {-1, +0}},
	WOU: {{+1, -1}, {+1, +0}, {+1, +1}, {+0, -1}, {+0, +1}, {-1, -1}, {-1, +0}, {-1, +1}},
	WUM: {{+1, +0}, {+0, -1}, {+0, +1}, {-1, +0}},
	WRY: {{+1, -1}, {+1, +1}, {-1, -1}, {-1, +1}},
}

var stepMap = map[Piece][]diff{
	BKY: {{-1, +0}},
	BKA: {{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1}},
	BHI: {{-1, +0}, {+0, -1}, {+0, +1}, {+1, +0}},
	WKY: {{+1, +0}},
	WKA: {{+1, -1}, {+1, +1}, {-1, -1}, {-1, +1}},
	WHI: {{+1, +0}, {+0, -1}, {+0, +1}, {-1, +0}},
}

func init() {
	for _, p := range []Piece{BTO, BNY, BNK, BNG} {
		reachableMap[p] = reachableMap[BKI]
	}
	for _, p := range []Piece{WTO, WNY, WNK, WNG} {
		reachableMap[p] = reachableMap[WKI]
	}
	stepMap[BUM] = stepMap[BKA]
	stepMap[BRY] = stepMap[BHI]
	stepMap[WUM] = stepMap[WKA]
	stepMap[WRY] = stepMap[WHI]
}

// LegalMoves method
func (s *State) LegalMoves() []*Move {
	moves := []*Move{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			p := s.Board[i][j]
			if p != EMP && p.Turn() == s.Turn {
				file, rank := 9-j, i+1
				for _, d := range reachableMap[p] {
					ii, jj := i+d.i, j+d.j
					if ii >= 0 && ii < 9 && jj >= 0 && jj < 9 {
						pp := s.Board[ii][jj]
						if pp == EMP || pp.Turn() != s.Turn {
							moves = append(moves, &Move{Position{file, rank}, Position{9 - jj, ii + 1}, p})
						}
					}
				}
				for _, d := range stepMap[p] {
					for ii, jj := i+d.i, j+d.j; ii >= 0 && ii < 9 && jj >= 0 && jj < 9; {
						pp := s.Board[ii][jj]
						if pp == EMP || pp.Turn() != s.Turn {
							moves = append(moves, &Move{Position{file, rank}, Position{9 - jj, ii + 1}, p})
						}
						if pp != EMP {
							break
						}
						ii += d.i
						jj += d.j
					}
				}
			}
		}
	}
	return moves
}
