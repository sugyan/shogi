package logic

import (
	"math/rand"
	"sort"

	"github.com/sugyan/shogi"
)

type zobrist struct {
	board    map[shogi.Piece][9][9]uint64
	captured map[shogi.RawPiece]map[shogi.Turn]uint64
	turn     map[shogi.Turn]uint64
}

var hasher *zobrist

func init() {
	r := rand.New(rand.NewSource(0)) // TODO
	hasher = &zobrist{
		board:    map[shogi.Piece][9][9]uint64{},
		captured: map[shogi.RawPiece]map[shogi.Turn]uint64{},
		turn:     map[shogi.Turn]uint64{},
	}
	// board
	pieces := make([]shogi.Piece, 0, 28)
	for piece := range shogi.PieceStringMap {
		if piece != shogi.EMP {
			pieces = append(pieces, piece)
		}
	}
	sort.Slice(pieces, func(i, j int) bool { return pieces[i] < pieces[j] })
	for _, piece := range pieces {
		board := [9][9]uint64{}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				board[i][j] = r.Uint64()
			}
		}
		hasher.board[piece] = board
	}
	// captured
	for _, raw := range []shogi.RawPiece{shogi.FU, shogi.KY, shogi.KE, shogi.GI, shogi.KI, shogi.KA, shogi.HI} {
		hasher.captured[raw] = map[shogi.Turn]uint64{
			shogi.TurnBlack: r.Uint64(),
			shogi.TurnWhite: r.Uint64(),
		}
	}
	// turn
	for _, turn := range []shogi.Turn{shogi.TurnBlack, shogi.TurnWhite} {
		hasher.turn[turn] = r.Uint64()
	}
}
