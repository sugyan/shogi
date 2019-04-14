package shogi_test

import (
	"strings"
	"testing"

	"github.com/sugyan/shogi"
)

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
		s := *shogi.InitialState
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

func TestMoveStrings(t *testing.T) {
	move := func(srcFile, srcRank, dstFile, dstRank int, piece shogi.Piece) *shogi.Move {
		return &shogi.Move{
			Src:   shogi.Position{File: srcFile, Rank: srcRank},
			Dst:   shogi.Position{File: dstFile, Rank: dstRank},
			Piece: piece,
		}
	}
	type testCase struct {
		moves    []*shogi.Move
		expected []string
	}
	test := func(s *shogi.State, testCases []*testCase) {
		for i, testCase := range testCases {
			results, err := shogi.MoveStrings(s, testCase.moves...)
			if err != nil {
				t.Error(err)
				continue
			}
			if len(results) != len(testCase.expected) {
				t.Errorf("length got: %d, expected: %d", len(results), len(testCase.expected))
				continue
			}
			result := strings.Join(results, " ")
			expected := strings.Join(testCase.expected, " ")
			if result != expected {
				t.Errorf("#%d: move string got: %s, expected: %s", i, result, expected)
			}
			t.Logf(result)
		}
	}
	// 打
	{
		state := &shogi.State{
			Board: [9][9]shogi.Piece{
				{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKA, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.BKA, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
			},
			Captured: [2]shogi.Captured{
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
			},
		}
		sb := *state
		sb.Turn = shogi.TurnBlack
		sw := *state
		sw.Turn = shogi.TurnWhite
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 9, 1, 8, shogi.BKY)}, []string{"▲1八香"}},
				{[]*shogi.Move{move(0, 0, 1, 8, shogi.BKY)}, []string{"▲1八香打"}},
				{[]*shogi.Move{move(0, 0, 2, 8, shogi.BKY)}, []string{"▲2八香"}},
				{[]*shogi.Move{move(2, 9, 1, 7, shogi.BKE)}, []string{"▲1七桂"}},
				{[]*shogi.Move{move(0, 0, 1, 7, shogi.BKE)}, []string{"▲1七桂打"}},
				{[]*shogi.Move{move(0, 0, 2, 7, shogi.BKE)}, []string{"▲2七桂"}},
				{[]*shogi.Move{move(3, 9, 3, 8, shogi.BGI)}, []string{"▲3八銀"}},
				{[]*shogi.Move{move(0, 0, 3, 8, shogi.BGI)}, []string{"▲3八銀打"}},
				{[]*shogi.Move{move(0, 0, 3, 7, shogi.BGI)}, []string{"▲3七銀"}},
				{[]*shogi.Move{move(4, 9, 5, 8, shogi.BKI)}, []string{"▲5八金"}},
				{[]*shogi.Move{move(0, 0, 5, 8, shogi.BKI)}, []string{"▲5八金打"}},
				{[]*shogi.Move{move(0, 0, 5, 7, shogi.BKI)}, []string{"▲5七金"}},
				{[]*shogi.Move{move(8, 8, 7, 7, shogi.BKA)}, []string{"▲7七角"}},
				{[]*shogi.Move{move(0, 0, 7, 7, shogi.BKA)}, []string{"▲7七角打"}},
				{[]*shogi.Move{move(0, 0, 8, 7, shogi.BKA)}, []string{"▲8七角"}},
				{[]*shogi.Move{move(2, 8, 2, 7, shogi.BHI)}, []string{"▲2七飛"}},
				{[]*shogi.Move{move(0, 0, 2, 7, shogi.BHI)}, []string{"▲2七飛打"}},
				{[]*shogi.Move{move(0, 0, 3, 7, shogi.BHI)}, []string{"▲3七飛"}},
			}
			test(&sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 1, 9, 2, shogi.WKY)}, []string{"△9二香"}},
				{[]*shogi.Move{move(0, 0, 9, 2, shogi.WKY)}, []string{"△9二香打"}},
				{[]*shogi.Move{move(0, 0, 8, 2, shogi.WKY)}, []string{"△8二香"}},
				{[]*shogi.Move{move(8, 1, 9, 3, shogi.WKE)}, []string{"△9三桂"}},
				{[]*shogi.Move{move(0, 0, 9, 3, shogi.WKE)}, []string{"△9三桂打"}},
				{[]*shogi.Move{move(0, 0, 8, 3, shogi.WKE)}, []string{"△8三桂"}},
				{[]*shogi.Move{move(7, 1, 7, 2, shogi.WGI)}, []string{"△7二銀"}},
				{[]*shogi.Move{move(0, 0, 7, 2, shogi.WGI)}, []string{"△7二銀打"}},
				{[]*shogi.Move{move(0, 0, 7, 3, shogi.WGI)}, []string{"△7三銀"}},
				{[]*shogi.Move{move(6, 1, 5, 2, shogi.WKI)}, []string{"△5二金"}},
				{[]*shogi.Move{move(0, 0, 5, 2, shogi.WKI)}, []string{"△5二金打"}},
				{[]*shogi.Move{move(0, 0, 5, 3, shogi.WKI)}, []string{"△5三金"}},
				{[]*shogi.Move{move(2, 2, 3, 3, shogi.WKA)}, []string{"△3三角"}},
				{[]*shogi.Move{move(0, 0, 3, 3, shogi.WKA)}, []string{"△3三角打"}},
				{[]*shogi.Move{move(0, 0, 2, 3, shogi.WKA)}, []string{"△2三角"}},
				{[]*shogi.Move{move(8, 2, 8, 3, shogi.WHI)}, []string{"△8三飛"}},
				{[]*shogi.Move{move(0, 0, 8, 3, shogi.WHI)}, []string{"△8三飛打"}},
				{[]*shogi.Move{move(0, 0, 7, 3, shogi.WHI)}, []string{"△7三飛"}},
			}
			test(&sw, testCases)
		}
	}
}
