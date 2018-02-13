package generator

import (
	"bytes"
	"testing"
	"time"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/format/csa"
)

type data struct {
	input    string
	expected bool
}

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

func TestIsCheckmate(t *testing.T) {
	g := &generator{
		timeout: time.Second,
	}
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

var isValidProblemTestData = map[int][]*data{
	1: []*data{
		&data{
			input: `
P1 *  *  *  *  *  *  * -OU-KE
P2 *  *  *  *  * +KI *  *  * 
P3 *  *  *  *  *  *  * +KI * 
P4 *  *  *  *  *  *  *  *  * 
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
P1 *  *  *  *  *  * +KI * -OU
P2 *  *  *  *  *  *  * +FU-KE
P3 *  *  *  *  *  *  *  *  * 
P4 *  *  *  *  *  *  *  *  * 
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
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  * -FU-FU-GI *  * 
P3 *  *  * +GI * -OU-HI *  * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  * +HI *  * +KI * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00KY
P-00AL
`,
			expected: false,
		},
	},
	3: []*data{
		&data{
			input: `
P1 *  *  *  *  *  *  * -OU-KY
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  *  *  * +TO * -FU
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00GI
P-00AL
`,
			expected: true,
		},
		&data{
			input: `
P1 *  *  *  *  *  *  * -GI-OU
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  *  *  *  * +KI * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  * +KA *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00KI
P-00AL
`,
			expected: false,
		},
		&data{
			input: `
P1 *  *  *  * -KE-OU *  * +UM
P2 *  *  *  * -KY-FU *  *  * 
P3 *  *  *  *  * -KY-KI *  * 
P4 *  *  *  *  *  *  * +HI * 
P5 *  *  *  *  * +KE *  *  * 
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
P1 *  * +GI *  * -OU * -HI * 
P2 *  * +HI * -FU *  *  *  * 
P3 *  *  *  *  * +KI+FU *  * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00KE
P-00AL
`,
			expected: false,
		},
		&data{
			input: `
P1 *  *  *  *  * +KI *  *  * 
P2 *  *  *  *  *  * -OU *  * 
P3 *  *  *  * +KA * -KY *  * 
P4 *  *  *  *  *  *  *  * +HI
P5 *  *  *  *  *  *  *  * +KE
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
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  * +HI *  * +FU-OU
P3 *  *  *  *  *  *  *  * -HI
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  * +KY * 
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
P1 *  *  *  *  *  *  *  * -OU
P2 *  *  *  *  *  * +HI *  * 
P3 *  *  *  * -KA *  *  *  * 
P4 *  *  *  *  *  *  *  * +KI
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  * +RY
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
			expected: false,
		},
		&data{
			input: `
P1 *  *  *  *  *  *  *  * -OU
P2 *  *  *  *  *  * +HI *  * 
P3 *  *  *  *  * -KA *  *  * 
P4 *  *  *  *  *  *  *  * +KI
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  * +RY
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
			expected: true,
		},
	},
}

func TestIsValidProblem(t *testing.T) {
	for steps, cases := range isValidProblemTestData {
		t.Logf("--- %d steps:", steps)
		g := &generator{
			steps: steps,
		}
		for i, data := range cases {
			record, err := csa.Parse(bytes.NewBufferString(data.input))
			if err != nil {
				t.Fatal(err)
			}
			start := time.Now()
			result := g.isValidProblem(record.State)
			elapsed := time.Since(start)
			if result != data.expected {
				t.Errorf("error: (result: %v, expected: %v)", result, data.expected)
				continue
			}
			t.Logf("%d: OK (elapsed time: %v)", i+1, elapsed)
		}
	}
}

func TestCandidatePrevStateS(t *testing.T) {
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
