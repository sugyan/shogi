package shogi

// Board definition
type Board [9][9]*Piece

// NewBoard function
func NewBoard() Board {
	board := [9][9]*Piece{}
	return board
}

func (b Board) String() string {
	result := make([]byte, 0, 9*(9*3+1))
	for _, row := range b {
		for _, p := range row {
			if p != nil {
				if p.First {
					result = append(result, '+')
				} else {
					result = append(result, '-')
				}
				result = append(result, []byte(p.Type)...)
			} else {
				result = append(result, []byte(` * `)...)
			}
		}
		result = append(result, '\n')
	}
	return string(result)
}
