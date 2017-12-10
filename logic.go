package shogi

// Clone method
func (s *State) Clone() *State {
	state := NewState()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := s.Board[i][j]
			if b != nil {
				state.Board[i][j] = &BoardPiece{
					Turn:  b.Turn,
					Piece: b.Piece,
				}
			}
		}
	}
	capB, capW := *s.Captured[TurnBlack], *s.Captured[TurnWhite]
	state.Captured = map[Turn]*CapturedPieces{
		TurnBlack: &capB,
		TurnWhite: &capW,
	}
	return state
}

// Check method
func (s *State) Check(turn Turn) *Move {
	var targetPos Position
searchTarget:
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := s.Board[i][j]
			if b != nil && b.Piece == OU && b.Turn != turn {
				targetPos = Position{9 - j, i + 1}
				break searchTarget
			}
		}
	}
	for _, move := range s.CandidateMoves(turn) {
		if move.Dst == targetPos {
			return move
		}
	}
	return nil
}

// CandidateMoves method
func (s *State) CandidateMoves(turn Turn) []*Move {
	results := []*Move{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			b := s.Board[i][j]
			if b != nil && b.Turn == turn {
				src := Pos(9-j, i+1)
				for _, dst := range s.movable(b, src) {
					mustPromote := false
					if turn == TurnBlack && (src.Rank <= 3 || dst.Rank <= 3) {
						switch b.Piece {
						case FU, KY:
							if dst.Rank <= 1 {
								mustPromote = true
							}
						case KE:
							if dst.Rank <= 2 {
								mustPromote = true
							}
						}
						if promoted, ok := promoteMap[b.Piece]; ok {
							results = append(results, &Move{
								Turn:  turn,
								Src:   src,
								Dst:   dst,
								Piece: promoted,
							})
						}
					}
					if turn == TurnWhite && (src.Rank >= 7 || dst.Rank >= 7) {
						switch b.Piece {
						case FU, KY:
							if dst.Rank >= 9 {
								mustPromote = true
							}
						case KE:
							if dst.Rank >= 8 {
								mustPromote = true
							}
						}
						if promoted, ok := promoteMap[b.Piece]; ok {
							results = append(results, &Move{
								Turn:  turn,
								Src:   src,
								Dst:   dst,
								Piece: promoted,
							})
						}
					}
					if !mustPromote {
						results = append(results, &Move{
							Turn:  turn,
							Src:   src,
							Dst:   dst,
							Piece: b.Piece,
						})
					}
				}
			}
		}
	}
	return results
}

func (s *State) movable(b *BoardPiece, src Position) []Position {
	positions := []Position{}
	switch b.Piece {
	case FU:
		switch b.Turn {
		case TurnBlack:
			positions = append(positions, Position{src.File, src.Rank - 1})
		case TurnWhite:
			positions = append(positions, Position{src.File, src.Rank + 1})
		}
	case KY:
		switch b.Turn {
		case TurnBlack:
			for i := 1; src.Rank-i > 0; i++ {
				dst := Position{src.File, src.Rank - i}
				positions = append(positions, dst)
				if s.GetBoard(dst.File, dst.Rank) != nil {
					break
				}
			}
		case TurnWhite:
			for i := 1; src.Rank+i < 10; i++ {
				dst := Position{src.File, src.Rank + i}
				positions = append(positions, dst)
				if s.GetBoard(dst.File, dst.Rank) != nil {
					break
				}
			}
		}
	case KE:
		switch b.Turn {
		case TurnBlack:
			positions = append(positions, Position{src.File - 1, src.Rank - 2})
			positions = append(positions, Position{src.File + 1, src.Rank - 2})
		case TurnWhite:
			positions = append(positions, Position{src.File - 1, src.Rank + 2})
			positions = append(positions, Position{src.File + 1, src.Rank + 2})
		}
	case GI:
		positions = append(positions, Position{src.File - 1, src.Rank - 1})
		positions = append(positions, Position{src.File + 1, src.Rank - 1})
		positions = append(positions, Position{src.File - 1, src.Rank + 1})
		positions = append(positions, Position{src.File + 1, src.Rank + 1})
		switch b.Turn {
		case TurnBlack:
			positions = append(positions, Position{src.File, src.Rank - 1})
		case TurnWhite:
			positions = append(positions, Position{src.File, src.Rank + 1})
		}
	case TO, NY, NK, NG, KI:
		positions = append(positions, Position{src.File - 1, src.Rank})
		positions = append(positions, Position{src.File + 1, src.Rank})
		positions = append(positions, Position{src.File, src.Rank - 1})
		positions = append(positions, Position{src.File, src.Rank + 1})
		switch b.Turn {
		case TurnBlack:
			positions = append(positions, Position{src.File - 1, src.Rank - 1})
			positions = append(positions, Position{src.File + 1, src.Rank - 1})
		case TurnWhite:
			positions = append(positions, Position{src.File - 1, src.Rank + 1})
			positions = append(positions, Position{src.File + 1, src.Rank + 1})
		}
	case UM:
		for _, d := range []Position{Position{0, -1}, Position{0, 1}, Position{-1, 0}, Position{1, 0}} {
			positions = append(positions, Position{src.File + d.File, src.Rank + d.Rank})
		}
		fallthrough
	case KA:
		for i := 1; src.File-i > 0 && src.Rank-i > 0; i++ {
			dst := Position{src.File - i, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File-i > 0 && src.Rank+i < 10; i++ {
			dst := Position{src.File - i, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File+i < 10 && src.Rank-i > 0; i++ {
			dst := Position{src.File + i, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File+i < 10 && src.Rank+i < 10; i++ {
			dst := Position{src.File + i, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
	case RY:
		for _, d := range []Position{Position{1, 1}, Position{1, -1}, Position{-1, 1}, Position{-1, -1}} {
			positions = append(positions, Position{src.File + d.File, src.Rank + d.Rank})
		}
		fallthrough
	case HI:
		for i := 1; src.File+i < 10; i++ {
			dst := Position{src.File + i, src.Rank}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.File-i > 0; i++ {
			dst := Position{src.File - i, src.Rank}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.Rank+i < 10; i++ {
			dst := Position{src.File, src.Rank + i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
		for i := 1; src.Rank-i > 0; i++ {
			dst := Position{src.File, src.Rank - i}
			positions = append(positions, dst)
			if s.GetBoard(dst.File, dst.Rank) != nil {
				break
			}
		}
	case OU:
		for _, d := range []Position{
			Position{1, 1}, Position{1, -1}, Position{-1, 1}, Position{-1, -1},
			Position{0, 1}, Position{0, -1}, Position{-1, 0}, Position{1, 0},
		} {
			positions = append(positions, Position{src.File + d.File, src.Rank + d.Rank})
		}
	}
	// filtering to valid positions
	results := []Position{}
	for _, pos := range positions {
		if pos.File > 0 && pos.File < 10 && pos.Rank > 0 && pos.Rank < 10 {
			dstB := s.Board[pos.Rank-1][9-pos.File]
			if dstB != nil && dstB.Turn == b.Turn {
				continue
			}
			results = append(results, pos)
		}
	}
	return results
}
