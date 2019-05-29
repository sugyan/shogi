package shogi

import (
	"math/rand"
	"sort"
)

type zobrist struct {
	board    map[Piece][9][9]uint64
	captured map[rawPiece]map[Turn]uint64
	turn     map[Turn]uint64
}

var hasher *zobrist

func init() {
	r := rand.New(rand.NewSource(0)) // TODO
	hasher = &zobrist{
		board:    map[Piece][9][9]uint64{},
		captured: map[rawPiece]map[Turn]uint64{},
		turn:     map[Turn]uint64{},
	}
	// board
	pieces := make([]Piece, 0, 28)
	for piece := range pieceStringMap {
		if piece != EMP {
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
	for _, raw := range []rawPiece{fu, ky, ke, gi, ki, ka, hi} {
		hasher.captured[raw] = map[Turn]uint64{
			TurnBlack: r.Uint64(),
			TurnWhite: r.Uint64(),
		}
	}
	// turn
	for _, turn := range []Turn{TurnBlack, TurnWhite} {
		hasher.turn[turn] = r.Uint64()
	}

}
