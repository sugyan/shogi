package shogi

import (
	"strings"
)

// Position struct
type Position struct {
	File, Rank int
}

// Move struct
type Move struct {
	Src   Position
	Dst   Position
	Piece Piece
}

// MoveStrings function
func MoveStrings(state *State, moves ...*Move) ([]string, error) {
	result := make([]string, 0, len(moves))
	s := *state
	var prev *Move
	for _, m := range moves {
		str, err := moveString(&s, m, prev)
		if err != nil {
			return nil, err
		}
		result = append(result, str)
		prev = m
		s.Move(m)
	}
	return result, nil
}

func moveString(state *State, move, prev *Move) (string, error) {
	pieceMap := map[rawPiece]string{
		fu: "歩",
		ky: "香",
		ke: "桂",
		gi: "銀",
		ki: "金",
		ka: "角",
		hi: "飛",
		ou: "玉",
		to: "と",
		ny: "成香",
		nk: "成桂",
		ng: "成銀",
		um: "馬",
		ry: "竜",
	}
	files := []rune("123456789")
	ranks := []rune("一二三四五六七八九")
	b := &strings.Builder{}
	switch move.Piece.Turn() {
	case TurnBlack:
		b.WriteRune('▲')
	case TurnWhite:
		b.WriteRune('△')
	}
	if prev != nil && move.Dst == prev.Dst {
		// 同
	} else {
		b.WriteRune(files[move.Dst.File-1])
		b.WriteRune(ranks[move.Dst.Rank-1])
	}
	b.WriteString(pieceMap[move.Piece.raw()])
	if move.Src.File == 0 && move.Src.Rank == 0 {
		for _, m := range state.LegalMoves() {
			if m.Src != move.Src && m.Dst == move.Dst && m.Piece == move.Piece {
				b.WriteRune('打')
				break
			}
		}
	} else {
		// TODO
	}
	return b.String(), nil
}
