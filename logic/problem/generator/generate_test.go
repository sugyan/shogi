package generator

import (
	"bytes"
	"testing"
	"time"

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
			steps:   steps,
			timeout: time.Second * 10,
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
