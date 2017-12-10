package shogi

import (
	"testing"
)

type testData struct {
	move     *Move
	expected string
}

func TestMoveString(t *testing.T) {
	// 打
	{
		state := NewState()
		state.SetBoard(1, 9, &BoardPiece{TurnBlack, KY})
		state.SetBoard(2, 9, &BoardPiece{TurnBlack, KE})
		state.SetBoard(3, 9, &BoardPiece{TurnBlack, GI})
		state.SetBoard(4, 9, &BoardPiece{TurnBlack, KI})
		state.SetBoard(8, 8, &BoardPiece{TurnBlack, KA})
		state.SetBoard(2, 8, &BoardPiece{TurnBlack, HI})
		state.SetBoard(9, 1, &BoardPiece{TurnWhite, KY})
		state.SetBoard(8, 1, &BoardPiece{TurnWhite, KE})
		state.SetBoard(7, 1, &BoardPiece{TurnWhite, GI})
		state.SetBoard(6, 1, &BoardPiece{TurnWhite, KI})
		state.SetBoard(2, 2, &BoardPiece{TurnWhite, KA})
		state.SetBoard(8, 2, &BoardPiece{TurnWhite, HI})
		state.Captured[TurnBlack] = &CapturedPieces{KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1}
		state.Captured[TurnWhite] = &CapturedPieces{KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1}
		tests := []testData{
			// KY
			testData{
				move:     &Move{TurnBlack, Pos(1, 9), Pos(1, 8), KY},
				expected: "▲1八香",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(1, 8), KY},
				expected: "▲1八香打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(2, 8), KY},
				expected: "▲2八香",
			},
			testData{
				move:     &Move{TurnWhite, Pos(9, 1), Pos(9, 2), KY},
				expected: "△9二香",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(9, 2), KY},
				expected: "△9二香打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(8, 2), KY},
				expected: "△8二香",
			},
			// KE
			testData{
				move:     &Move{TurnBlack, Pos(2, 9), Pos(1, 7), KE},
				expected: "▲1七桂",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(1, 7), KE},
				expected: "▲1七桂打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(2, 7), KE},
				expected: "▲2七桂",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 1), Pos(9, 3), KE},
				expected: "△9三桂",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(9, 3), KE},
				expected: "△9三桂打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(8, 3), KE},
				expected: "△8三桂",
			},
			// GI
			testData{
				move:     &Move{TurnBlack, Pos(3, 9), Pos(3, 8), GI},
				expected: "▲3八銀",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(3, 8), GI},
				expected: "▲3八銀打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(3, 7), GI},
				expected: "▲3七銀",
			},
			testData{
				move:     &Move{TurnWhite, Pos(7, 1), Pos(7, 2), GI},
				expected: "△7二銀",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(7, 2), GI},
				expected: "△7二銀打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(7, 3), GI},
				expected: "△7三銀",
			},
			// KI
			testData{
				move:     &Move{TurnBlack, Pos(4, 9), Pos(5, 8), KI},
				expected: "▲5八金",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(5, 8), KI},
				expected: "▲5八金打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(5, 7), KI},
				expected: "▲5七金",
			},
			testData{
				move:     &Move{TurnWhite, Pos(4, 9), Pos(5, 2), KI},
				expected: "△5二金",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(5, 2), KI},
				expected: "△5二金打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(5, 3), KI},
				expected: "△5三金",
			},
			// KA
			testData{
				move:     &Move{TurnBlack, Pos(8, 8), Pos(7, 7), KA},
				expected: "▲7七角",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(7, 7), KA},
				expected: "▲7七角打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(8, 7), KA},
				expected: "▲8七角",
			},
			testData{
				move:     &Move{TurnWhite, Pos(2, 2), Pos(3, 3), KA},
				expected: "△3三角",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(3, 3), KA},
				expected: "△3三角打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(2, 3), KA},
				expected: "△2三角",
			},
			// HI
			testData{
				move:     &Move{TurnBlack, Pos(2, 8), Pos(2, 7), HI},
				expected: "▲2七飛",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(2, 7), HI},
				expected: "▲2七飛打",
			},
			testData{
				move:     &Move{TurnBlack, Pos(0, 0), Pos(3, 7), HI},
				expected: "▲3七飛",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 2), Pos(8, 3), HI},
				expected: "△8三飛",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(8, 3), HI},
				expected: "△8三飛打",
			},
			testData{
				move:     &Move{TurnWhite, Pos(0, 0), Pos(7, 3), HI},
				expected: "△7三飛",
			},
		}
		test(t, state, tests)
	}
	// 成・不成
	{
		state := NewState()
		state.SetBoard(1, 4, &BoardPiece{TurnBlack, FU})
		state.SetBoard(2, 4, &BoardPiece{TurnBlack, KY})
		state.SetBoard(3, 5, &BoardPiece{TurnBlack, KE})
		state.SetBoard(4, 4, &BoardPiece{TurnBlack, GI})
		state.SetBoard(5, 4, &BoardPiece{TurnBlack, KA})
		state.SetBoard(6, 4, &BoardPiece{TurnBlack, HI})
		state.SetBoard(7, 3, &BoardPiece{TurnBlack, GI})
		state.SetBoard(8, 3, &BoardPiece{TurnBlack, KA})
		state.SetBoard(9, 3, &BoardPiece{TurnBlack, HI})
		state.SetBoard(1, 1, &BoardPiece{TurnBlack, TO})
		state.SetBoard(2, 1, &BoardPiece{TurnBlack, NY})
		state.SetBoard(3, 1, &BoardPiece{TurnBlack, NK})
		state.SetBoard(4, 1, &BoardPiece{TurnBlack, NG})
		state.SetBoard(5, 1, &BoardPiece{TurnBlack, UM})
		state.SetBoard(6, 1, &BoardPiece{TurnBlack, RY})
		state.SetBoard(9, 6, &BoardPiece{TurnWhite, FU})
		state.SetBoard(8, 6, &BoardPiece{TurnWhite, KY})
		state.SetBoard(7, 5, &BoardPiece{TurnWhite, KE})
		state.SetBoard(6, 6, &BoardPiece{TurnWhite, GI})
		state.SetBoard(5, 6, &BoardPiece{TurnWhite, KA})
		state.SetBoard(4, 6, &BoardPiece{TurnWhite, HI})
		state.SetBoard(3, 7, &BoardPiece{TurnWhite, GI})
		state.SetBoard(2, 7, &BoardPiece{TurnWhite, KA})
		state.SetBoard(1, 7, &BoardPiece{TurnWhite, HI})
		state.SetBoard(9, 9, &BoardPiece{TurnWhite, TO})
		state.SetBoard(8, 9, &BoardPiece{TurnWhite, NY})
		state.SetBoard(7, 9, &BoardPiece{TurnWhite, NK})
		state.SetBoard(6, 9, &BoardPiece{TurnWhite, NG})
		state.SetBoard(5, 9, &BoardPiece{TurnWhite, UM})
		state.SetBoard(4, 9, &BoardPiece{TurnWhite, RY})
		tests := []testData{
			// FU
			testData{
				move:     &Move{TurnBlack, Pos(1, 4), Pos(1, 3), FU},
				expected: "▲1三歩不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(1, 4), Pos(1, 3), TO},
				expected: "▲1三歩成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(1, 1), Pos(1, 2), TO},
				expected: "▲1二と",
			},
			testData{
				move:     &Move{TurnWhite, Pos(9, 6), Pos(9, 7), FU},
				expected: "△9七歩不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(9, 6), Pos(9, 7), TO},
				expected: "△9七歩成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(9, 9), Pos(9, 8), TO},
				expected: "△9八と",
			},
			// KY
			testData{
				move:     &Move{TurnBlack, Pos(2, 4), Pos(2, 3), KY},
				expected: "▲2三香不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(2, 4), Pos(2, 3), NY},
				expected: "▲2三香成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(2, 1), Pos(2, 2), NY},
				expected: "▲2二成香",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 6), Pos(8, 7), KY},
				expected: "△8七香不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 6), Pos(8, 7), NY},
				expected: "△8七香成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 9), Pos(8, 8), NY},
				expected: "△8八成香",
			},
			// KE
			testData{
				move:     &Move{TurnBlack, Pos(3, 5), Pos(2, 3), KE},
				expected: "▲2三桂不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(3, 5), Pos(2, 3), NK},
				expected: "▲2三桂成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(3, 1), Pos(3, 2), NK},
				expected: "▲3二成桂",
			},
			testData{
				move:     &Move{TurnWhite, Pos(7, 5), Pos(8, 7), KE},
				expected: "△8七桂不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(7, 5), Pos(8, 7), NK},
				expected: "△8七桂成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(7, 9), Pos(7, 8), NK},
				expected: "△7八成桂",
			},
			// GI
			testData{
				move:     &Move{TurnBlack, Pos(4, 4), Pos(4, 3), GI},
				expected: "▲4三銀不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(4, 4), Pos(4, 3), NG},
				expected: "▲4三銀成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(4, 1), Pos(4, 2), NG},
				expected: "▲4二成銀",
			},
			testData{
				move:     &Move{TurnBlack, Pos(7, 3), Pos(8, 4), GI},
				expected: "▲8四銀不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(7, 3), Pos(8, 4), NG},
				expected: "▲8四銀成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 6), Pos(6, 7), GI},
				expected: "△6七銀不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 6), Pos(6, 7), NG},
				expected: "△6七銀成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 9), Pos(6, 8), NG},
				expected: "△6八成銀",
			},
			testData{
				move:     &Move{TurnWhite, Pos(3, 7), Pos(2, 6), GI},
				expected: "△2六銀不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(3, 7), Pos(2, 6), NG},
				expected: "△2六銀成",
			},
			// KA
			testData{
				move:     &Move{TurnBlack, Pos(5, 4), Pos(4, 3), KA},
				expected: "▲4三角不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(5, 4), Pos(4, 3), UM},
				expected: "▲4三角成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(5, 1), Pos(5, 2), UM},
				expected: "▲5二馬",
			},
			testData{
				move:     &Move{TurnBlack, Pos(8, 3), Pos(7, 4), KA},
				expected: "▲7四角不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(8, 3), Pos(7, 4), UM},
				expected: "▲7四角成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(5, 6), Pos(6, 7), KA},
				expected: "△6七角不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(5, 6), Pos(6, 7), UM},
				expected: "△6七角成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(5, 9), Pos(5, 8), UM},
				expected: "△5八馬",
			},
			testData{
				move:     &Move{TurnWhite, Pos(2, 7), Pos(3, 6), KA},
				expected: "△3六角不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(2, 7), Pos(3, 6), UM},
				expected: "△3六角成",
			},
			// HI
			testData{
				move:     &Move{TurnBlack, Pos(6, 4), Pos(6, 3), HI},
				expected: "▲6三飛不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(6, 4), Pos(6, 3), RY},
				expected: "▲6三飛成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(6, 1), Pos(6, 2), RY},
				expected: "▲6二竜",
			},
			testData{
				move:     &Move{TurnBlack, Pos(9, 3), Pos(9, 5), HI},
				expected: "▲9五飛不成",
			},
			testData{
				move:     &Move{TurnBlack, Pos(9, 3), Pos(9, 5), RY},
				expected: "▲9五飛成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(4, 6), Pos(4, 7), HI},
				expected: "△4七飛不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(4, 6), Pos(4, 7), RY},
				expected: "△4七飛成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(4, 9), Pos(4, 8), RY},
				expected: "△4八竜",
			},
			testData{
				move:     &Move{TurnWhite, Pos(1, 7), Pos(1, 5), HI},
				expected: "△1五飛不成",
			},
			testData{
				move:     &Move{TurnWhite, Pos(1, 7), Pos(1, 5), RY},
				expected: "△1五飛成",
			},
		}
		test(t, state, tests)
	}
	// 上・寄・引
	{
		// P1 *  *  * -GI *  * +KI-GI *
		// P2 *  * +KI *  *  *  *  *  *
		// P3+KI-GI *  *  * +KI-GI *  *
		// P4 *  *  *  * -KI *  *  *  *
		// P5 *  *  * -KI * +KI *  *  *
		// P6 *  *  *  * +KI *  *  *  *
		// P7 *  * +GI-KI *  *  * +GI-KI
		// P8 *  *  *  *  *  * -KI *  *
		// P9 * +GI-KI *  * +GI *  *  *
		state := NewState()
		state.SetBoard(9, 3, &BoardPiece{TurnBlack, KI})
		state.SetBoard(7, 2, &BoardPiece{TurnBlack, KI})
		state.SetBoard(4, 3, &BoardPiece{TurnBlack, KI})
		state.SetBoard(3, 1, &BoardPiece{TurnBlack, KI})
		state.SetBoard(5, 6, &BoardPiece{TurnBlack, KI})
		state.SetBoard(4, 5, &BoardPiece{TurnBlack, KI})
		state.SetBoard(8, 9, &BoardPiece{TurnBlack, GI})
		state.SetBoard(7, 7, &BoardPiece{TurnBlack, GI})
		state.SetBoard(4, 9, &BoardPiece{TurnBlack, GI})
		state.SetBoard(2, 7, &BoardPiece{TurnBlack, GI})
		state.SetBoard(1, 7, &BoardPiece{TurnWhite, KI})
		state.SetBoard(3, 8, &BoardPiece{TurnWhite, KI})
		state.SetBoard(6, 7, &BoardPiece{TurnWhite, KI})
		state.SetBoard(7, 9, &BoardPiece{TurnWhite, KI})
		state.SetBoard(5, 4, &BoardPiece{TurnWhite, KI})
		state.SetBoard(6, 5, &BoardPiece{TurnWhite, KI})
		state.SetBoard(2, 1, &BoardPiece{TurnWhite, GI})
		state.SetBoard(3, 3, &BoardPiece{TurnWhite, GI})
		state.SetBoard(6, 1, &BoardPiece{TurnWhite, GI})
		state.SetBoard(8, 3, &BoardPiece{TurnWhite, GI})
		tests := []testData{
			testData{
				move:     &Move{TurnBlack, Pos(9, 3), Pos(8, 2), KI},
				expected: "▲8二金上",
			},
			testData{
				move:     &Move{TurnBlack, Pos(7, 2), Pos(8, 2), KI},
				expected: "▲8二金寄",
			},
			testData{
				move:     &Move{TurnBlack, Pos(4, 3), Pos(3, 2), KI},
				expected: "▲3二金上",
			},
			testData{
				move:     &Move{TurnBlack, Pos(3, 1), Pos(3, 2), KI},
				expected: "▲3二金引",
			},
			testData{
				move:     &Move{TurnBlack, Pos(5, 6), Pos(5, 5), KI},
				expected: "▲5五金上",
			},
			testData{
				move:     &Move{TurnBlack, Pos(4, 5), Pos(5, 5), KI},
				expected: "▲5五金寄",
			},
			testData{
				move:     &Move{TurnBlack, Pos(8, 9), Pos(8, 8), GI},
				expected: "▲8八銀上",
			},
			testData{
				move:     &Move{TurnBlack, Pos(7, 7), Pos(8, 8), GI},
				expected: "▲8八銀引",
			},
			testData{
				move:     &Move{TurnBlack, Pos(4, 9), Pos(3, 8), GI},
				expected: "▲3八銀上",
			},
			testData{
				move:     &Move{TurnBlack, Pos(2, 7), Pos(3, 8), GI},
				expected: "▲3八銀引",
			},
			testData{
				move:     &Move{TurnWhite, Pos(1, 7), Pos(2, 8), KI},
				expected: "△2八金上",
			},
			testData{
				move:     &Move{TurnWhite, Pos(3, 8), Pos(2, 8), KI},
				expected: "△2八金寄",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 7), Pos(7, 8), KI},
				expected: "△7八金上",
			},
			testData{
				move:     &Move{TurnWhite, Pos(7, 9), Pos(7, 8), KI},
				expected: "△7八金引",
			},
			testData{
				move:     &Move{TurnWhite, Pos(5, 4), Pos(5, 5), KI},
				expected: "△5五金上",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 5), Pos(5, 5), KI},
				expected: "△5五金寄",
			},
			testData{
				move:     &Move{TurnWhite, Pos(2, 1), Pos(2, 2), GI},
				expected: "△2二銀上",
			},
			testData{
				move:     &Move{TurnWhite, Pos(3, 3), Pos(2, 2), GI},
				expected: "△2二銀引",
			},
			testData{
				move:     &Move{TurnWhite, Pos(6, 1), Pos(7, 2), GI},
				expected: "△7二銀上",
			},
			testData{
				move:     &Move{TurnWhite, Pos(8, 3), Pos(7, 2), GI},
				expected: "△7二銀引",
			},
		}
		test(t, state, tests)
	}
	// 左・右・直
	{
		{
			// P1 * -GI-GI *  *  * -KI-KI *
			// P2+KI * +KI *  *  * +KI * +KI
			// P3 *  *  *  *  *  *  *  *  *
			// P4 *  *  * +GI * +GI *  *  *
			// P5 *  *  *  *  *  *  *  *  *
			// P6 *  *  * -GI * -GI *  *  *
			// P7 *  *  *  *  *  *  *  *  *
			// P8-KI * -KI *  *  * -KI * -KI
			// P9 * +KI+KI *  *  * +GI+GI *
			state := NewState()
			state.SetBoard(9, 2, &BoardPiece{TurnBlack, KI})
			state.SetBoard(7, 2, &BoardPiece{TurnBlack, KI})
			state.SetBoard(3, 2, &BoardPiece{TurnBlack, KI})
			state.SetBoard(1, 2, &BoardPiece{TurnBlack, KI})
			state.SetBoard(6, 4, &BoardPiece{TurnBlack, GI})
			state.SetBoard(4, 4, &BoardPiece{TurnBlack, GI})
			state.SetBoard(8, 9, &BoardPiece{TurnBlack, KI})
			state.SetBoard(7, 9, &BoardPiece{TurnBlack, KI})
			state.SetBoard(3, 9, &BoardPiece{TurnBlack, GI})
			state.SetBoard(2, 9, &BoardPiece{TurnBlack, GI})
			state.SetBoard(1, 8, &BoardPiece{TurnWhite, KI})
			state.SetBoard(3, 8, &BoardPiece{TurnWhite, KI})
			state.SetBoard(7, 8, &BoardPiece{TurnWhite, KI})
			state.SetBoard(9, 8, &BoardPiece{TurnWhite, KI})
			state.SetBoard(6, 6, &BoardPiece{TurnWhite, GI})
			state.SetBoard(4, 6, &BoardPiece{TurnWhite, GI})
			state.SetBoard(2, 1, &BoardPiece{TurnWhite, KI})
			state.SetBoard(3, 1, &BoardPiece{TurnWhite, KI})
			state.SetBoard(7, 1, &BoardPiece{TurnWhite, GI})
			state.SetBoard(8, 1, &BoardPiece{TurnWhite, GI})
			tests := []testData{
				testData{
					move:     &Move{TurnBlack, Pos(9, 2), Pos(8, 1), KI},
					expected: "▲8一金左",
				},
				testData{
					move:     &Move{TurnBlack, Pos(7, 2), Pos(8, 1), KI},
					expected: "▲8一金右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(3, 2), Pos(2, 2), KI},
					expected: "▲2二金左",
				},
				testData{
					move:     &Move{TurnBlack, Pos(1, 2), Pos(2, 2), KI},
					expected: "▲2二金右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(6, 4), Pos(5, 5), GI},
					expected: "▲5五銀左",
				},
				testData{
					move:     &Move{TurnBlack, Pos(4, 4), Pos(5, 5), GI},
					expected: "▲5五銀右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(8, 9), Pos(7, 8), KI},
					expected: "▲7八金左",
				},
				testData{
					move:     &Move{TurnBlack, Pos(7, 9), Pos(7, 8), KI},
					expected: "▲7八金直",
				},
				testData{
					move:     &Move{TurnBlack, Pos(3, 9), Pos(3, 8), GI},
					expected: "▲3八銀直",
				},
				testData{
					move:     &Move{TurnBlack, Pos(2, 9), Pos(3, 8), GI},
					expected: "▲3八銀右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(1, 8), Pos(2, 9), KI},
					expected: "△2九金左",
				},
				testData{
					move:     &Move{TurnWhite, Pos(3, 8), Pos(2, 9), KI},
					expected: "△2九金右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(7, 8), Pos(8, 8), KI},
					expected: "△8八金左",
				},
				testData{
					move:     &Move{TurnWhite, Pos(9, 8), Pos(8, 8), KI},
					expected: "△8八金右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(4, 6), Pos(5, 5), GI},
					expected: "△5五銀左",
				},
				testData{
					move:     &Move{TurnWhite, Pos(6, 6), Pos(5, 5), GI},
					expected: "△5五銀右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(2, 1), Pos(3, 2), KI},
					expected: "△3二金左",
				},
				testData{
					move:     &Move{TurnWhite, Pos(3, 1), Pos(3, 2), KI},
					expected: "△3二金直",
				},
				testData{
					move:     &Move{TurnWhite, Pos(7, 1), Pos(7, 2), GI},
					expected: "△7二銀直",
				},
				testData{
					move:     &Move{TurnWhite, Pos(8, 1), Pos(7, 2), GI},
					expected: "△7二銀右",
				},
			}
			test(t, state, tests)
		}
		{
			// P1 * -GI-GI *  *  * -TO-TO-TO
			// P2 *  *  *  *  *  *  *  * -TO
			// P3-GI * -GI+KI+KI+KI * -TO *
			// P4 *  *  *  *  *  *  *  *  *
			// P5 *  *  *  *  *  *  *  *  *
			// P6 *  *  *  *  *  *  *  *  *
			// P7 * +TO * -KI-KI-KI+GI * +GI
			// P8+TO *  *  *  *  *  *  *  *
			// P9+TO+TO+TO *  *  * +GI+GI *
			state := NewState()
			state.SetBoard(6, 3, &BoardPiece{TurnBlack, KI})
			state.SetBoard(5, 3, &BoardPiece{TurnBlack, KI})
			state.SetBoard(4, 3, &BoardPiece{TurnBlack, KI})
			state.SetBoard(7, 9, &BoardPiece{TurnBlack, TO})
			state.SetBoard(8, 9, &BoardPiece{TurnBlack, TO})
			state.SetBoard(9, 9, &BoardPiece{TurnBlack, TO})
			state.SetBoard(9, 8, &BoardPiece{TurnBlack, TO})
			state.SetBoard(8, 7, &BoardPiece{TurnBlack, TO})
			state.SetBoard(2, 9, &BoardPiece{TurnBlack, GI})
			state.SetBoard(1, 7, &BoardPiece{TurnBlack, GI})
			state.SetBoard(3, 9, &BoardPiece{TurnBlack, GI})
			state.SetBoard(3, 7, &BoardPiece{TurnBlack, GI})
			state.SetBoard(4, 7, &BoardPiece{TurnWhite, KI})
			state.SetBoard(5, 7, &BoardPiece{TurnWhite, KI})
			state.SetBoard(6, 7, &BoardPiece{TurnWhite, KI})
			state.SetBoard(3, 1, &BoardPiece{TurnWhite, TO})
			state.SetBoard(2, 1, &BoardPiece{TurnWhite, TO})
			state.SetBoard(1, 1, &BoardPiece{TurnWhite, TO})
			state.SetBoard(1, 2, &BoardPiece{TurnWhite, TO})
			state.SetBoard(2, 3, &BoardPiece{TurnWhite, TO})
			state.SetBoard(8, 1, &BoardPiece{TurnWhite, GI})
			state.SetBoard(9, 3, &BoardPiece{TurnWhite, GI})
			state.SetBoard(7, 1, &BoardPiece{TurnWhite, GI})
			state.SetBoard(7, 3, &BoardPiece{TurnWhite, GI})
			tests := []testData{
				testData{
					move:     &Move{TurnBlack, Pos(6, 3), Pos(5, 2), KI},
					expected: "▲5二金左",
				},
				testData{
					move:     &Move{TurnBlack, Pos(5, 3), Pos(5, 2), KI},
					expected: "▲5二金直",
				},
				testData{
					move:     &Move{TurnBlack, Pos(4, 3), Pos(5, 2), KI},
					expected: "▲5二金右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(7, 9), Pos(8, 8), TO},
					expected: "▲8八と右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(8, 9), Pos(8, 8), TO},
					expected: "▲8八と直",
				},
				testData{
					move:     &Move{TurnBlack, Pos(9, 9), Pos(8, 8), TO},
					expected: "▲8八と左上",
				},
				testData{
					move:     &Move{TurnBlack, Pos(9, 8), Pos(8, 8), TO},
					expected: "▲8八と寄",
				},
				testData{
					move:     &Move{TurnBlack, Pos(8, 7), Pos(8, 8), TO},
					expected: "▲8八と引",
				},
				testData{
					move:     &Move{TurnBlack, Pos(2, 9), Pos(2, 8), GI},
					expected: "▲2八銀直",
				},
				testData{
					move:     &Move{TurnBlack, Pos(1, 7), Pos(2, 8), GI},
					expected: "▲2八銀右",
				},
				testData{
					move:     &Move{TurnBlack, Pos(3, 9), Pos(2, 8), GI},
					expected: "▲2八銀左上",
				},
				testData{
					move:     &Move{TurnBlack, Pos(3, 7), Pos(2, 8), GI},
					expected: "▲2八銀左引",
				},
				testData{
					move:     &Move{TurnWhite, Pos(4, 7), Pos(5, 8), KI},
					expected: "△5八金左",
				},
				testData{
					move:     &Move{TurnWhite, Pos(5, 7), Pos(5, 8), KI},
					expected: "△5八金直",
				},
				testData{
					move:     &Move{TurnWhite, Pos(6, 7), Pos(5, 8), KI},
					expected: "△5八金右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(3, 1), Pos(2, 2), TO},
					expected: "△2二と右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(2, 1), Pos(2, 2), TO},
					expected: "△2二と直",
				},
				testData{
					move:     &Move{TurnWhite, Pos(1, 1), Pos(2, 2), TO},
					expected: "△2二と左上",
				},
				testData{
					move:     &Move{TurnWhite, Pos(1, 2), Pos(2, 2), TO},
					expected: "△2二と寄",
				},
				testData{
					move:     &Move{TurnWhite, Pos(2, 3), Pos(2, 2), TO},
					expected: "△2二と引",
				},
				testData{
					move:     &Move{TurnWhite, Pos(8, 1), Pos(8, 2), GI},
					expected: "△8二銀直",
				},
				testData{
					move:     &Move{TurnWhite, Pos(9, 3), Pos(8, 2), GI},
					expected: "△8二銀右",
				},
				testData{
					move:     &Move{TurnWhite, Pos(7, 1), Pos(8, 2), GI},
					expected: "△8二銀左上",
				},
				testData{
					move:     &Move{TurnWhite, Pos(7, 3), Pos(8, 2), GI},
					expected: "△8二銀左引",
				},
			}
			test(t, state, tests)
		}
		{
			{
				// P1+RY *  *  *  *  *  *  *  *
				// P2 *  *  *  * +RY *  *  *  *
				// P3 *  *  *  *  *  *  * +RY *
				// P4 * +RY *  *  *  *  *  *  *
				// P5 *  *  *  * +RY *  *  * +RY
				// P6 *  *  *  *  *  *  *  *  *
				// P7 *  *  *  *  *  *  *  *  *
				// P8 *  *  *  *  *  *  * +RY *
				// P9+RY+RY *  *  *  *  *  * +RY
				{
					state := NewState()
					state.SetBoard(9, 1, &BoardPiece{TurnBlack, RY})
					state.SetBoard(8, 4, &BoardPiece{TurnBlack, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(9, 1), Pos(8, 2), RY},
							expected: "▲8二竜引",
						},
						testData{
							move:     &Move{TurnBlack, Pos(8, 4), Pos(8, 2), RY},
							expected: "▲8二竜上",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(2, 3, &BoardPiece{TurnBlack, RY})
					state.SetBoard(5, 2, &BoardPiece{TurnBlack, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(2, 3), Pos(4, 3), RY},
							expected: "▲4三竜寄",
						},
						testData{
							move:     &Move{TurnBlack, Pos(5, 2), Pos(4, 3), RY},
							expected: "▲4三竜引",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(5, 5, &BoardPiece{TurnBlack, RY})
					state.SetBoard(1, 5, &BoardPiece{TurnBlack, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(5, 5), Pos(3, 5), RY},
							expected: "▲3五竜左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(1, 5), Pos(3, 5), RY},
							expected: "▲3五竜右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(9, 9, &BoardPiece{TurnBlack, RY})
					state.SetBoard(8, 9, &BoardPiece{TurnBlack, RY})
					tests := []testData{
						{
							move:     &Move{TurnBlack, Pos(9, 9), Pos(8, 8), RY},
							expected: "▲8八竜左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(8, 9), Pos(8, 8), RY},
							expected: "▲8八竜右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(2, 8, &BoardPiece{TurnBlack, RY})
					state.SetBoard(1, 9, &BoardPiece{TurnBlack, RY})
					tests := []testData{
						{
							move:     &Move{TurnBlack, Pos(2, 8), Pos(1, 7), RY},
							expected: "▲1七竜左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(1, 9), Pos(1, 7), RY},
							expected: "▲1七竜右",
						},
					}
					test(t, state, tests)
				}
			}
			{
				// P1-RY *  *  *  *  *  * -RY-RY
				// P2 * -RY *  *  *  *  *  *  *
				// P3 *  *  *  *  *  *  *  *  *
				// P4 *  *  *  *  *  *  *  *  *
				// P5-RY * -RY *  *  *  *  *  *
				// P6 *  *  *  *  *  *  * -RY *
				// P7 * -RY *  *  *  *  *  *  *
				// P8 *  *  *  * -RY *  *  *  *
				// P9 *  *  *  *  *  *  *  * -RY
				{
					state := NewState()
					state.SetBoard(1, 9, &BoardPiece{TurnWhite, RY})
					state.SetBoard(2, 6, &BoardPiece{TurnWhite, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(1, 9), Pos(2, 8), RY},
							expected: "△2八竜引",
						},
						testData{
							move:     &Move{TurnWhite, Pos(2, 6), Pos(2, 8), RY},
							expected: "△2八竜上",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(8, 7, &BoardPiece{TurnWhite, RY})
					state.SetBoard(5, 8, &BoardPiece{TurnWhite, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(8, 7), Pos(6, 7), RY},
							expected: "△6七竜寄",
						},
						testData{
							move:     &Move{TurnWhite, Pos(5, 8), Pos(6, 7), RY},
							expected: "△6七竜引",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(5, 5, &BoardPiece{TurnWhite, RY})
					state.SetBoard(9, 5, &BoardPiece{TurnWhite, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(5, 5), Pos(7, 5), RY},
							expected: "△7五竜左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(9, 5), Pos(7, 5), RY},
							expected: "△7五竜右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(1, 1, &BoardPiece{TurnWhite, RY})
					state.SetBoard(2, 1, &BoardPiece{TurnWhite, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(1, 1), Pos(2, 2), RY},
							expected: "△2二竜左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(2, 1), Pos(2, 2), RY},
							expected: "△2二竜右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(8, 2, &BoardPiece{TurnWhite, RY})
					state.SetBoard(9, 1, &BoardPiece{TurnWhite, RY})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(8, 2), Pos(9, 3), RY},
							expected: "△9三竜左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(9, 1), Pos(9, 3), RY},
							expected: "△9三竜右",
						},
					}
					test(t, state, tests)
				}
			}
		}
		{
			{
				// P1+UM+UM *  *  *  *  *  * +UM
				// P2 *  *  *  *  *  *  *  *  *
				// P3 *  *  * +UM *  *  *  *  *
				// P4 *  *  *  *  *  * +UM *  *
				// P5+UM *  *  *  *  *  *  *  *
				// P6 *  *  *  *  *  *  *  *  *
				// P7 *  *  *  *  * +UM *  *  *
				// P8 *  *  *  *  *  *  *  * +UM
				// P9+UM *  *  * +UM *  *  *  *
				{
					state := NewState()
					state.SetBoard(9, 1, &BoardPiece{TurnBlack, UM})
					state.SetBoard(8, 1, &BoardPiece{TurnBlack, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(9, 1), Pos(8, 2), UM},
							expected: "▲8二馬左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(8, 1), Pos(8, 2), UM},
							expected: "▲8二馬右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(9, 5, &BoardPiece{TurnBlack, UM})
					state.SetBoard(6, 3, &BoardPiece{TurnBlack, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(9, 5), Pos(8, 5), UM},
							expected: "▲8五馬寄",
						},
						testData{
							move:     &Move{TurnBlack, Pos(6, 3), Pos(8, 5), UM},
							expected: "▲8五馬引",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(1, 1, &BoardPiece{TurnBlack, UM})
					state.SetBoard(3, 4, &BoardPiece{TurnBlack, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(1, 1), Pos(1, 2), UM},
							expected: "▲1二馬引",
						},
						testData{
							move:     &Move{TurnBlack, Pos(3, 4), Pos(1, 2), UM},
							expected: "▲1二馬上",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(9, 9, &BoardPiece{TurnBlack, UM})
					state.SetBoard(5, 9, &BoardPiece{TurnBlack, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(9, 9), Pos(7, 7), UM},
							expected: "▲7七馬左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(5, 9), Pos(7, 7), UM},
							expected: "▲7七馬右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(4, 7, &BoardPiece{TurnBlack, UM})
					state.SetBoard(1, 8, &BoardPiece{TurnBlack, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnBlack, Pos(4, 7), Pos(2, 9), UM},
							expected: "▲2九馬左",
						},
						testData{
							move:     &Move{TurnBlack, Pos(1, 8), Pos(2, 9), UM},
							expected: "▲2九馬右",
						},
					}
					test(t, state, tests)
				}
			}
			{
				// P1 *  *  *  * -UM *  *  * -UM
				// P2-UM *  *  *  *  *  *  *  *
				// P3 *  *  * -UM *  *  *  *  *
				// P4 *  *  *  *  *  *  *  *  *
				// P5 *  *  *  *  *  *  *  * -UM
				// P6 *  * -UM *  *  *  *  *  *
				// P7 *  *  *  *  * -UM *  *  *
				// P8 *  *  *  *  *  *  *  *  *
				// P9-UM *  *  *  *  *  * -UM-UM
				{
					state := NewState()
					state.SetBoard(1, 9, &BoardPiece{TurnWhite, UM})
					state.SetBoard(2, 9, &BoardPiece{TurnWhite, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(1, 9), Pos(2, 8), UM},
							expected: "△2八馬左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(2, 9), Pos(2, 8), UM},
							expected: "△2八馬右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(1, 5, &BoardPiece{TurnWhite, UM})
					state.SetBoard(4, 7, &BoardPiece{TurnWhite, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(1, 5), Pos(2, 5), UM},
							expected: "△2五馬寄",
						},
						testData{
							move:     &Move{TurnWhite, Pos(4, 7), Pos(2, 5), UM},
							expected: "△2五馬引",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(9, 9, &BoardPiece{TurnWhite, UM})
					state.SetBoard(7, 6, &BoardPiece{TurnWhite, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(9, 9), Pos(9, 8), UM},
							expected: "△9八馬引",
						},
						testData{
							move:     &Move{TurnWhite, Pos(7, 6), Pos(9, 8), UM},
							expected: "△9八馬上",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(1, 1, &BoardPiece{TurnWhite, UM})
					state.SetBoard(5, 1, &BoardPiece{TurnWhite, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(1, 1), Pos(3, 3), UM},
							expected: "△3三馬左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(5, 1), Pos(3, 3), UM},
							expected: "△3三馬右",
						},
					}
					test(t, state, tests)
				}
				{
					state := NewState()
					state.SetBoard(6, 3, &BoardPiece{TurnWhite, UM})
					state.SetBoard(9, 2, &BoardPiece{TurnWhite, UM})
					tests := []testData{
						testData{
							move:     &Move{TurnWhite, Pos(6, 3), Pos(8, 1), UM},
							expected: "△8一馬左",
						},
						testData{
							move:     &Move{TurnWhite, Pos(9, 2), Pos(8, 1), UM},
							expected: "△8一馬右",
						},
					}
					test(t, state, tests)
				}
			}
		}
	}
}

func test(t *testing.T, state *State, tests []testData) {
	for _, test := range tests {
		s := state.Clone()
		result, err := s.MoveString(test.move)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Log(result)
		if result != test.expected {
			t.Errorf("error: expected: %s, actual: %s", test.expected, result)
		}
	}
}
