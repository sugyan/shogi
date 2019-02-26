package shogi_test

import (
	"testing"

	"github.com/sugyan/shogi"
)

var initialState = &shogi.State{
	Board: [9][9]shogi.Piece{
		{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
		{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKA, shogi.EMP},
		{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU},
		{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
		{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
		{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
		{shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
		{shogi.EMP, shogi.BKA, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
		{shogi.BKY, shogi.BKE, shogi.BGI, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
	},
}

func TestMove(t *testing.T) {
	type result struct {
		err   error
		state shogi.State
	}
	testCases := []struct {
		moves    []*shogi.Move
		expected result
	}{
		{
			[]*shogi.Move{
				{Src: shogi.Position{File: 5, Rank: 7}, Dst: shogi.Position{File: 5, Rank: 6}, Piece: shogi.BFU},
				{Src: shogi.Position{File: 3, Rank: 3}, Dst: shogi.Position{File: 3, Rank: 4}, Piece: shogi.WFU},
				{Src: shogi.Position{File: 7, Rank: 7}, Dst: shogi.Position{File: 7, Rank: 6}, Piece: shogi.BFU},
				{Src: shogi.Position{File: 2, Rank: 2}, Dst: shogi.Position{File: 8, Rank: 8}, Piece: shogi.WUM},
				{Src: shogi.Position{File: 7, Rank: 9}, Dst: shogi.Position{File: 8, Rank: 8}, Piece: shogi.BGI},
				{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 5, Rank: 7}, Piece: shogi.WKA},
			},
			result{
				err: nil,
				state: shogi.State{
					Board: [9][9]shogi.Piece{
						{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
						{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.EMP, shogi.WFU, shogi.WFU},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WFU, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.BFU, shogi.EMP, shogi.BFU, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.BFU, shogi.BFU, shogi.EMP, shogi.BFU, shogi.WKA, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
						{shogi.EMP, shogi.BGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
						{shogi.BKY, shogi.BKE, shogi.EMP, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
					},
					Captured: [2]shogi.Captured{
						{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 1, HI: 0},
						{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
					},
				},
			},
		},
		{
			[]*shogi.Move{
				{Src: shogi.Position{File: 5, Rank: 9}, Dst: shogi.Position{File: 5, Rank: 8}, Piece: shogi.BFU},
			},
			result{
				err: shogi.ErrInvalidMove,
			},
		},
	}
	for i, testCase := range testCases {
		s := *initialState
		err := s.Move(testCase.moves...)
		if err != testCase.expected.err {
			t.Errorf("#%d, err got: %v, expected: %v", i, err, testCase.expected.err)
			continue
		}
		if err != nil {
			continue
		}
		if s != testCase.expected.state {
			t.Errorf("#%d: state got: %v, expected: %v", i, s, testCase.expected.state)
		}
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		state    *shogi.State
		expected string
	}{
		{
			state: &shogi.State{},
			expected: `
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  *  *  *  *  *  * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * `[1:],
		},
		{
			state: initialState,
			expected: `
P1-KY-KE-GI-KI-OU-KI-GI-KE-KY
P2 * -HI *  *  *  *  * -KA * 
P3-FU-FU-FU-FU-FU-FU-FU-FU-FU
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7+FU+FU+FU+FU+FU+FU+FU+FU+FU
P8 * +KA *  *  *  *  * +HI * 
P9+KY+KE+GI+KI+OU+KI+GI+KE+KY`[1:],
		},
	}
	for i, testCase := range testCases {
		s := testCase.state.String()
		if s != testCase.expected {
			t.Errorf("#%d: got %v, expected: %v", i, s, testCase.expected)
		}
	}
}
