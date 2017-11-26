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
	if move.Src.File > 0 && move.Src.Rank > 0 {
		bp := s.GetBoardPiece(move.Dst.File, move.Dst.Rank)
		if bp != nil {
			s.Captured[move.Turn].AddPieces(bp.Piece)
		}
		s.SetBoardPiece(move.Src.File, move.Src.Rank, nil)
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
	} else {
		s.SetBoardPiece(move.Dst.File, move.Dst.Rank, &BoardPiece{
			Turn:  move.Turn,
			Piece: move.Piece,
		})
		s.Captured[move.Turn].SubPieces(move.Piece)
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
	if move.Turn == TurnSecond {
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
		ok := false
		for _, m := range s.CandidateMoves(move.Turn) {
			if m.Dst == move.Dst && m.Piece == move.Piece {
				ok = true
				break
			}
		}
		if !ok {
			return "", fmt.Errorf("piece %s does not exist which move to (%d, %d)", move.Piece, move.Dst.File, move.Dst.Rank)
		}
		bp := s.GetBoardPiece(move.Src.File, move.Src.Rank)
		if bp == nil {
			return "", fmt.Errorf("piece does not exist in (%d, %d)", move.Src.File, move.Src.Rank)
		}
		if bp.Piece != move.Piece {
			result += nameMap[bp.Piece] + "成"
		} else {
			result += nameMap[move.Piece]
			if !move.Piece.Promoted() && ((move.Turn == TurnFirst && move.Dst.Rank <= 3) || (move.Turn == TurnSecond && move.Dst.Rank >= 7)) {
				result += "不成"
			}
		}
	}
	// TODO special case
	return result, nil
}
