package solver

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sugyan/shogi/format/csa"
)

func TestSolver(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	matches, err := filepath.Glob(filepath.Join(dir, "testdata", "*.csa"))
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range matches {
		file, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}
		record, err := csa.Parse(file)
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		answer := Solve(record.State)
		elapsed := time.Since(start)
		if len(answer) != len(record.Moves) {
			t.Errorf("error answer length: %d (expected: %d)", len(answer), len(record.Moves))
			continue
		}
		ok := true
		for i, move := range answer {
			if len(answer) >= 3 && i >= len(answer)-2 {
				continue
			}
			if *move != *record.Moves[i] {
				t.Errorf("error A[%d]: %v != %v", i+1, *move, *record.Moves[i])
				ok = false
				break
			}
		}
		if ok {
			t.Logf("%15s: OK (elapsed time: %v)", filepath.Base(filename), elapsed)
		}
	}
}
