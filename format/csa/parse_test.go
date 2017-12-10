package csa

import "testing"
import "bytes"
import "github.com/sugyan/shogi"

func TestParse(t *testing.T) {
	data := `'----------棋譜ファイルの例"example.csa"-----------------
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
'---------------------------------------------------------`

	record, err := Parse(bytes.NewBufferString(data))
	if err != nil {
		t.Fatal(err)
	}
	{
		b := record.State.GetBoard(5, 1)
		if b != nil && b.Turn == shogi.TurnWhite && b.Piece == shogi.OU {
		} else {
			t.Error("5, 1 is not OU")
		}
	}
	{
		b := record.State.GetBoard(5, 9)
		if b != nil && b.Turn == shogi.TurnBlack && b.Piece == shogi.OU {
		} else {
			t.Error("5, 9 is not OU")
		}
	}
	if len(record.Moves) != 2 {
		t.Error("moves count is not 2")
	}
}
