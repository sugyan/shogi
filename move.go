package shogi

import (
	"fmt"
)

// Move type
type Move struct {
	Turn  Turn
	Src   Position
	Dst   Position
	Piece Piece
}

// Apply move method
func (s *State) Apply(move *Move) {
	// update state
	s.latestMove = move
	if move.Src.IsCaptured() {
		s.SetBoard(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
		s.Captured[move.Turn].Sub(move.Piece)
	} else {
		b := s.GetBoard(move.Dst.File, move.Dst.Rank)
		if b != nil {
			s.Captured[move.Turn].Add(b.Piece)
		}
		s.SetBoard(move.Src.File, move.Src.Rank, nil)
		s.SetBoard(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
	}
}

// MoveStrings method
func (s *State) MoveStrings(moves []*Move) ([]string, error) {
	results := []string{}
	state := s.Clone()
	for _, move := range moves {
		ms, err := state.MoveString(move)
		if err != nil {
			return nil, err
		}
		state.Apply(move)
		results = append(results, ms)
	}
	return results, nil
}

// MoveString method
func (s *State) MoveString(move *Move) (string, error) {
	// move string
	nameMap := map[Piece]string{
		FU: "歩",
		TO: "と",
		KY: "香",
		NY: "成香",
		KE: "桂",
		NK: "成桂",
		GI: "銀",
		NG: "成銀",
		KI: "金",
		KA: "角",
		UM: "馬",
		HI: "飛",
		RY: "竜",
		OU: "玉",
	}
	result := "▲"
	if move.Turn == TurnWhite {
		result = "△"
	}
	if s.latestMove != nil && move.Dst == s.latestMove.Dst {
		result += "同"
	} else {
		result += fmt.Sprintf("%c%c",
			[]rune("123456789")[move.Dst.File-1],
			[]rune("一二三四五六七八九")[move.Dst.Rank-1],
		)
	}
	if move.Src.IsCaptured() {
		result += nameMap[move.Piece]
		for _, m := range s.CandidateMoves(move.Turn) {
			if m.Dst == move.Dst && m.Piece == move.Piece {
				result += "打"
			}
		}
	} else {
		// check movable
		srcPiece := s.GetBoard(move.Src.File, move.Src.Rank).Piece
		moves := []*Move{}
		ok := false
		for _, m := range s.CandidateMoves(move.Turn) {
			if m.Dst == move.Dst && m.Piece == move.Piece {
				ok = true
				if s.GetBoard(m.Src.File, m.Src.Rank).Piece == srcPiece {
					moves = append(moves, m)
				}
			}
		}
		if !ok {
			return "", fmt.Errorf("piece %s does not exist which move to (%d, %d)", move.Piece, move.Dst.File, move.Dst.Rank)
		}
		// piece name
		b := s.GetBoard(move.Src.File, move.Src.Rank)
		if b == nil {
			return "", fmt.Errorf("piece does not exist in (%d, %d)", move.Src.File, move.Src.Rank)
		}
		result += nameMap[b.Piece]
		// relative position and movements
		if len(moves) > 1 {
			sameFile := false
			sameRank := false
			switch move.Piece {
			case KA, UM, HI, RY:
				if moves[0].Src.File == moves[1].Src.File {
					sameFile = true
				}
				if (moves[0].Src.Rank == moves[1].Src.Rank) ||
					(moves[0].Src.Rank > move.Dst.Rank && moves[1].Src.Rank > move.Dst.Rank) ||
					(moves[0].Src.Rank < move.Dst.Rank && moves[1].Src.Rank < move.Dst.Rank) {
					sameRank = true
				}
			default:
				for _, m := range moves {
					if m.Src != move.Src && m.Src.File == move.Src.File {
						sameFile = true
					}
					if m.Src != move.Src && m.Src.Rank == move.Src.Rank {
						sameRank = true
					}
				}
			}
			if sameRank {
				d := move.Src.File - move.Dst.File
				if move.Turn == TurnWhite {
					d *= -1
				}
				switch {
				case d == 0:
					if move.Piece == RY || move.Piece == UM {
						right := false
						if move.Src != moves[0].Src {
							if move.Src.File < moves[0].Src.File {
								right = true
							}
						} else {
							if move.Src.File < moves[1].Src.File {
								right = true
							}
						}
						if move.Turn == TurnWhite {
							right = !right
						}
						if right {
							result += "右"
						} else {
							result += "左"
						}
					} else {
						result += "直"
					}
				case d > 0:
					result += "左"
				case d < 0:
					result += "右"
				}
			}
			if !sameRank || (sameFile && move.Src.File != move.Dst.File) {
				d := move.Src.Rank - move.Dst.Rank
				if move.Turn == TurnWhite {
					d *= -1
				}
				switch {
				case d == 0:
					result += "寄"
				case d > 0:
					result += "上"
				case d < 0:
					result += "引"
				}
			}
		}
		if b.Piece != move.Piece {
			result += "成"
		} else if !move.Piece.Promoted() {
			switch move.Piece {
			case KI, OU:
			// noop
			default:
				if (move.Turn == TurnBlack && (move.Src.Rank <= 3 || move.Dst.Rank <= 3)) ||
					(move.Turn == TurnWhite && (move.Src.Rank >= 7 || move.Dst.Rank >= 7)) {
					result += "不成"
				}
			}
		}
	}
	return result, nil
}
