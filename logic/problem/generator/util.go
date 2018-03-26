package generator

import (
	"math/rand"

	"github.com/sugyan/shogi"
)

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
			if rand.Intn(4) > 0 {
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
				if turn == shogi.TurnWhite && p.Rank >= 9 {
					piece = shogi.TO
				}
			case shogi.KY:
				if turn == shogi.TurnBlack {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NY
					} else if p.Rank <= 1 {
						turn = shogi.TurnWhite
					}
				}
				if turn == shogi.TurnWhite && p.Rank >= 9 {
					turn = shogi.TurnBlack
				}
			case shogi.KE:
				if turn == shogi.TurnBlack {
					if p.Rank <= 3 && rand.Intn(4) == 0 {
						piece = shogi.NK
					} else if p.Rank <= 2 {
						turn = shogi.TurnWhite
					}
				}
				if turn == shogi.TurnWhite && p.Rank >= 8 {
					turn = shogi.TurnBlack
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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
