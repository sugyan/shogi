package generator

import (
	"bytes"
	"testing"

	"github.com/sugyan/shogi/format/csa"
	"github.com/sugyan/shogi/logic/problem/solver"
)

func TestCheckSolvable(t *testing.T) {
	problemType := ProblemType3
	g := &generator{
		steps:  problemType.Steps(),
		solver: solver.NewSolver(3),
	}
	for i, data := range checkSolvableTestData {
		state, err := csa.Parse(bytes.NewBufferString(data.input))
		if err != nil {
			t.Fatal(err)
		}
		result := g.checkSolvable(state)
		if result != data.expected {
			t.Fatalf("error: (result: %v, expected: %v)", result, data.expected)
		}
		t.Logf("%d: OK", i+1)
	}
}

type checkSolvableData struct {
	input    string
	expected bool
}

var checkSolvableTestData = []*checkSolvableData{
	&checkSolvableData{
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
	&checkSolvableData{
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
}
