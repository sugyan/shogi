package logic_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/format/csa"
	"github.com/sugyan/shogi/logic"
)

func TestLegalMovesFromInitialState(t *testing.T) {
	s := logic.NewInitialState()
	results := s.LegalMoves()
	if len(results) != 30 {
		t.Errorf("error!")
	}
}

func TestLegalMovesInRecord(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "testdata", "*.csa"))
	if err != nil {
		t.Fatal(err)
	}
	for i, match := range matches {
		file, err := os.Open(match)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		record, err := csa.Parse(file)
		if err != nil {
			t.Fatal(err)
		}
		s := record.State
		for j, move := range record.Moves {
			ok := false
			for _, legal := range s.LegalMoves() {
				if *move == *legal {
					ok = true
				}
			}
			if !ok {
				t.Errorf("#%d-%d: move %v is not in legal moves", i, j, move)
			}
			s.Move(move)
		}
	}
}

func TestLegalMovesImpossibleMoves(t *testing.T) {
	// TurnBlack
	{
		s := logic.NewState(
			[9][9]shogi.Piece{
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
			[2]shogi.Captured{
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
				{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
			},
			shogi.TurnBlack,
		)
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
		s := logic.NewState(
			[9][9]shogi.Piece{
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
			[2]shogi.Captured{
				{FU: 0, KY: 0, KE: 0, GI: 0, KI: 0, KA: 0, HI: 0},
				{FU: 1, KY: 1, KE: 1, GI: 1, KI: 1, KA: 1, HI: 1},
			},
			shogi.TurnWhite,
		)
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
