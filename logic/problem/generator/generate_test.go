package generator

import (
	"bytes"
	"testing"

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
	// 	&data{
	// 		input: `
	// P1 *  *  *  *  *  * +TO * -KI
	// P2 *  *  *  *  * +RY * -OU *
	// P3 *  *  *  *  *  *  * -KI-GI
	// P4 *  *  *  *  *  * -FU *  *
	// P5 *  *  *  *  *  *  *  *  *
	// P6 *  *  *  *  *  *  *  *  *
	// P7 *  *  *  *  *  *  *  *  *
	// P8 *  *  *  *  *  *  *  *  *
	// P9 *  *  *  *  *  *  *  *  *
	// P-00AL
	// `,
	// 		expected: true,
	// 	},
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
	for i, data := range isCheckmateData {
		record, err := csa.Parse(bytes.NewBufferString(data.input))
		if err != nil {
			t.Fatal(err)
		}
		result := isCheckmate(record.State)
		if result != data.expected {
			t.Errorf("error: (result: %v, expected: %v)", result, data.expected)
			continue
		}
		t.Logf("%d: OK", i+1)
	}
}

func TestCheckSolvable(t *testing.T) {
	problemType := Type3
	g := &generator{
		steps: problemType.Steps(),
	}
	for i, data := range checkSolvableTestData {
		record, err := csa.Parse(bytes.NewBufferString(data.input))
		if err != nil {
			t.Fatal(err)
		}
		result := g.isValidProblem(record.State)
		if result != data.expected {
			t.Errorf("error: (result: %v, expected: %v)", result, data.expected)
			continue
		}
		t.Logf("%d: OK", i+1)
	}
}

var checkSolvableTestData = []*data{
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
}
