package shogi

// Clone method
func (s *State) Clone() *State {
	state := NewState()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := s.Board[i][j]
			if bp != nil {
				state.Board[i][j] = &BoardPiece{
					Turn:  bp.Turn,
					Piece: NewPiece(bp.Piece.Code()),
				}
			}
		}
	}
	capF, capS := *s.Captured[TurnFirst], *s.Captured[TurnSecond]
	state.Captured = map[Turn]*CapturedPieces{
		TurnFirst:  &capF,
		TurnSecond: &capS,
	}
	return state
}

// Check method
func (s *State) Check(turn Turn) *Move {
	var targetPos *Position
searchTarget:
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := s.Board[i][j]
			if bp != nil && bp.Piece.Code() == OU && bp.Turn != turn {
				targetPos = &Position{9 - j, i + 1}
				break searchTarget
			}
		}
	}
	for _, m := range s.CandidateMoves(turn) {
		if *m.Dst == *targetPos {
			return m
		}
	}
	return nil
}

// CandidateMoves method
func (s *State) CandidateMoves(turn Turn) []*Move {
	results := []*Move{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			bp := s.Board[i][j]
			if bp != nil && bp.Turn == turn {
				src := &Position{9 - j, i + 1}
				for _, dst := range s.movable(bp, src) {
					results = append(results, &Move{
						Src:   src,
						Dst:   dst,
						Piece: bp.Piece,
					})
				}
			}
		}
	}
	return results
}

func (s *State) movable(bp *BoardPiece, src *Position) []*Position {
	positions := []*Position{}
	switch bp.Piece.Code() {
	case FU:
		switch bp.Turn {
		case TurnFirst:
			positions = append(positions, &Position{src.File, src.Rank - 1})
		case TurnSecond:
			positions = append(positions, &Position{src.File, src.Rank + 1})
		}
	case KY:
		switch bp.Turn {
		case TurnFirst:
			for i := 1; src.Rank-i > 0; i++ {
				dst := &Position{src.File, src.Rank - i}
				positions = append(positions, dst)
				if s.GetBoardPiece(dst.File, dst.Rank) != nil {
					break
				}
			}
		case TurnSecond:
			for i := 1; src.Rank+i < 10; i++ {
				dst := &Position{src.File, src.Rank + i}
				positions = append(positions, dst)
				if s.GetBoardPiece(dst.File, dst.Rank) != nil {
					break
				}
			}
		}
	case KE:
		switch bp.Turn {
		case TurnFirst:
			positions = append(positions, &Position{src.File - 1, src.Rank - 2})
			positions = append(positions, &Position{src.File + 1, src.Rank - 2})
		case TurnSecond:
			positions = append(positions, &Position{src.File - 1, src.Rank + 2})
			positions = append(positions, &Position{src.File + 1, src.Rank + 2})
		}
	case GI:
		positions = append(positions, &Position{src.File - 1, src.Rank - 1})
		positions = append(positions, &Position{src.File + 1, src.Rank - 1})
		positions = append(positions, &Position{src.File - 1, src.Rank + 1})
		positions = append(positions, &Position{src.File + 1, src.Rank + 1})
		switch bp.Turn {
		case TurnFirst:
			positions = append(positions, &Position{src.File, src.Rank - 1})
		case TurnSecond:
			positions = append(positions, &Position{src.File, src.Rank + 1})
		}
	case TO:
		fallthrough
	case NY:
		fallthrough
	case NK:
		fallthrough
	case NG:
		fallthrough
	case KI:
		positions = append(positions, &Position{src.File - 1, src.Rank})
		positions = append(positions, &Position{src.File + 1, src.Rank})
		positions = append(positions, &Position{src.File, src.Rank - 1})
		positions = append(positions, &Position{src.File, src.Rank + 1})
		switch bp.Turn {
		case TurnFirst:
			positions = append(positions, &Position{src.File - 1, src.Rank - 1})
			positions = append(positions, &Position{src.File + 1, src.Rank - 1})
		case TurnSecond:
			positions = append(positions, &Position{src.File - 1, src.Rank + 1})
			positions = append(positions, &Position{src.File + 1, src.Rank + 1})
		}
	case UM:
		for _, d := range []*Position{&Position{0, -1}, &Position{0, 1}, &Position{-1, 0}, &Position{1, 0}} {
			positions = append(positions, &Position{src.File + d.File, src.Rank + d.Rank})
		}
		fallthrough
	case KA:
		for i := 1; src.File-i > 0 && src.Rank-i > 0; i++ {
			dst := &Position{src.File - i, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File-i > 0 && src.Rank+i < 10; i++ {
			dst := &Position{src.File - i, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File+i < 10 && src.Rank-i > 0; i++ {
			dst := &Position{src.File + i, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File+i < 10 && src.Rank+i < 10; i++ {
			dst := &Position{src.File + i, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
	case RY:
		for _, d := range []*Position{&Position{1, 1}, &Position{1, -1}, &Position{-1, 1}, &Position{-1, -1}} {
			positions = append(positions, &Position{src.File + d.File, src.Rank + d.Rank})
		}
		fallthrough
	case HI:
		for i := 1; src.File+i < 10; i++ {
			dst := &Position{src.File + i, src.Rank}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File-i > 0; i++ {
			dst := &Position{src.File - i, src.Rank}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.Rank+i < 10; i++ {
			dst := &Position{src.File, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.Rank-i > 0; i++ {
			dst := &Position{src.File, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoardPiece(dst.File, dst.Rank) != nil {
				break
			}
		}
	case OU:
		for _, d := range []*Position{
			&Position{1, 1}, &Position{1, -1}, &Position{-1, 1}, &Position{-1, -1},
			&Position{0, 1}, &Position{0, -1}, &Position{-1, 0}, &Position{1, 0},
		} {
			positions = append(positions, &Position{src.File + d.File, src.Rank + d.Rank})
		}
	}
	// filtering to valid positions
	results := []*Position{}
	for _, pos := range positions {
		if pos.File > 0 && pos.File < 10 && pos.Rank > 0 && pos.Rank < 10 {
			dstBp := s.Board[pos.Rank-1][9-pos.File]
			if dstBp != nil && dstBp.Turn == bp.Turn {
				continue
			}
			results = append(results, pos)
		}
	}
	return results
}
