package shogi_test

import (
	"strings"
	"testing"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic"
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
				state: logic.NewState(
					[9][9]shogi.Piece{
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
					[2]shogi.Captured{
						{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 1, HI: 0},
						{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
					},
					shogi.TurnBlack,
				),
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
		s := logic.NewInitialState()
		err := s.Move(testCase.moves...)
		if err != testCase.expected.err {
			t.Errorf("#%d, err got: %v, expected: %v", i, err, testCase.expected.err)
			continue
		}
		if err != nil {
			continue
		}
		if !s.Equals(testCase.expected.state) {
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
	test := func(s shogi.State, testCases []*testCase) {
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
				continue
			}
			t.Logf(result)
		}
	}
	// 打
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
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
			[2]shogi.Captured{
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
			},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 9, 1, 8, shogi.BKY)}, []string{"▲1八香"}},
				{[]*shogi.Move{move(0, 0, 1, 8, shogi.BKY)}, []string{"▲1八香打"}},
				{[]*shogi.Move{move(0, 0, 2, 7, shogi.BKY)}, []string{"▲2七香"}},
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
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 1, 9, 2, shogi.WKY)}, []string{"△9二香"}},
				{[]*shogi.Move{move(0, 0, 9, 2, shogi.WKY)}, []string{"△9二香打"}},
				{[]*shogi.Move{move(0, 0, 8, 3, shogi.WKY)}, []string{"△8三香"}},
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
			test(sw, testCases)
		}
	}
	// 成・不成
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.BUM, shogi.BNG, shogi.BNK, shogi.BNY, shogi.BTO},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BHI, shogi.BKA, shogi.BGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.BKA, shogi.BGI, shogi.EMP, shogi.BKY, shogi.BFU},
				{shogi.EMP, shogi.EMP, shogi.WKE, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKE, shogi.EMP, shogi.EMP},
				{shogi.WFU, shogi.WKY, shogi.EMP, shogi.WGI, shogi.WKA, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WGI, shogi.WKA, shogi.WHI},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.WTO, shogi.WNY, shogi.WNK, shogi.WNG, shogi.WUM, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 4, 1, 3, shogi.BFU)}, []string{"▲1三歩不成"}},
				{[]*shogi.Move{move(1, 4, 1, 3, shogi.BTO)}, []string{"▲1三歩成"}},
				{[]*shogi.Move{move(1, 1, 1, 2, shogi.BTO)}, []string{"▲1二と"}},
				{[]*shogi.Move{move(2, 4, 2, 3, shogi.BKY)}, []string{"▲2三香不成"}},
				{[]*shogi.Move{move(2, 4, 2, 3, shogi.BNY)}, []string{"▲2三香成"}},
				{[]*shogi.Move{move(2, 1, 2, 2, shogi.BNY)}, []string{"▲2二成香"}},
				{[]*shogi.Move{move(3, 5, 2, 3, shogi.BKE)}, []string{"▲2三桂不成"}},
				{[]*shogi.Move{move(3, 5, 2, 3, shogi.BNK)}, []string{"▲2三桂成"}},
				{[]*shogi.Move{move(3, 1, 3, 2, shogi.BNK)}, []string{"▲3二成桂"}},
				{[]*shogi.Move{move(4, 4, 4, 3, shogi.BGI)}, []string{"▲4三銀不成"}},
				{[]*shogi.Move{move(4, 4, 4, 3, shogi.BNG)}, []string{"▲4三銀成"}},
				{[]*shogi.Move{move(4, 1, 4, 2, shogi.BNG)}, []string{"▲4二成銀"}},
				{[]*shogi.Move{move(7, 3, 6, 4, shogi.BGI)}, []string{"▲6四銀不成"}},
				{[]*shogi.Move{move(7, 3, 6, 4, shogi.BNG)}, []string{"▲6四銀成"}},
				{[]*shogi.Move{move(5, 4, 4, 3, shogi.BKA)}, []string{"▲4三角不成"}},
				{[]*shogi.Move{move(5, 4, 4, 3, shogi.BUM)}, []string{"▲4三角成"}},
				{[]*shogi.Move{move(5, 1, 5, 2, shogi.BUM)}, []string{"▲5二馬"}},
				{[]*shogi.Move{move(8, 3, 7, 4, shogi.BKA)}, []string{"▲7四角不成"}},
				{[]*shogi.Move{move(8, 3, 7, 4, shogi.BUM)}, []string{"▲7四角成"}},
				{[]*shogi.Move{move(6, 4, 6, 3, shogi.BHI)}, []string{"▲6三飛不成"}},
				{[]*shogi.Move{move(6, 4, 6, 3, shogi.BRY)}, []string{"▲6三飛成"}},
				{[]*shogi.Move{move(6, 1, 6, 2, shogi.BRY)}, []string{"▲6二竜"}},
				{[]*shogi.Move{move(9, 3, 9, 5, shogi.BHI)}, []string{"▲9五飛不成"}},
				{[]*shogi.Move{move(9, 3, 9, 5, shogi.BRY)}, []string{"▲9五飛成"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 6, 9, 7, shogi.WFU)}, []string{"△9七歩不成"}},
				{[]*shogi.Move{move(9, 6, 9, 7, shogi.WTO)}, []string{"△9七歩成"}},
				{[]*shogi.Move{move(9, 9, 9, 8, shogi.WTO)}, []string{"△9八と"}},
				{[]*shogi.Move{move(8, 6, 8, 7, shogi.WKY)}, []string{"△8七香不成"}},
				{[]*shogi.Move{move(8, 6, 8, 7, shogi.WNY)}, []string{"△8七香成"}},
				{[]*shogi.Move{move(8, 9, 8, 8, shogi.WNY)}, []string{"△8八成香"}},
				{[]*shogi.Move{move(7, 5, 8, 7, shogi.WKE)}, []string{"△8七桂不成"}},
				{[]*shogi.Move{move(7, 5, 8, 7, shogi.WNK)}, []string{"△8七桂成"}},
				{[]*shogi.Move{move(7, 9, 7, 8, shogi.WNK)}, []string{"△7八成桂"}},
				{[]*shogi.Move{move(6, 6, 6, 7, shogi.WGI)}, []string{"△6七銀不成"}},
				{[]*shogi.Move{move(6, 6, 6, 7, shogi.WNG)}, []string{"△6七銀成"}},
				{[]*shogi.Move{move(6, 9, 6, 8, shogi.WNG)}, []string{"△6八成銀"}},
				{[]*shogi.Move{move(3, 7, 4, 6, shogi.WGI)}, []string{"△4六銀不成"}},
				{[]*shogi.Move{move(3, 7, 4, 6, shogi.WNG)}, []string{"△4六銀成"}},
				{[]*shogi.Move{move(5, 6, 6, 7, shogi.WKA)}, []string{"△6七角不成"}},
				{[]*shogi.Move{move(5, 6, 6, 7, shogi.WUM)}, []string{"△6七角成"}},
				{[]*shogi.Move{move(5, 9, 5, 8, shogi.WUM)}, []string{"△5八馬"}},
				{[]*shogi.Move{move(2, 7, 3, 6, shogi.WKA)}, []string{"△3六角不成"}},
				{[]*shogi.Move{move(2, 7, 3, 6, shogi.WUM)}, []string{"△3六角成"}},
				{[]*shogi.Move{move(4, 6, 4, 7, shogi.WHI)}, []string{"△4七飛不成"}},
				{[]*shogi.Move{move(4, 6, 4, 7, shogi.WRY)}, []string{"△4七飛成"}},
				{[]*shogi.Move{move(4, 9, 4, 8, shogi.WRY)}, []string{"△4八竜"}},
				{[]*shogi.Move{move(1, 7, 1, 5, shogi.WHI)}, []string{"△1五飛不成"}},
				{[]*shogi.Move{move(1, 7, 1, 5, shogi.WRY)}, []string{"△1五飛成"}},
			}
			test(sw, testCases)
		}
	}
	// 上・寄・引
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.WGI, shogi.EMP, shogi.EMP, shogi.BKI, shogi.WGI, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.BKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BKI, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKI, shogi.WGI, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKI, shogi.EMP, shogi.BKI, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.BGI, shogi.WKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.WKI},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKI, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.BGI, shogi.WKI, shogi.EMP, shogi.EMP, shogi.BGI, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 3, 8, 2, shogi.BKI)}, []string{"▲8二金上"}},
				{[]*shogi.Move{move(7, 2, 8, 2, shogi.BKI)}, []string{"▲8二金寄"}},
				{[]*shogi.Move{move(4, 3, 3, 2, shogi.BKI)}, []string{"▲3二金上"}},
				{[]*shogi.Move{move(3, 1, 3, 2, shogi.BKI)}, []string{"▲3二金引"}},
				{[]*shogi.Move{move(5, 6, 5, 5, shogi.BKI)}, []string{"▲5五金上"}},
				{[]*shogi.Move{move(4, 5, 5, 5, shogi.BKI)}, []string{"▲5五金寄"}},
				{[]*shogi.Move{move(8, 9, 8, 8, shogi.BGI)}, []string{"▲8八銀上"}},
				{[]*shogi.Move{move(7, 7, 8, 8, shogi.BGI)}, []string{"▲8八銀引"}},
				{[]*shogi.Move{move(4, 9, 3, 8, shogi.BGI)}, []string{"▲3八銀上"}},
				{[]*shogi.Move{move(2, 7, 3, 8, shogi.BGI)}, []string{"▲3八銀引"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 7, 2, 8, shogi.WKI)}, []string{"△2八金上"}},
				{[]*shogi.Move{move(3, 8, 2, 8, shogi.WKI)}, []string{"△2八金寄"}},
				{[]*shogi.Move{move(6, 7, 7, 8, shogi.WKI)}, []string{"△7八金上"}},
				{[]*shogi.Move{move(7, 9, 7, 8, shogi.WKI)}, []string{"△7八金引"}},
				{[]*shogi.Move{move(5, 4, 5, 5, shogi.WKI)}, []string{"△5五金上"}},
				{[]*shogi.Move{move(6, 5, 5, 5, shogi.WKI)}, []string{"△5五金寄"}},
				{[]*shogi.Move{move(2, 1, 2, 2, shogi.WGI)}, []string{"△2二銀上"}},
				{[]*shogi.Move{move(3, 3, 2, 2, shogi.WGI)}, []string{"△2二銀引"}},
				{[]*shogi.Move{move(6, 1, 7, 2, shogi.WGI)}, []string{"△7二銀上"}},
				{[]*shogi.Move{move(8, 3, 7, 2, shogi.WGI)}, []string{"△7二銀引"}},
			}
			test(sw, testCases)
		}
	}
	// 左・右・直 (2枚)
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.WGI, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKI, shogi.WKI, shogi.EMP},
				{shogi.BKI, shogi.EMP, shogi.BKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKI, shogi.EMP, shogi.BKI},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.EMP, shogi.BGI, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.WGI, shogi.EMP, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.WKI, shogi.EMP, shogi.WKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKI, shogi.EMP, shogi.WKI},
				{shogi.EMP, shogi.BKI, shogi.BKI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.BGI, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 2, 8, 1, shogi.BKI)}, []string{"▲8一金左"}},
				{[]*shogi.Move{move(7, 2, 8, 1, shogi.BKI)}, []string{"▲8一金右"}},
				{[]*shogi.Move{move(3, 2, 2, 2, shogi.BKI)}, []string{"▲2二金左"}},
				{[]*shogi.Move{move(1, 2, 2, 2, shogi.BKI)}, []string{"▲2二金右"}},
				{[]*shogi.Move{move(6, 4, 5, 5, shogi.BGI)}, []string{"▲5五銀左"}},
				{[]*shogi.Move{move(4, 4, 5, 5, shogi.BGI)}, []string{"▲5五銀右"}},
				{[]*shogi.Move{move(8, 9, 7, 8, shogi.BKI)}, []string{"▲7八金左"}},
				{[]*shogi.Move{move(7, 9, 7, 8, shogi.BKI)}, []string{"▲7八金直"}},
				{[]*shogi.Move{move(3, 9, 3, 8, shogi.BGI)}, []string{"▲3八銀直"}},
				{[]*shogi.Move{move(2, 9, 3, 8, shogi.BGI)}, []string{"▲3八銀右"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 8, 2, 9, shogi.WKI)}, []string{"△2九金左"}},
				{[]*shogi.Move{move(3, 8, 2, 9, shogi.WKI)}, []string{"△2九金右"}},
				{[]*shogi.Move{move(7, 8, 8, 8, shogi.WKI)}, []string{"△8八金左"}},
				{[]*shogi.Move{move(9, 8, 8, 8, shogi.WKI)}, []string{"△8八金右"}},
				{[]*shogi.Move{move(4, 6, 5, 5, shogi.WGI)}, []string{"△5五銀左"}},
				{[]*shogi.Move{move(6, 6, 5, 5, shogi.WGI)}, []string{"△5五銀右"}},
				{[]*shogi.Move{move(2, 1, 3, 2, shogi.WKI)}, []string{"△3二金左"}},
				{[]*shogi.Move{move(3, 1, 3, 2, shogi.WKI)}, []string{"△3二金直"}},
				{[]*shogi.Move{move(7, 1, 7, 2, shogi.WGI)}, []string{"△7二銀直"}},
				{[]*shogi.Move{move(8, 1, 7, 2, shogi.WGI)}, []string{"△7二銀右"}},
			}
			test(sw, testCases)
		}
	}
	// 左・右・直 (3枚以上)
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.WGI, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WTO, shogi.WTO, shogi.WTO},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WTO},
				{shogi.WGI, shogi.EMP, shogi.WGI, shogi.BKI, shogi.BKI, shogi.BKI, shogi.EMP, shogi.WTO, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.BTO, shogi.EMP, shogi.WKI, shogi.WKI, shogi.WKI, shogi.BGI, shogi.EMP, shogi.BGI},
				{shogi.BTO, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BTO, shogi.BTO, shogi.BTO, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.BGI, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(6, 3, 5, 2, shogi.BKI)}, []string{"▲5二金左"}},
				{[]*shogi.Move{move(5, 3, 5, 2, shogi.BKI)}, []string{"▲5二金直"}},
				{[]*shogi.Move{move(4, 3, 5, 2, shogi.BKI)}, []string{"▲5二金右"}},
				{[]*shogi.Move{move(7, 9, 8, 8, shogi.BTO)}, []string{"▲8八と右"}},
				{[]*shogi.Move{move(8, 9, 8, 8, shogi.BTO)}, []string{"▲8八と直"}},
				{[]*shogi.Move{move(9, 9, 8, 8, shogi.BTO)}, []string{"▲8八と左上"}},
				{[]*shogi.Move{move(9, 8, 8, 8, shogi.BTO)}, []string{"▲8八と寄"}},
				{[]*shogi.Move{move(8, 7, 8, 8, shogi.BTO)}, []string{"▲8八と引"}},
				{[]*shogi.Move{move(2, 9, 2, 8, shogi.BGI)}, []string{"▲2八銀直"}},
				{[]*shogi.Move{move(1, 7, 2, 8, shogi.BGI)}, []string{"▲2八銀右"}},
				{[]*shogi.Move{move(3, 9, 2, 8, shogi.BGI)}, []string{"▲2八銀左上"}},
				{[]*shogi.Move{move(3, 7, 2, 8, shogi.BGI)}, []string{"▲2八銀左引"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(4, 7, 5, 8, shogi.WKI)}, []string{"△5八金左"}},
				{[]*shogi.Move{move(5, 7, 5, 8, shogi.WKI)}, []string{"△5八金直"}},
				{[]*shogi.Move{move(6, 7, 5, 8, shogi.WKI)}, []string{"△5八金右"}},
				{[]*shogi.Move{move(3, 1, 2, 2, shogi.WTO)}, []string{"△2二と右"}},
				{[]*shogi.Move{move(2, 1, 2, 2, shogi.WTO)}, []string{"△2二と直"}},
				{[]*shogi.Move{move(1, 1, 2, 2, shogi.WTO)}, []string{"△2二と左上"}},
				{[]*shogi.Move{move(1, 2, 2, 2, shogi.WTO)}, []string{"△2二と寄"}},
				{[]*shogi.Move{move(2, 3, 2, 2, shogi.WTO)}, []string{"△2二と引"}},
				{[]*shogi.Move{move(8, 1, 8, 2, shogi.WGI)}, []string{"△8二銀直"}},
				{[]*shogi.Move{move(9, 3, 8, 2, shogi.WGI)}, []string{"△8二銀右"}},
				{[]*shogi.Move{move(7, 1, 8, 2, shogi.WGI)}, []string{"△8二銀左上"}},
				{[]*shogi.Move{move(7, 3, 8, 2, shogi.WGI)}, []string{"△8二銀左引"}},
			}
			test(sw, testCases)
		}
	}
	// 竜
	{
		{
			sb := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.BRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.BRY, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				},
				[2]shogi.Captured{},
				shogi.TurnBlack,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 1, 8, 2, shogi.BRY)}, []string{"▲8二竜引"}},
				{[]*shogi.Move{move(8, 4, 8, 2, shogi.BRY)}, []string{"▲8二竜上"}},
				{[]*shogi.Move{move(2, 5, 4, 5, shogi.BRY)}, []string{"▲4五竜寄"}},
				{[]*shogi.Move{move(5, 4, 4, 5, shogi.BRY)}, []string{"▲4五竜引"}},
				{[]*shogi.Move{move(5, 7, 3, 7, shogi.BRY)}, []string{"▲3七竜左"}},
				{[]*shogi.Move{move(1, 7, 3, 7, shogi.BRY)}, []string{"▲3七竜右"}},
			}
			test(sb, testCases)
		}
		{
			sb := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BRY},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.BRY, shogi.BRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				},
				[2]shogi.Captured{},
				shogi.TurnBlack,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 9, 8, 8, shogi.BRY)}, []string{"▲8八竜左"}},
				{[]*shogi.Move{move(8, 9, 8, 8, shogi.BRY)}, []string{"▲8八竜右"}},
				{[]*shogi.Move{move(2, 6, 1, 5, shogi.BRY)}, []string{"▲1五竜左"}},
				{[]*shogi.Move{move(1, 7, 1, 5, shogi.BRY)}, []string{"▲1五竜右"}},
				{[]*shogi.Move{move(3, 4, 3, 2, shogi.BRY)}, []string{"▲3二竜左"}},
				{[]*shogi.Move{move(2, 3, 3, 2, shogi.BRY)}, []string{"▲3二竜右"}},
			}
			test(sb, testCases)
		}
		{
			sw := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.WRY, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WRY},
				},
				[2]shogi.Captured{},
				shogi.TurnWhite,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 9, 2, 8, shogi.WRY)}, []string{"△2八竜引"}},
				{[]*shogi.Move{move(2, 6, 2, 8, shogi.WRY)}, []string{"△2八竜上"}},
				{[]*shogi.Move{move(8, 5, 6, 5, shogi.WRY)}, []string{"△6五竜寄"}},
				{[]*shogi.Move{move(5, 6, 6, 5, shogi.WRY)}, []string{"△6五竜引"}},
				{[]*shogi.Move{move(5, 3, 7, 3, shogi.WRY)}, []string{"△7三竜左"}},
				{[]*shogi.Move{move(9, 3, 7, 3, shogi.WRY)}, []string{"△7三竜右"}},
			}
			test(sw, testCases)
		}
		{
			sw := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WRY, shogi.WRY},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.WRY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				},
				[2]shogi.Captured{},
				shogi.TurnWhite,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 1, 2, 2, shogi.WRY)}, []string{"△2二竜左"}},
				{[]*shogi.Move{move(2, 1, 2, 2, shogi.WRY)}, []string{"△2二竜右"}},
				{[]*shogi.Move{move(8, 4, 9, 5, shogi.WRY)}, []string{"△9五竜左"}},
				{[]*shogi.Move{move(9, 3, 9, 5, shogi.WRY)}, []string{"△9五竜右"}},
				{[]*shogi.Move{move(7, 6, 7, 8, shogi.WRY)}, []string{"△7八竜左"}},
				{[]*shogi.Move{move(8, 7, 7, 8, shogi.WRY)}, []string{"△7八竜右"}},
			}
			test(sw, testCases)
		}
	}
	// 馬
	{
		{
			sb := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.BUM, shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM, shogi.EMP, shogi.EMP},
					{shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM},
					{shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				},
				[2]shogi.Captured{},
				shogi.TurnBlack,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 1, 8, 2, shogi.BUM)}, []string{"▲8二馬左"}},
				{[]*shogi.Move{move(8, 1, 8, 2, shogi.BUM)}, []string{"▲8二馬右"}},
				{[]*shogi.Move{move(9, 5, 8, 5, shogi.BUM)}, []string{"▲8五馬寄"}},
				{[]*shogi.Move{move(6, 3, 8, 5, shogi.BUM)}, []string{"▲8五馬引"}},
				{[]*shogi.Move{move(1, 1, 1, 2, shogi.BUM)}, []string{"▲1二馬引"}},
				{[]*shogi.Move{move(3, 4, 1, 2, shogi.BUM)}, []string{"▲1二馬上"}},
				{[]*shogi.Move{move(9, 9, 7, 7, shogi.BUM)}, []string{"▲7七馬左"}},
				{[]*shogi.Move{move(5, 9, 7, 7, shogi.BUM)}, []string{"▲7七馬右"}},
				{[]*shogi.Move{move(4, 7, 2, 9, shogi.BUM)}, []string{"▲2九馬左"}},
				{[]*shogi.Move{move(1, 8, 2, 9, shogi.BUM)}, []string{"▲2九馬右"}},
			}
			test(sb, testCases)
		}
		{
			sw := logic.NewState(
				[9][9]shogi.Piece{
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM},
					{shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM},
					{shogi.EMP, shogi.EMP, shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					{shogi.WUM, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WUM, shogi.WUM},
				},
				[2]shogi.Captured{},
				shogi.TurnWhite,
			)
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 9, 2, 8, shogi.WUM)}, []string{"△2八馬左"}},
				{[]*shogi.Move{move(2, 9, 2, 8, shogi.WUM)}, []string{"△2八馬右"}},
				{[]*shogi.Move{move(1, 5, 2, 5, shogi.WUM)}, []string{"△2五馬寄"}},
				{[]*shogi.Move{move(4, 7, 2, 5, shogi.WUM)}, []string{"△2五馬引"}},
				{[]*shogi.Move{move(9, 9, 9, 8, shogi.WUM)}, []string{"△9八馬引"}},
				{[]*shogi.Move{move(7, 6, 9, 8, shogi.WUM)}, []string{"△9八馬上"}},
				{[]*shogi.Move{move(1, 1, 3, 3, shogi.WUM)}, []string{"△3三馬左"}},
				{[]*shogi.Move{move(5, 1, 3, 3, shogi.WUM)}, []string{"△3三馬右"}},
				{[]*shogi.Move{move(6, 3, 8, 1, shogi.WUM)}, []string{"△8一馬左"}},
				{[]*shogi.Move{move(9, 2, 8, 1, shogi.WUM)}, []string{"△8一馬右"}},
			}
			test(sw, testCases)
		}
	}
	// 上成・引成・左成・右成
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
				{shogi.EMP, shogi.WKE, shogi.EMP, shogi.WKE, shogi.EMP, shogi.BKE, shogi.EMP, shogi.BKE, shogi.EMP},
				{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.WHI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(2, 2, 2, 3, shogi.BRY)}, []string{"▲2三飛引成"}},
				{[]*shogi.Move{move(2, 4, 2, 3, shogi.BRY)}, []string{"▲2三飛上成"}},
				{[]*shogi.Move{move(4, 5, 3, 3, shogi.BNK)}, []string{"▲3三桂左成"}},
				{[]*shogi.Move{move(2, 5, 3, 3, shogi.BNK)}, []string{"▲3三桂右成"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(8, 8, 8, 7, shogi.WRY)}, []string{"△8七飛引成"}},
				{[]*shogi.Move{move(8, 6, 8, 7, shogi.WRY)}, []string{"△8七飛上成"}},
				{[]*shogi.Move{move(6, 5, 7, 7, shogi.WNK)}, []string{"△7七桂左成"}},
				{[]*shogi.Move{move(8, 5, 7, 7, shogi.WNK)}, []string{"△7七桂右成"}},
			}
			test(sw, testCases)
		}
	}
	// 玉
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BOU, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BOU},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.WOU, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.WOU, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(1, 4, 1, 3, shogi.BOU)}, []string{"▲1三玉"}},
				{[]*shogi.Move{move(3, 3, 2, 2, shogi.BOU)}, []string{"▲2二玉"}},
				{[]*shogi.Move{move(3, 3, 4, 4, shogi.BOU)}, []string{"▲4四玉"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{move(9, 6, 9, 7, shogi.WOU)}, []string{"△9七玉"}},
				{[]*shogi.Move{move(7, 7, 8, 8, shogi.WOU)}, []string{"△8八玉"}},
				{[]*shogi.Move{move(7, 7, 6, 6, shogi.WOU)}, []string{"△6六玉"}},
			}
			test(sw, testCases)
		}
	}
	// 同
	{
		state := logic.NewState(
			[9][9]shogi.Piece{
				{shogi.WKY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKY},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.WFU},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.WFU, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BGI, shogi.BFU},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BFU, shogi.WGI, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BKY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKY},
			},
			[2]shogi.Captured{},
			shogi.TurnBlack,
		)
		sb := state.Clone()
		sb.SetTurn(shogi.TurnBlack)
		sw := state.Clone()
		sw.SetTurn(shogi.TurnWhite)
		{
			testCases := []*testCase{
				{[]*shogi.Move{
					move(1, 5, 1, 4, shogi.BFU),
					move(1, 3, 1, 4, shogi.WFU),
					move(1, 9, 1, 4, shogi.BKY),
				}, []string{"▲1四歩", "△同歩", "▲同香"}},
				{[]*shogi.Move{
					move(1, 5, 1, 4, shogi.BFU),
					move(1, 3, 1, 4, shogi.WFU),
					move(2, 5, 1, 4, shogi.BGI),
				}, []string{"▲1四歩", "△同歩", "▲同銀上"}},
				{[]*shogi.Move{
					move(1, 5, 1, 4, shogi.BFU),
					move(1, 3, 1, 4, shogi.WFU),
					move(2, 3, 1, 4, shogi.BGI),
				}, []string{"▲1四歩", "△同歩", "▲同銀引不成"}},
				{[]*shogi.Move{
					move(1, 5, 1, 4, shogi.BFU),
					move(1, 3, 1, 4, shogi.WFU),
					move(2, 3, 1, 4, shogi.BNG),
				}, []string{"▲1四歩", "△同歩", "▲同銀引成"}},
			}
			test(sb, testCases)
		}
		{
			testCases := []*testCase{
				{[]*shogi.Move{
					move(9, 5, 9, 6, shogi.WFU),
					move(9, 7, 9, 6, shogi.BFU),
					move(9, 1, 9, 6, shogi.WKY),
				}, []string{"△9六歩", "▲同歩", "△同香"}},
				{[]*shogi.Move{
					move(9, 5, 9, 6, shogi.WFU),
					move(9, 7, 9, 6, shogi.BFU),
					move(8, 5, 9, 6, shogi.WGI),
				}, []string{"△9六歩", "▲同歩", "△同銀上"}},
				{[]*shogi.Move{
					move(9, 5, 9, 6, shogi.WFU),
					move(9, 7, 9, 6, shogi.BFU),
					move(8, 7, 9, 6, shogi.WGI),
				}, []string{"△9六歩", "▲同歩", "△同銀引不成"}},
				{[]*shogi.Move{
					move(9, 5, 9, 6, shogi.WFU),
					move(9, 7, 9, 6, shogi.BFU),
					move(8, 7, 9, 6, shogi.WNG),
				}, []string{"△9六歩", "▲同歩", "△同銀引成"}},
			}
			test(sw, testCases)
		}
	}
}
