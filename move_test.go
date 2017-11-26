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
		state.SetBoardPiece(1, 9, &BoardPiece{TurnFirst, KY})
		state.SetBoardPiece(2, 9, &BoardPiece{TurnFirst, KE})
		state.SetBoardPiece(3, 9, &BoardPiece{TurnFirst, GI})
		state.SetBoardPiece(4, 9, &BoardPiece{TurnFirst, KI})
		state.SetBoardPiece(8, 8, &BoardPiece{TurnFirst, KA})
		state.SetBoardPiece(2, 8, &BoardPiece{TurnFirst, HI})
		state.SetBoardPiece(9, 1, &BoardPiece{TurnSecond, KY})
		state.SetBoardPiece(8, 1, &BoardPiece{TurnSecond, KE})
		state.SetBoardPiece(7, 1, &BoardPiece{TurnSecond, GI})
		state.SetBoardPiece(6, 1, &BoardPiece{TurnSecond, KI})
		state.SetBoardPiece(2, 2, &BoardPiece{TurnSecond, KA})
		state.SetBoardPiece(8, 2, &BoardPiece{TurnSecond, HI})
		state.Captured[TurnFirst] = &CapturedPieces{KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1}
		state.Captured[TurnSecond] = &CapturedPieces{KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1}
		tests := []testData{
			// KY
			testData{
				move:     &Move{TurnFirst, Pos(1, 9), Pos(1, 8), KY},
				expected: "▲1八香",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(1, 8), KY},
				expected: "▲1八香打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(2, 8), KY},
				expected: "▲2八香",
			},
			testData{
				move:     &Move{TurnSecond, Pos(9, 1), Pos(9, 2), KY},
				expected: "△9二香",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(9, 2), KY},
				expected: "△9二香打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(8, 2), KY},
				expected: "△8二香",
			},
			// KE
			testData{
				move:     &Move{TurnFirst, Pos(2, 9), Pos(1, 7), KE},
				expected: "▲1七桂",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(1, 7), KE},
				expected: "▲1七桂打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(2, 7), KE},
				expected: "▲2七桂",
			},
			testData{
				move:     &Move{TurnSecond, Pos(8, 1), Pos(9, 3), KE},
				expected: "△9三桂",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(9, 3), KE},
				expected: "△9三桂打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(8, 3), KE},
				expected: "△8三桂",
			},
			// GI
			testData{
				move:     &Move{TurnFirst, Pos(3, 9), Pos(3, 8), GI},
				expected: "▲3八銀",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(3, 8), GI},
				expected: "▲3八銀打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(3, 7), GI},
				expected: "▲3七銀",
			},
			testData{
				move:     &Move{TurnSecond, Pos(7, 1), Pos(7, 2), GI},
				expected: "△7二銀",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(7, 2), GI},
				expected: "△7二銀打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(7, 3), GI},
				expected: "△7三銀",
			},
			// KI
			testData{
				move:     &Move{TurnFirst, Pos(4, 9), Pos(5, 8), KI},
				expected: "▲5八金",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(5, 8), KI},
				expected: "▲5八金打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(5, 7), KI},
				expected: "▲5七金",
			},
			testData{
				move:     &Move{TurnSecond, Pos(4, 9), Pos(5, 2), KI},
				expected: "△5二金",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(5, 2), KI},
				expected: "△5二金打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(5, 3), KI},
				expected: "△5三金",
			},
			// KA
			testData{
				move:     &Move{TurnFirst, Pos(8, 8), Pos(7, 7), KA},
				expected: "▲7七角",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(7, 7), KA},
				expected: "▲7七角打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(8, 7), KA},
				expected: "▲8七角",
			},
			testData{
				move:     &Move{TurnSecond, Pos(2, 2), Pos(3, 3), KA},
				expected: "△3三角",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(3, 3), KA},
				expected: "△3三角打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(2, 3), KA},
				expected: "△2三角",
			},
			// HI
			testData{
				move:     &Move{TurnFirst, Pos(2, 8), Pos(2, 7), HI},
				expected: "▲2七飛",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(2, 7), HI},
				expected: "▲2七飛打",
			},
			testData{
				move:     &Move{TurnFirst, Pos(0, 0), Pos(3, 7), HI},
				expected: "▲3七飛",
			},
			testData{
				move:     &Move{TurnSecond, Pos(8, 2), Pos(8, 3), HI},
				expected: "△8三飛",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(8, 3), HI},
				expected: "△8三飛打",
			},
			testData{
				move:     &Move{TurnSecond, Pos(0, 0), Pos(7, 3), HI},
				expected: "△7三飛",
			},
		}
		for _, test := range tests {
			s := state.Clone()
			result, err := s.MoveString(test.move)
			if err != nil {
				t.Error(err)
				continue
			}
			if result != test.expected {
				t.Errorf("error: expected: %s, actual: %s", test.expected, result)
			}
		}
	}
	// 成・不成
	{
		state := NewState()
		state.SetBoardPiece(1, 4, &BoardPiece{TurnFirst, FU})
		state.SetBoardPiece(2, 4, &BoardPiece{TurnFirst, KY})
		state.SetBoardPiece(3, 5, &BoardPiece{TurnFirst, KE})
		state.SetBoardPiece(4, 4, &BoardPiece{TurnFirst, GI})
		state.SetBoardPiece(5, 4, &BoardPiece{TurnFirst, KA})
		state.SetBoardPiece(6, 4, &BoardPiece{TurnFirst, HI})
		state.SetBoardPiece(9, 6, &BoardPiece{TurnSecond, FU})
		state.SetBoardPiece(8, 6, &BoardPiece{TurnSecond, KY})
		state.SetBoardPiece(7, 5, &BoardPiece{TurnSecond, KE})
		state.SetBoardPiece(6, 6, &BoardPiece{TurnSecond, GI})
		state.SetBoardPiece(5, 6, &BoardPiece{TurnSecond, KA})
		state.SetBoardPiece(4, 6, &BoardPiece{TurnSecond, HI})
		tests := []testData{
			// FU
			testData{
				move:     &Move{TurnFirst, Pos(1, 4), Pos(1, 3), FU},
				expected: "▲1三歩不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(1, 4), Pos(1, 3), TO},
				expected: "▲1三歩成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(9, 6), Pos(9, 7), FU},
				expected: "△9七歩不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(9, 6), Pos(9, 7), TO},
				expected: "△9七歩成",
			},
			// KY
			testData{
				move:     &Move{TurnFirst, Pos(2, 4), Pos(2, 3), KY},
				expected: "▲2三香不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(2, 4), Pos(2, 3), NY},
				expected: "▲2三香成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(8, 6), Pos(8, 7), KY},
				expected: "△8七香不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(8, 6), Pos(8, 7), NY},
				expected: "△8七香成",
			},
			// KE
			testData{
				move:     &Move{TurnFirst, Pos(3, 5), Pos(2, 3), KE},
				expected: "▲2三桂不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(3, 5), Pos(2, 3), NK},
				expected: "▲2三桂成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(7, 5), Pos(8, 7), KE},
				expected: "△8七桂不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(7, 5), Pos(8, 7), NK},
				expected: "△8七桂成",
			},
			// GI
			testData{
				move:     &Move{TurnFirst, Pos(4, 4), Pos(4, 3), GI},
				expected: "▲4三銀不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(4, 4), Pos(4, 3), NG},
				expected: "▲4三銀成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(6, 6), Pos(6, 7), GI},
				expected: "△6七銀不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(6, 6), Pos(6, 7), NG},
				expected: "△6七銀成",
			},
			// KA
			testData{
				move:     &Move{TurnFirst, Pos(5, 4), Pos(4, 3), KA},
				expected: "▲4三角不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(5, 4), Pos(4, 3), UM},
				expected: "▲4三角成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(5, 6), Pos(6, 7), KA},
				expected: "△6七角不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(5, 6), Pos(6, 7), UM},
				expected: "△6七角成",
			},
			// HI
			testData{
				move:     &Move{TurnFirst, Pos(6, 4), Pos(6, 3), HI},
				expected: "▲6三飛不成",
			},
			testData{
				move:     &Move{TurnFirst, Pos(6, 4), Pos(6, 3), RY},
				expected: "▲6三飛成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(4, 6), Pos(4, 7), HI},
				expected: "△4七飛不成",
			},
			testData{
				move:     &Move{TurnSecond, Pos(4, 6), Pos(4, 7), RY},
				expected: "△4七飛成",
			},
		}
		for _, test := range tests {
			s := state.Clone()
			result, err := s.MoveString(test.move)
			if err != nil {
				t.Error(err)
				continue
			}
			if result != test.expected {
				t.Errorf("error: expected: %s, actual: %s", test.expected, result)
			}
		}
	}
}
