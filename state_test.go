package shogi_test

import (
	"testing"

	"github.com/sugyan/shogi"
)

func TestString(t *testing.T) {
	testCases := []struct {
		state    shogi.State
		expected string
	}{
		{
			state: shogi.State{},
			expected: `
P1 *  *  *  *  *  *  *  *  * 
P2 *  *  *  *  *  *  *  *  * 
P3 *  *  *  *  *  *  *  *  * 
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7 *  *  *  *  *  *  *  *  * 
P8 *  *  *  *  *  *  *  *  * 
P9 *  *  *  *  *  *  *  *  * `[1:],
		},
		{
			state: *shogi.InitialState,
			expected: `
P1-KY-KE-GI-KI-OU-KI-GI-KE-KY
P2 * -HI *  *  *  *  * -KA * 
P3-FU-FU-FU-FU-FU-FU-FU-FU-FU
P4 *  *  *  *  *  *  *  *  * 
P5 *  *  *  *  *  *  *  *  * 
P6 *  *  *  *  *  *  *  *  * 
P7+FU+FU+FU+FU+FU+FU+FU+FU+FU
P8 * +KA *  *  *  *  * +HI * 
P9+KY+KE+GI+KI+OU+KI+GI+KE+KY`[1:],
		},
	}
	for i, testCase := range testCases {
		s := testCase.state.String()
		if s != testCase.expected {
			t.Errorf("#%d: got %v, expected: %v", i, s, testCase.expected)
		}
	}
}
