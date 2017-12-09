package csa

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/record"
)

var codeMap = map[string]shogi.Piece{
	"FU": shogi.FU,
	"KY": shogi.KY,
	"KE": shogi.KE,
	"GI": shogi.GI,
	"KI": shogi.KI,
	"KA": shogi.KA,
	"HI": shogi.HI,
	"OU": shogi.OU,
	"TO": shogi.TO,
	"NY": shogi.NY,
	"NK": shogi.NK,
	"NG": shogi.NG,
	"UM": shogi.UM,
	"RY": shogi.RY,
}

// Parse function
func Parse(r io.Reader) (*record.Record, error) {
	piecesPattern := `(?:FU|KY|KE|GI|KI|KA|HI|OU|TO|NY|NK|NG|UM|RY)`
	rePos := regexp.MustCompile(`\d\d(?:` + piecesPattern + `|AL)`)
	reRow := regexp.MustCompile(`^P[1-9](?:[\+\-]` + piecesPattern + `| \* ){9}$`)
	reMove := regexp.MustCompile(`^[\+\-]\d{4}` + piecesPattern + `$`)

	state := shogi.NewState()
	moves := []*shogi.Move{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "P"):
			switch line[1] {
			case '+':
				fallthrough
			case '-':
				var turn shogi.Turn
				switch line[1] {
				case '+':
					turn = shogi.TurnBlack
				case '-':
					turn = shogi.TurnWhite
				}
				for _, e := range rePos.FindAllStringSubmatch(line, -1) {
					file, rank := int(e[0][0]-'0'), int(e[0][1]-'0')
					piece := codeMap[e[0][2:]]
					if !(file == 0 && rank == 0) {
						state.SetBoard(file, rank, &shogi.BoardPiece{Turn: turn, Piece: piece})
					} else {
						if e[0][2:] != "AL" {
							state.Captured[turn].Add(piece)
						} else {
							state.Captured[turn] = &shogi.CapturedPieces{
								FU: 18 - state.Captured[!turn].FU,
								KY: 4 - state.Captured[!turn].KY,
								KE: 4 - state.Captured[!turn].KE,
								GI: 4 - state.Captured[!turn].GI,
								KI: 4 - state.Captured[!turn].KI,
								KA: 2 - state.Captured[!turn].KA,
								HI: 2 - state.Captured[!turn].HI,
							}
							for i := 0; i < 9; i++ {
								for j := 0; j < 9; j++ {
									b := state.Board[i][j]
									if b != nil {
										state.Captured[turn].Sub(b.Piece)
									}
								}
							}
						}
					}
				}
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if reRow.MatchString(line) {
					rank := int(line[1] - '0')
					for i := 0; i < 9; i++ {
						e := line[2+3*i : 2+3*(i+1)]
						if e != " * " {
							var turn shogi.Turn
							switch e[0] {
							case '+':
								turn = shogi.TurnBlack
							case '-':
								turn = shogi.TurnWhite
							}
							file := 9 - i
							state.SetBoard(file, rank, &shogi.BoardPiece{Turn: turn, Piece: codeMap[e[1:]]})
						}
					}
				}
			case 'I':
				state.SetBoard(1, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KY})
				state.SetBoard(2, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KE})
				state.SetBoard(3, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.GI})
				state.SetBoard(4, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KI})
				state.SetBoard(5, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.OU})
				state.SetBoard(6, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KI})
				state.SetBoard(7, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.GI})
				state.SetBoard(8, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KE})
				state.SetBoard(9, 1, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KY})
				state.SetBoard(2, 2, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.KA})
				state.SetBoard(8, 2, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.HI})
				state.SetBoard(1, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(2, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(3, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(4, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(5, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(6, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(7, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(8, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(9, 3, &shogi.BoardPiece{Turn: shogi.TurnWhite, Piece: shogi.FU})
				state.SetBoard(1, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(2, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(3, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(4, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(5, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(6, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(7, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(8, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(9, 7, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.FU})
				state.SetBoard(2, 8, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.HI})
				state.SetBoard(8, 8, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KA})
				state.SetBoard(1, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KY})
				state.SetBoard(2, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KE})
				state.SetBoard(3, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.GI})
				state.SetBoard(4, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KI})
				state.SetBoard(5, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.OU})
				state.SetBoard(6, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KI})
				state.SetBoard(7, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.GI})
				state.SetBoard(8, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KE})
				state.SetBoard(9, 9, &shogi.BoardPiece{Turn: shogi.TurnBlack, Piece: shogi.KY})
				for _, e := range rePos.FindAllStringSubmatch(line, -1) {
					file, rank := int(e[0][0]-'0'), int(e[0][1]-'0')
					code := codeMap[e[0][2:]]
					if state.Board[rank-1][9-file].Piece == code {
						state.Board[rank-1][9-file] = nil
					}
				}
			}
		case strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-"):
			if reMove.MatchString(line) {
				move := &shogi.Move{}
				switch line[0] {
				case '+':
					move.Turn = shogi.TurnBlack
				case '-':
					move.Turn = shogi.TurnWhite
				}
				move.Src = shogi.Pos(int(line[1]-'0'), int(line[2]-'0'))
				move.Dst = shogi.Pos(int(line[3]-'0'), int(line[4]-'0'))
				move.Piece = codeMap[line[5:7]]
				moves = append(moves, move)
			}
		}
	}
	return &record.Record{
		State: state,
		Moves: moves,
	}, nil
}
