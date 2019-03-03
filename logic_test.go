package shogi_test

import (
	"testing"

	"github.com/sugyan/shogi"
)

func TestLegalMoves(t *testing.T) {
	s := *shogi.InitialState
	results := s.LegalMoves()
	if len(results) != 30 {
		t.Errorf("error!")
	}
}
