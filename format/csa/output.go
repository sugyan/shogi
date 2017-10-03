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
				if p.IsFirst() {
					result = append(result, '+')
				} else {
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
	pieces := make(map[shogi.Move][]*position)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			p := state.Board[i][j]
			if p != nil {
				pos := &position{
					Rank:  i + 1,
					File:  9 - j,
					Piece: p,
				}
				if p.IsFirst() {
					pieces[shogi.MoveFirst] = append(pieces[shogi.MoveFirst], pos)
				} else {
					pieces[shogi.MoveSecond] = append(pieces[shogi.MoveSecond], pos)
				}
			}
		}
	}
	for move, positions := range pieces {
		result = append(result, 'P')
		switch move {
		case shogi.MoveFirst:
			result = append(result, '+')
		case shogi.MoveSecond:
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
	var useAll shogi.Move
	if len(state.Hands[shogi.MoveFirst]) > len(state.Hands[shogi.MoveSecond]) {
		useAll = shogi.MoveFirst
	} else {
		useAll = shogi.MoveSecond
	}
	for move, pieces := range state.Hands {
		if move == useAll {
			continue
		}
		if len(pieces) == 0 {
			continue
		}
		result = append(result, 'P')
		switch move {
		case shogi.MoveFirst:
			result = append(result, '+')
		case shogi.MoveSecond:
			result = append(result, '-')
		}
		for _, p := range pieces {
			result = append(result, []byte(`00`+p.Code())...)
		}
		result = append(result, '\n')
	}
	if len(state.Hands[useAll]) > 0 {
		result = append(result, 'P')
		switch useAll {
		case shogi.MoveFirst:
			result = append(result, '+')
		case shogi.MoveSecond:
			result = append(result, '-')
		}
		result = append(result, []byte(`00AL`)...)
		result = append(result, '\n')
	}
	return string(result)
}
