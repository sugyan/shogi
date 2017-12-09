package csa

import (
	"fmt"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/record"
)

type initialStateOption int

// constant variables
const (
	InitialStateOption1 initialStateOption = iota
	InitialStateOption2
)

// Converter type
type Converter struct {
	opt *ConvertOption
}

// ConvertOption type
type ConvertOption struct {
	InitialState initialStateOption
}

// NewConverter function
func NewConverter(opt *ConvertOption) *Converter {
	if opt == nil {
		opt = &ConvertOption{}
	}
	return &Converter{opt: opt}
}

// ConvertToString function
func (c *Converter) ConvertToString(record *record.Record) string {
	result := ""
	switch c.opt.InitialState {
	case InitialStateOption1:
		result += InitialState1(record.State)
	case InitialStateOption2:
		result += InitialState2(record.State)
	}
	// TODO
	result += "+\n"
	for _, move := range record.Moves {
		result += moveToString(move) + "\n"
	}
	return result
}

func moveToString(move *shogi.Move) string {
	result := make([]byte, 0, 6)
	switch move.Turn {
	case shogi.TurnBlack:
		result = append(result, '+')
	case shogi.TurnWhite:
		result = append(result, '-')
	}
	result = append(result, '0'+byte(move.Src.File))
	result = append(result, '0'+byte(move.Src.Rank))
	result = append(result, '0'+byte(move.Dst.File))
	result = append(result, '0'+byte(move.Dst.Rank))
	result = append(result, []byte(move.Piece.String())...)
	return string(result)
}

// InitialState1 function
func InitialState1(state *shogi.State) string {
	result := make([]byte, 0, 9*(2+9*3+1))
	for i := 0; i < 9; i++ {
		rowName := fmt.Sprintf("P%d", i+1)
		result = append(result, []byte(rowName)...)
		for j := 0; j < 9; j++ {
			b := state.Board[i][j]
			if b != nil {
				switch b.Turn {
				case shogi.TurnBlack:
					result = append(result, '+')
				case shogi.TurnWhite:
					result = append(result, '-')
				}
				result = append(result, []byte(b.Piece.String())...)
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
			b := state.Board[i][j]
			if b != nil {
				pos := &position{
					Rank:  i + 1,
					File:  9 - j,
					Piece: b.Piece,
				}
				pieces[b.Turn] = append(pieces[b.Turn], pos)
			}
		}
	}
	for move, positions := range pieces {
		result = append(result, 'P')
		switch move {
		case shogi.TurnBlack:
			result = append(result, '+')
		case shogi.TurnWhite:
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
		less := shogi.TurnBlack
		if state.Captured[shogi.TurnBlack].Num() > state.Captured[shogi.TurnWhite].Num() {
			less = shogi.TurnWhite
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
				b := state.Board[i][j]
				if b != nil {
					switch b.Piece {
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
		case shogi.TurnBlack:
			result = append(result, '+')
		case shogi.TurnWhite:
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
		case shogi.TurnBlack:
			result = append(result, '+')
		case shogi.TurnWhite:
			result = append(result, '-')
		}
		result = append(result, []byte(`00AL`)...)
		result = append(result, '\n')
	}
	return string(result)
}
