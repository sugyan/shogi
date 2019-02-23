package csa_test

import (
	"reflect"
	"testing"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/format/csa"
)

func TestParse(t *testing.T) {
	tests := []struct {
		data     string
		expected *shogi.Record
	}{
		{
			`'----------棋譜ファイルの例"example.csa"-----------------
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
				State: &shogi.State{
					Board: [9][9]shogi.Piece{
						{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
						{shogi.BLANK, shogi.WHI, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.WKA, shogi.BLANK},
						{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
						{shogi.BLANK, shogi.BKA, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BHI, shogi.BLANK},
						{shogi.BKY, shogi.BKE, shogi.BGI, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
					},
				},
				Moves: []*shogi.Move{
					{Src: shogi.Position{File: 2, Rank: 7}, Dst: shogi.Position{File: 2, Rank: 6}, Piece: shogi.BFU},
					{Src: shogi.Position{File: 3, Rank: 3}, Dst: shogi.Position{File: 3, Rank: 4}, Piece: shogi.WFU},
				},
			},
		},
		{
			`PI82HI22KA`,
			&shogi.Record{
				State: &shogi.State{
					Board: [9][9]shogi.Piece{
						{shogi.WKY, shogi.WKE, shogi.WGI, shogi.WKI, shogi.WOU, shogi.WKI, shogi.WGI, shogi.WKE, shogi.WKY},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU, shogi.WFU},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU, shogi.BFU},
						{shogi.BLANK, shogi.BKA, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BHI, shogi.BLANK},
						{shogi.BKY, shogi.BKE, shogi.BGI, shogi.BKI, shogi.BOU, shogi.BKI, shogi.BGI, shogi.BKE, shogi.BKY},
					},
				},
				Moves: []*shogi.Move{},
			},
		},
		{
			`'
P-22KA
P+99KY89KE
P+00KI00FU
P-00AL`,
			&shogi.Record{
				State: &shogi.State{
					Board: [9][9]shogi.Piece{
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.WKA, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
						{shogi.BKY, shogi.BKE, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK, shogi.BLANK},
					},
					Captured: [2]shogi.Captured{
						{FU: 1, KY: 0, KE: 0, GI: 0, KI: 1, KA: 0, HI: 0},
						{FU: 17, KY: 3, KE: 3, GI: 4, KI: 3, KA: 1, HI: 2},
					},
				},
				Moves: []*shogi.Move{},
			},
		},
	}
	for _, tc := range tests {
		record, err := csa.ParseString(tc.data)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(record, tc.expected) {
			t.Errorf("got: %#v, expected: %#v", record, tc.expected)
		}
	}
}
