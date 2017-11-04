package solver

import (
	"bytes"
	"testing"

	"github.com/sugyan/shogi/format/csa"
)

func TestSolvers(t *testing.T) {
	for i, data := range testData {
		t.Logf("Q%d...", i+1)
		state, err := csa.Parse(bytes.NewBufferString(data.q))
		if err != nil {
			t.Fatal(err)
		}
		answer, err := Solve(state)
		if err != nil {
			t.Error(err)
			continue
		}
		if len(answer) != len(data.a) {
			t.Errorf("answer length mismatch: %d (expected: %d)", len(answer), len(data.a))
			continue
		}
		for j, move := range answer {
			if move != data.a[j] {
				t.Errorf("error Q%d - A%d: %s != %s", i+1, j+1, move, data.a[j])
			}
		}
	}
}

type data struct {
	q string
	a []string
}

var testData = []*data{
	&data{
		q: `
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  *  *  * +HI *  * 
P3 *  *  *  *  *  * -KE-OU-GI
P4 *  *  *  *  *  * +KE * -FU
P5 *  *  *  *  * +KA * +FU * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
		a: []string{"▲2二桂成"},
	},
	&data{
		q: `
P1 *  *  *  * +TO * -OU-KE * 
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  * +HI *  *  *  * 
P4 *  *  * +KA *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
		a: []string{"▲3三飛成"},
	},
	&data{
		q: `
P1 *  *  *  * -HI *  *  *  * 
P2 *  *  *  *  * -OU+GI *  * 
P3 *  *  *  * -KI *  *  *  * 
P4 *  *  *  *  * +GI *  *  * 
P5 *  * +KA *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P-00AL
`,
		a: []string{"▲4三銀成"},
	},
	&data{
		q: `
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
		a: []string{"▲3二銀", "△1二玉", "▲2三と"},
	},
	&data{
		q: `
P1 *  *  *  *  *  *  * -OU-KY
P2 *  *  *  *  *  *  * -KI * 
P3 *  *  *  *  *  *  * +TO * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * 
P+00KI00KE
P-00AL
`,
		a: []string{"▲3三桂", "△同金", "▲2二金"},
	},
}
