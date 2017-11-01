package csa

import (
	"fmt"

	"github.com/sugyan/shogi"
)

// InitialState1 function
func InitialState1(state *shogi.State) string {
	result := make([]byte, 0, 9*(2+9*3+1))
	for i := 0; i < 9; i++ {
		rowName := fmt.Sprintf("P%d", i+1)
		result = append(result, []byte(rowName)...)
		for j := 0; j < 9; j++ {
			bp := state.Board[i][j]
			if bp != nil {
				switch bp.Turn {
				case shogi.TurnFirst:
					result = append(result, '+')
				case shogi.TurnSecond:
					result = append(result, '-')
				}
				result = append(result, []byte(bp.Piece.String())...)
			} else {
				result = append(result, []byte(` * `)...)
			}
		}
		result = append(result, '\n')
	}
	return string(result) + handPieces(state)
}

// InitialState2 function
func InitialState2(state *shogi.State) string {
	result := make([]byte, 0, 10)
	type position struct {
		Rank  int
		File  int
		Piece shogi.Piece
	}
	pieces := make(map[shogi.Turn][]*position)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := state.Board[i][j]
			if bp != nil {
				pos := &position{
					Rank:  i + 1,
					File:  9 - j,
					Piece: bp.Piece,
				}
				pieces[bp.Turn] = append(pieces[bp.Turn], pos)
			}
		}
	}
	for move, positions := range pieces {
		result = append(result, 'P')
		switch move {
		case shogi.TurnFirst:
			result = append(result, '+')
		case shogi.TurnSecond:
			result = append(result, '-')
		}
		for _, pos := range positions {
			s := fmt.Sprintf("%d%d%s", pos.File, pos.Rank, pos.Piece.String())
			result = append(result, []byte(s)...)
		}
		result = append(result, '\n')
	}
	return string(result) + handPieces(state)
}

func handPieces(state *shogi.State) string {
	result := make([]byte, 0)
	useAll := false
	var useAllTarget shogi.Turn
	{
		less := shogi.TurnFirst
		if state.Captured[shogi.TurnFirst].Num() > state.Captured[shogi.TurnSecond].Num() {
			less = shogi.TurnSecond
		}
		lessCap := state.Captured[less]
		remains := &shogi.CapturedPieces{
			FU: 18 - lessCap.FU,
			KY: 4 - lessCap.KY,
			KE: 4 - lessCap.KE,
			GI: 4 - lessCap.GI,
			KI: 4 - lessCap.KI,
			KA: 2 - lessCap.KA,
			HI: 2 - lessCap.HI,
		}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				bp := state.Board[i][j]
				if bp != nil {
					switch bp.Piece {
					case shogi.FU:
						fallthrough
					case shogi.TO:
						remains.FU--
					case shogi.KY:
						fallthrough
					case shogi.NY:
						remains.KY--
					case shogi.KE:
						fallthrough
					case shogi.NK:
						remains.KE--
					case shogi.GI:
						fallthrough
					case shogi.NG:
						remains.GI--
					case shogi.KI:
						remains.KI--
					case shogi.KA:
						fallthrough
					case shogi.UM:
						remains.KA--
					case shogi.HI:
						fallthrough
					case shogi.RY:
						remains.HI--
					}
				}
			}
		}
		if remains.Num() == state.Captured[!less].Num() {
			useAll = true
			useAllTarget = !less
		}
	}
	for move, pieces := range state.Captured {
		if useAll && move == useAllTarget {
			continue
		}
		if pieces.Num() == 0 {
			continue
		}
		result = append(result, 'P')
		switch move {
		case shogi.TurnFirst:
			result = append(result, '+')
		case shogi.TurnSecond:
			result = append(result, '-')
		}
		for i := 0; i < pieces.HI; i++ {
			result = append(result, []byte(`00HI`)...)
		}
		for i := 0; i < pieces.KA; i++ {
			result = append(result, []byte(`00KA`)...)
		}
		for i := 0; i < pieces.KI; i++ {
			result = append(result, []byte(`00KI`)...)
		}
		for i := 0; i < pieces.GI; i++ {
			result = append(result, []byte(`00GI`)...)
		}
		for i := 0; i < pieces.KE; i++ {
			result = append(result, []byte(`00KE`)...)
		}
		for i := 0; i < pieces.KY; i++ {
			result = append(result, []byte(`00KY`)...)
		}
		for i := 0; i < pieces.FU; i++ {
			result = append(result, []byte(`00FU`)...)
		}
		result = append(result, '\n')
	}
	if useAll {
		result = append(result, 'P')
		switch useAllTarget {
		case shogi.TurnFirst:
			result = append(result, '+')
		case shogi.TurnSecond:
			result = append(result, '-')
		}
		result = append(result, []byte(`00AL`)...)
		result = append(result, '\n')
	}
	return string(result)
}
