package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/format/csa"
)

type data struct {
	input    string
	expected bool
}

func TestRandom(t *testing.T) {
	for n := 0; n < 100; n++ {
		state := random()
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				file, rank := 9-j, i+1
				b := state.GetBoard(file, rank)
				if b == nil {
					continue
				}
				switch b.Turn {
				case shogi.TurnBlack:
					if rank <= 1 && (b.Piece == shogi.FU || b.Piece == shogi.KY) {
						t.Errorf("invalid BoardPiece %v in (%d, %d)", b.Piece, file, rank)
						break
					}
					if rank <= 2 && (b.Piece == shogi.KE) {
						t.Errorf("invalid BoardPiece %v in (%d, %d)", b.Piece, file, rank)
						break
					}
				case shogi.TurnWhite:
					if rank >= 9 && (b.Piece == shogi.FU || b.Piece == shogi.KY) {
						t.Errorf("invalid BoardPiece %v in (%d, %d)", b.Piece, file, rank)
						break
					}
					if rank >= 8 && (b.Piece == shogi.KE) {
						t.Errorf("invalid BoardPiece %v in (%d, %d)", b.Piece, file, rank)
						break
					}
				}
			}
		}
	}
}

func TestIsCheckmate(t *testing.T) {
	var isCheckmateData = []*data{
		&data{
			input: `
P1 *  *  *  *  *  * +TO * -KI
P2 *  *  *  *  *  *  * -OU * 
P3 *  *  *  * +RY *  * +UM-GI
P4 *  *  *  *  *  * -FU-KI * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
			expected: false,
		},
		&data{
			input: `
P1 *  *  *  *  *  * +TO * -KI
P2 *  *  *  *  *  *  * -OU * 
P3 *  *  *  * +RY *  * -KI-GI
P4 *  *  *  *  *  * -FU *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
			expected: false,
		},
		&data{
			input: `
P1 *  *  *  *  *  * +TO * -KI
P2 *  *  *  *  * +RY * -OU * 
P3 *  *  *  *  *  *  * -KI-GI
P4 *  *  *  *  *  * -FU *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
	`,
			expected: true,
		},
		&data{
			input: `
P1 *  *  *  *  *  * +TO * -KI
P2 *  * +RY *  *  *  * -OU * 
P3 *  *  *  * +KA *  * -KI-GI
P4 *  *  *  *  *  * +FU *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
	`,
			expected: true,
		},
		&data{
			input: `
P1 *  *  *  *  *  *  *  * -KI
P2 *  *  *  *  * +RY * -OU * 
P3 *  *  *  *  *  *  * -KI-GI
P4 *  *  *  *  *  * -FU *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
			expected: false,
		},
	}

	g := &generator{}
	for i, data := range isCheckmateData {
		record, err := csa.Parse(bytes.NewBufferString(data.input))
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		result := g.isCheckmate(record.State)
		elapsed := time.Since(start)
		if result != data.expected {
			t.Errorf("error: (result: %v, expected: %v)", result, data.expected)
			continue
		}
		t.Logf("%d: OK (elapsed time: %v)", i+1, elapsed)
	}
}

func TestIsValidProblem(t *testing.T) {
	g := &generator{}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "testdata", "isValidProblem", "*.csa"))
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range matches {
		basename := filepath.Base(filename)
		s := strings.Split(basename, "_")
		steps, err := strconv.Atoi(s[0])
		if err != nil {
			t.Fatal(err)
		}
		expected := false
		if strings.HasPrefix(s[2], "true") {
			expected = true
		}

		file, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		record, err := csa.Parse(file)
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		result := g.isValidProblem(record.State, steps)
		elapsed := time.Since(start)

		if result != expected {
			t.Errorf("error: %s (result: %v, expected: %v)", basename, result, expected)
			continue
		}
		t.Logf("%15s: OK (elapsed time: %v)", basename, elapsed)
	}
}

func TestCandidatePrevStateS(t *testing.T) {
	{
		record, err := csa.Parse(bytes.NewBufferString(`
P1 *  *  *  *  *  * +RY *  * 
P2 *  *  *  *  *  * -KY * -OU
P3 *  *  *  *  *  *  *  *  * 
P4 *  *  *  *  *  * +KI+GI * 
P5 *  *  *  *  *  *  * -FU * 
P6 *  *  *  *  *  *  *  * +FU
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`))
		if err != nil {
			t.Fatal(err)
		}
		candidates := candidatePrevStatesS(
			record.State,
			&posPiece{
				pos:   shogi.Pos(1, 2),
				piece: shogi.OU,
			},
			shogi.Pos(1, 2),
		)
		for _, state := range candidates {
			ok := true
			for i := 0; i < 9; i++ {
				for j := 0; j < 9; j++ {
					file, rank := 9-j, i+1
					b := state.GetBoard(file, rank)
					if b != nil {
						switch b.Turn {
						case shogi.TurnBlack:
							switch b.Piece {
							case shogi.FU, shogi.KY:
								if rank <= 1 {
									ok = false
								}
							case shogi.KE:
								if rank <= 2 {
									ok = false
								}
							}
						case shogi.TurnWhite:
							switch b.Piece {
							case shogi.FU, shogi.KY:
								if rank >= 9 {
									ok = false
								}
							case shogi.KE:
								if rank >= 8 {
									ok = false
								}
							}
						}
					}
				}
			}
			if !ok {
				t.Errorf("invalid state:\n%v", csa.InitialState1(state))
			}
		}
	}
	{
		record, err := csa.Parse(bytes.NewBufferString(`
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  *  *  * -OU *  * 
P4 *  *  *  *  * +KE+FU-GI+HI
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`))
		if err != nil {
			t.Fatal(err)
		}
		candidates := candidatePrevStatesS(
			record.State,
			&posPiece{
				pos:   shogi.Pos(3, 3),
				piece: shogi.OU,
			},
			shogi.Pos(3, 3),
		)
		for _, state := range candidates {
			ok := true
			for i := 0; i < 9; i++ {
				fu := map[shogi.Turn]int{}
				file := 9 - i
				for j := 0; j < 9; j++ {
					rank := j + 1
					b := state.GetBoard(file, rank)
					if b != nil && b.Piece == shogi.FU {
						fu[b.Turn]++
					}
				}
				if fu[shogi.TurnBlack] > 1 || fu[shogi.TurnWhite] > 1 {
					ok = false
					break
				}
			}
			if !ok {
				t.Errorf("invalid state:\n%v", csa.InitialState1(state))
			}
		}
	}
}
