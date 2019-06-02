package csa_test

import (
	"testing"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/format/csa"
	"github.com/sugyan/shogi/logic"
)

func TestParse(t *testing.T) {
	equals := func(a, b *shogi.Record) bool {
		for i := 0; i < 2; i++ {
			if (a.Players[i] != nil && b.Players[i] != nil && *a.Players[i] != *b.Players[i]) ||
				(a.Players[i] == nil && b.Players[i] != nil) ||
				(a.Players[i] != nil && b.Players[i] == nil) {
				return false
			}
		}
		if !a.State.Equals(b.State) {
			return false
		}
		if len(a.Moves) != len(b.Moves) {
			return false
		}
		for i := 0; i < len(a.Moves); i++ {
			if *a.Moves[i] != *b.Moves[i] {
				return false
			}
		}
		return true
	}
	tests := []struct {
		data     string
		expected *shogi.Record
	}{
		{
			`
'----------棋譜ファイルの例"example.csa"-----------------
'バージョン
V2.2
'対局者名
N+NAKAHARA
N-YONENAGA
'棋譜情報
'棋戦名
$EVENT:13th World Computer Shogi Championship
'対局場所
$SITE:KAZUSA ARC
'開始日時
$START_TIME:2003/05/03 10:30:00
'終了日時
$END_TIME:2003/05/03 11:11:05
'持ち時間:25分、切れ負け
$TIME_LIMIT:00:25+00
'戦型:矢倉
$OPENING:YAGURA
'平手の局面
P1-KY-KE-GI-KI-OU-KI-GI-KE-KY
P2 * -HI *  *  *  *  * -KA * 
P3-FU-FU-FU-FU-FU-FU-FU-FU-FU
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7+FU+FU+FU+FU+FU+FU+FU+FU+FU
P8 * +KA *  *  *  *  * +HI * 
P9+KY+KE+GI+KI+OU+KI+GI+KE+KY
'先手番
+
'指し手と消費時間
+2726FU
T12
-3334FU
T6
%CHUDAN
'---------------------------------------------------------`,
			&shogi.Record{
				Players: [2]*shogi.Player{
					{Name: "NAKAHARA"},
					{Name: "YONENAGA"},
				},
				State: logic.NewState(
					[9][9]shogi.Piece{
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
					[2]shogi.Captured{},
					shogi.TurnBlack,
				),
				Moves: []*shogi.Move{
					{Src: shogi.Position{File: 2, Rank: 7}, Dst: shogi.Position{File: 2, Rank: 6}, Piece: shogi.BFU},
					{Src: shogi.Position{File: 3, Rank: 3}, Dst: shogi.Position{File: 3, Rank: 4}, Piece: shogi.WFU},
				},
			},
		},
		{
			`
PI82HI22KA`,
			&shogi.Record{
				State: logic.NewState(
					[9][9]shogi.Piece{
						{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
						{shogi.EMP, shogi.BKA, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.BHI, shogi.EMP},
						{shogi.BKY, shogi.BKE, shogi.BGI, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
					},
					[2]shogi.Captured{},
					shogi.TurnBlack,
				),
				Moves: []*shogi.Move{},
			},
		},
		{
			`
P-22KA
P+99KY89KE
P+00KI00FU
P-00AL`,
			&shogi.Record{
				State: logic.NewState(
					[9][9]shogi.Piece{
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKA, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
						{shogi.BKY, shogi.BKE, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
					},
					[2]shogi.Captured{
						{FU: 1, KY: 0, KE: 0, GI: 0, KI: 1, KA: 0, HI: 0},
						{FU: 17, KY: 3, KE: 3, GI: 4, KI: 3, KA: 1, HI: 2},
					},
					shogi.TurnBlack,
				),
				Moves: []*shogi.Move{},
			},
		},
	}
	for i, tc := range tests {
		record, err := csa.ParseString(tc.data)
		if err != nil {
			t.Fatal(err)
		}
		if !equals(record, tc.expected) {
			t.Errorf("#%d: got: %v, expected: %v", i, record, tc.expected)
		}
	}
}
