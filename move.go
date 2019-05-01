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
	pieceMap := map[bool]map[rawPiece]string{
		false: {
			fu: "歩",
			ky: "香",
			ke: "桂",
			gi: "銀",
			ki: "金",
			ka: "角",
			hi: "飛",
			ou: "玉",
		},
		true: {
			fu: "と",
			ky: "成香",
			ke: "成桂",
			gi: "成銀",
			ka: "馬",
			hi: "竜",
		},
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
	// 打
	if move.Src.File == 0 && move.Src.Rank == 0 {
		b.WriteString(pieceMap[false][move.Piece.raw()])
		for _, m := range state.LegalMoves() {
			if m.Src != move.Src && m.Dst == move.Dst && m.Piece == move.Piece {
				b.WriteRune('打')
			}
		}
		return b.String(), nil
	}

	orig := state.Board[move.Src.Rank-1][9-move.Src.File]
	b.WriteString(pieceMap[orig.IsPromoted()][move.Piece.raw()])
	if orig.raw() != move.Piece.raw() {
		return "", ErrInvalidMove
	}
	// TODO
	// 成・不成
	if orig != move.Piece {
		b.WriteRune('成')
	} else if !move.Piece.IsPromoted() {
		switch move.Piece.raw() {
		case fu, ky, ke, gi, ka, hi:
			if (move.Piece.Turn() == TurnBlack && (move.Src.Rank <= 3 || move.Dst.Rank <= 3)) ||
				(move.Piece.Turn() == TurnWhite && (move.Src.Rank >= 7 || move.Dst.Rank >= 7)) {
				b.WriteString("不成")
			}
		}
	}
	return b.String(), nil
}
