package dfpn

import (
	"github.com/sugyan/shogi"
)

type hash struct {
	turn   shogi.Turn
	pn, dn uint32
}

func (h *hash) p() uint32 {
	switch h.turn {
	case shogi.TurnBlack:
		return h.dn
	case shogi.TurnWhite:
		return h.pn
	}
	return 0
}

func (h *hash) d() uint32 {
	switch h.turn {
	case shogi.TurnBlack:
		return h.pn
	case shogi.TurnWhite:
		return h.dn
	}
	return 0
}
