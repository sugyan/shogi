package shogi_test

import (
	"testing"

	"github.com/sugyan/shogi"
)

func TestPieceString(t *testing.T) {
	testCases := []struct {
		piece    shogi.Piece
		expected string
	}{
		{shogi.EMP, " * "},
		{shogi.BFU, "+FU"},
		{shogi.BKY, "+KY"},
		{shogi.BKE, "+KE"},
		{shogi.BGI, "+GI"},
		{shogi.BKI, "+KI"},
		{shogi.BKA, "+KA"},
		{shogi.BHI, "+HI"},
		{shogi.BOU, "+OU"},
		{shogi.BTO, "+TO"},
		{shogi.BNY, "+NY"},
		{shogi.BNK, "+NK"},
		{shogi.BNG, "+NG"},
		{shogi.BUM, "+UM"},
		{shogi.BRY, "+RY"},
		{shogi.WFU, "-FU"},
		{shogi.WKY, "-KY"},
		{shogi.WKE, "-KE"},
		{shogi.WGI, "-GI"},
		{shogi.WKI, "-KI"},
		{shogi.WKA, "-KA"},
		{shogi.WHI, "-HI"},
		{shogi.WOU, "-OU"},
		{shogi.WTO, "-TO"},
		{shogi.WNY, "-NY"},
		{shogi.WNK, "-NK"},
		{shogi.WNG, "-NG"},
		{shogi.WUM, "-UM"},
		{shogi.WRY, "-RY"},
	}
	for _, testCase := range testCases {
		if testCase.piece.String() != testCase.expected {
			t.Errorf("string got '%s', expected: '%s'", testCase.piece.String(), testCase.expected)
		}
	}
}

func TestPiecePromoted(t *testing.T) {
	testCases := []struct {
		piece    shogi.Piece
		expected bool
	}{
		{shogi.EMP, false},
		{shogi.BFU, false},
		{shogi.BKY, false},
		{shogi.BKE, false},
		{shogi.BGI, false},
		{shogi.BKI, false},
		{shogi.BKA, false},
		{shogi.BHI, false},
		{shogi.BOU, false},
		{shogi.BTO, true},
		{shogi.BNY, true},
		{shogi.BNK, true},
		{shogi.BNG, true},
		{shogi.BUM, true},
		{shogi.BRY, true},
		{shogi.WFU, false},
		{shogi.WKY, false},
		{shogi.WKE, false},
		{shogi.WGI, false},
		{shogi.WKI, false},
		{shogi.WKA, false},
		{shogi.WHI, false},
		{shogi.WOU, false},
		{shogi.WTO, true},
		{shogi.WNY, true},
		{shogi.WNK, true},
		{shogi.WNG, true},
		{shogi.WUM, true},
		{shogi.WRY, true},
	}
	for i, testCase := range testCases {
		if testCase.piece.IsPromoted() != testCase.expected {
			t.Errorf("#%d: piece promoted of %v, expected: %v", i, testCase.piece, testCase.expected)
		}
	}
}
