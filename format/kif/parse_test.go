package kif

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sugyan/shogi"
)

func TestParse(t *testing.T) {
	type parsed struct {
		moves []*shogi.Move
	}
	expectedMap := map[string]parsed{
		"everyday_170401.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 1), Dst: shogi.Pos(2, 2), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(2, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(4, 1), Dst: shogi.Pos(3, 1), Piece: shogi.UM},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(3, 1), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(5, 2), Dst: shogi.Pos(5, 1), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(3, 1), Dst: shogi.Pos(2, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(5, 1), Dst: shogi.Pos(1, 1), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(1, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(1, 2), Dst: shogi.Pos(2, 3), Piece: shogi.NG},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(2, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 4), Piece: shogi.KI},
			},
		},
		"everyday_170501.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 2), Piece: shogi.GI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 1), Dst: shogi.Pos(1, 2), Piece: shogi.KY},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 1), Piece: shogi.HI},
			},
		},
		"everyday_170601.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(5, 3), Dst: shogi.Pos(3, 1), Piece: shogi.UM},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(3, 1), Piece: shogi.GI},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 2), Piece: shogi.KI},
			},
		},
		"everyday_170701.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(3, 1), Piece: shogi.KA},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(1, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 4), Piece: shogi.KE},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(3, 3), Dst: shogi.Pos(2, 4), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 2), Dst: shogi.Pos(2, 3), Piece: shogi.NG},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 4), Dst: shogi.Pos(2, 3), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 4), Piece: shogi.KE},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(2, 4), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(5, 3), Dst: shogi.Pos(5, 2), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 2), Dst: shogi.Pos(2, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(5, 2), Dst: shogi.Pos(2, 2), Piece: shogi.RY},
			},
		},
		"everyday_170801.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 3), Dst: shogi.Pos(1, 4), Piece: shogi.NG},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(2, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 4), Dst: shogi.Pos(1, 2), Piece: shogi.UM},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(1, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 3), Piece: shogi.GI},
			},
		},
		"everyday_170901.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 4), Dst: shogi.Pos(2, 3), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(2, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 3), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(1, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 1), Dst: shogi.Pos(1, 2), Piece: shogi.UM},
			},
		},
		"everyday_171001.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 4), Dst: shogi.Pos(1, 2), Piece: shogi.NK},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(1, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 3), Piece: shogi.GI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 2), Dst: shogi.Pos(1, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 5), Piece: shogi.KE},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(1, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 4), Dst: shogi.Pos(3, 2), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(3, 2), Piece: shogi.KA},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 3), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 2), Dst: shogi.Pos(2, 1), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 5), Dst: shogi.Pos(3, 3), Piece: shogi.KE},
			},
		},
		"everyday_171101.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(2, 6), Dst: shogi.Pos(3, 4), Piece: shogi.KE},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(1, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(1, 2), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 3), Dst: shogi.Pos(1, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 1), Dst: shogi.Pos(1, 1), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 2), Dst: shogi.Pos(1, 1), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 2), Piece: shogi.KI},
			},
		},
		"everyday_171201.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(3, 1), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 1), Dst: shogi.Pos(3, 1), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(4, 4), Dst: shogi.Pos(5, 3), Piece: shogi.UM},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(5, 3), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(3, 2), Piece: shogi.KI},
			},
		},
		"everyday_180101.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 4), Dst: shogi.Pos(2, 2), Piece: shogi.NK},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(2, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 3), Piece: shogi.GI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(2, 3), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(3, 5), Piece: shogi.KE},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 3), Dst: shogi.Pos(2, 2), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(4, 4), Dst: shogi.Pos(4, 2), Piece: shogi.RY},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(1, 5), Dst: shogi.Pos(4, 2), Piece: shogi.KA},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 3), Piece: shogi.KI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 2), Dst: shogi.Pos(3, 1), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 5), Dst: shogi.Pos(4, 3), Piece: shogi.KE},
			},
		},
		"everyday_180201.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 3), Piece: shogi.HI},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(2, 4), Dst: shogi.Pos(2, 3), Piece: shogi.HI},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(3, 1), Dst: shogi.Pos(4, 2), Piece: shogi.UM},
			},
		},
		"everyday_180301.kif": parsed{
			moves: []*shogi.Move{
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(2, 4), Piece: shogi.KA},
				&shogi.Move{Turn: shogi.TurnWhite, Src: shogi.Pos(3, 3), Dst: shogi.Pos(2, 4), Piece: shogi.OU},
				&shogi.Move{Turn: shogi.TurnBlack, Src: shogi.Pos(0, 0), Dst: shogi.Pos(3, 4), Piece: shogi.KI},
			},
		},
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "testdata", "*.kif"))
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range matches {
		basename := filepath.Base(filename)
		expected, exist := expectedMap[basename]
		if !exist {
			continue
		}
		file, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		record, err := Parse(file)
		if err != nil {
			t.Fatal(err)
		}
		if len(record.Moves) != len(expected.moves) {
			t.Errorf("error length of moves: %d != %d", len(record.Moves), len(expected.moves))
			continue
		}
		ok := true
		for i, move := range record.Moves {
			if *expected.moves[i] != *move {
				t.Errorf("error move[%d]: %v != %v", i, move, expected.moves[i])
				ok = false
				break
			}
		}
		if ok {
			t.Logf("%15s: OK", basename)
		}
	}
}
