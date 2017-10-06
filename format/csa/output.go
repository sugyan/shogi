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
			p := state.Board[i][j]
			if p != nil {
				switch p.Turn() {
				case shogi.TurnFirst:
					result = append(result, '+')
				case shogi.TurnSecond:
					result = append(result, '-')
				}
				result = append(result, []byte(p.Code())...)
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
			p := state.Board[i][j]
			if p != nil {
				pos := &position{
					Rank:  i + 1,
					File:  9 - j,
					Piece: p,
				}
				pieces[p.Turn()] = append(pieces[p.Turn()], pos)
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
			s := fmt.Sprintf("%d%d%s", pos.File, pos.Rank, pos.Piece.Code())
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
			Fu: 18 - lessCap.Fu,
			Ky: 4 - lessCap.Ky,
			Ke: 4 - lessCap.Ke,
			Gi: 4 - lessCap.Gi,
			Ki: 4 - lessCap.Ki,
			Ka: 2 - lessCap.Ka,
			Hi: 2 - lessCap.Hi,
		}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				p := state.Board[i][j]
				if p != nil {
					switch shogi.PieceCode(p.Code()) {
					case shogi.FU:
						fallthrough
					case shogi.TO:
						remains.Fu--
					case shogi.KY:
						fallthrough
					case shogi.NY:
						remains.Ky--
					case shogi.KE:
						fallthrough
					case shogi.NK:
						remains.Ke--
					case shogi.GI:
						fallthrough
					case shogi.NG:
						remains.Gi--
					case shogi.KI:
						remains.Ki--
					case shogi.KA:
						fallthrough
					case shogi.UM:
						remains.Ka--
					case shogi.HI:
						fallthrough
					case shogi.RY:
						remains.Hi--
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
		for i := 0; i < pieces.Hi; i++ {
			result = append(result, []byte(`00HI`)...)
		}
		for i := 0; i < pieces.Ka; i++ {
			result = append(result, []byte(`00KA`)...)
		}
		for i := 0; i < pieces.Ki; i++ {
			result = append(result, []byte(`00KI`)...)
		}
		for i := 0; i < pieces.Gi; i++ {
			result = append(result, []byte(`00GI`)...)
		}
		for i := 0; i < pieces.Ke; i++ {
			result = append(result, []byte(`00KE`)...)
		}
		for i := 0; i < pieces.Ky; i++ {
			result = append(result, []byte(`00KY`)...)
		}
		for i := 0; i < pieces.Fu; i++ {
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
