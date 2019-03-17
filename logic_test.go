package shogi_test

import (
	"github.com/sugyan/shogi/format/csa"
	"testing"

	"github.com/sugyan/shogi"
)

func TestLegalMovesFromInitialState(t *testing.T) {
	s := *shogi.InitialState
	results := s.LegalMoves()
	if len(results) != 30 {
		t.Errorf("error!")
	}
}

func TestLegalMovesInRecord(t *testing.T) {
	record, err := csa.ParseString(`
PI
+
+5756FU
-3334FU
+2858HI
-8232HI
+7776FU
-5162OU
+5948OU
-6272OU
+4838OU
-7282OU
+3828OU
-7172GI
+3938GI
-9394FU
+8822UM
-3122GI
+7988GI
-4152KI
+8877GI
-2233GI
+5655FU
-1314FU
+1716FU
-3222HI
+7766GI
-2324FU
+6665GI
-2425FU
+5554FU
-5354FU
+6554GI
-2526FU
+2726FU
-2226HI
+0027FU
-2676HI
+5453NG
-0057FU
+5878HI
-7678RY
+6978KI
-5253KI
+0031HI
-0012GI
+3132RY
-0026FU
+8977KE
-2627TO
+3827GI
-0069HI
+4939KI
-5758TO
+7765KE
-5352KI
+0075KA
-0064KA
+7566KA
-3435FU
+0054FU
-0026FU
+2726GI
-3536FU
+5453TO
-5251KI
+5363TO
-3637TO
+2937KE
-0036FU
+6364TO
-3637TO
+2637GI
-0035KE
+6473TO
-8173KE
+6573NK
-8273OU
+3243RY
-0063FU
+0046KA
-0064KE
+4635KA
-0027FU
+2838OU
-3342GI
+0085KE
-7374OU
+4342RY
-5142KI
+0075GI
-7485OU
+8786FU
-8576OU
+7877KI
-7665OU
+6611UM
-6939RY
+3827OU
-0025HI
+2736OU
-2535HI
+3635OU
-3937RY
+0036FU
-0043KE
+3525OU
-3727RY
+0026FU
-0024GI
+2524OU
-2726RY
+0025FU
-0023KI`[1:])
	if err != nil {
		t.Fatal(err)
	}
	s := record.State
	for i, move := range record.Moves {
		ok := false
		for _, legal := range s.LegalMoves() {
			if *move == *legal {
				ok = true
			}
		}
		if !ok {
			t.Errorf("#%d: move %v is not in legal moves", i, move)
		}
		s.Move(move)
	}
}

func TestLegalMovesImpossibleMoves(t *testing.T) {
	// TurnBlack
	{
		s := &shogi.State{
			Board: [9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.BFU, shogi.BKY, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKE, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.BKE, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			Captured: [2]shogi.Captured{
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
				{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
			},
			Turn: shogi.TurnBlack,
		}
		impossibleMoves := []*shogi.Move{
			{Src: shogi.Position{File: 9, Rank: 2}, Dst: shogi.Position{File: 9, Rank: 1}, Piece: shogi.BFU},
			{Src: shogi.Position{File: 8, Rank: 2}, Dst: shogi.Position{File: 8, Rank: 1}, Piece: shogi.BKY},
			{Src: shogi.Position{File: 6, Rank: 3}, Dst: shogi.Position{File: 7, Rank: 1}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 6, Rank: 3}, Dst: shogi.Position{File: 5, Rank: 1}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 6, Rank: 4}, Dst: shogi.Position{File: 7, Rank: 2}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 6, Rank: 4}, Dst: shogi.Position{File: 5, Rank: 2}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 9, Rank: 1}, Piece: shogi.BFU},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 8, Rank: 1}, Piece: shogi.BFU},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 7, Rank: 1}, Piece: shogi.BKY},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 6, Rank: 1}, Piece: shogi.BKY},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 5, Rank: 1}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 4, Rank: 1}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 3, Rank: 2}, Piece: shogi.BKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 2, Rank: 2}, Piece: shogi.BKE},
		}
		for _, move := range s.LegalMoves() {
			for _, m := range impossibleMoves {
				if *m == *move {
					t.Errorf("move %v should be impossible", move)
				}
			}
		}
	}
	// TurnWhite
	{
		s := &shogi.State{
			Board: [9][9]shogi.Piece{
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKE, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKE, shogi.EMP, shogi.EMP, shogi.EMP},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.WKY, shogi.WFU},
				{shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP, shogi.EMP},
			},
			Captured: [2]shogi.Captured{
				{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
			},
			Turn: shogi.TurnWhite,
		}
		impossibleMoves := []*shogi.Move{
			{Src: shogi.Position{File: 1, Rank: 8}, Dst: shogi.Position{File: 1, Rank: 9}, Piece: shogi.WFU},
			{Src: shogi.Position{File: 2, Rank: 8}, Dst: shogi.Position{File: 2, Rank: 9}, Piece: shogi.WKY},
			{Src: shogi.Position{File: 4, Rank: 7}, Dst: shogi.Position{File: 3, Rank: 9}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 4, Rank: 7}, Dst: shogi.Position{File: 5, Rank: 9}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 4, Rank: 6}, Dst: shogi.Position{File: 3, Rank: 8}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 4, Rank: 6}, Dst: shogi.Position{File: 5, Rank: 8}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 1, Rank: 9}, Piece: shogi.WFU},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 2, Rank: 9}, Piece: shogi.WFU},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 3, Rank: 9}, Piece: shogi.WKY},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 4, Rank: 9}, Piece: shogi.WKY},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 5, Rank: 9}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 6, Rank: 9}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 7, Rank: 8}, Piece: shogi.WKE},
			{Src: shogi.Position{File: 0, Rank: 0}, Dst: shogi.Position{File: 8, Rank: 8}, Piece: shogi.WKE},
		}
		for _, move := range s.LegalMoves() {
			for _, m := range impossibleMoves {
				if *m == *move {
					t.Errorf("move %v should be impossible", move)
				}
			}
		}
	}
}
