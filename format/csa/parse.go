package csa

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/sugyan/shogi"
)

var codeMap = map[string]shogi.PieceCode{
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
func Parse(r io.Reader) (*shogi.State, error) {
	rePos := regexp.MustCompile(`\d\d(?:FU|KY|KE|GI|KI|KA|HI|OU|TO|NY|NK|NG|UM|RY|AL)`)
	reRow := regexp.MustCompile(`^P[1-9](?:[\+\-](?:FU|KY|KE|GI|KI|KA|HI|OU|TO|NY|NK|NG|UM|RY)| \* ){9}$`)

	state := shogi.NewState()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "P") {
			switch line[1] {
			case '+':
				fallthrough
			case '-':
				var turn shogi.Turn
				switch line[1] {
				case '+':
					turn = shogi.TurnFirst
				case '-':
					turn = shogi.TurnSecond
				}
				for _, e := range rePos.FindAllStringSubmatch(line, -1) {
					file, rank := int(e[0][0]-'0'), int(e[0][1]-'0')
					code := codeMap[e[0][2:]]
					if !(file == 0 && rank == 0) {
						state.SetBoardPiece(file, rank, turn, shogi.NewPiece(code))
					} else {
						if e[0][2:] != "AL" {
							state.Captured[turn].AddPieces(shogi.NewPiece(code))
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
									bp := state.Board[i][j]
									if bp != nil {
										switch bp.Piece.Code() {
										case shogi.FU:
											fallthrough
										case shogi.TO:
											state.Captured[turn].FU--
										case shogi.KY:
											fallthrough
										case shogi.NY:
											state.Captured[turn].KY--
										case shogi.KE:
											fallthrough
										case shogi.NK:
											state.Captured[turn].KE--
										case shogi.GI:
											fallthrough
										case shogi.NG:
											state.Captured[turn].GI--
										case shogi.KI:
											state.Captured[turn].KI--
										case shogi.KA:
											fallthrough
										case shogi.UM:
											state.Captured[turn].KA--
										case shogi.HI:
											fallthrough
										case shogi.RY:
											state.Captured[turn].HI--
										}
									}
								}
							}
						}
					}
				}
			case '1':
				fallthrough
			case '2':
				fallthrough
			case '3':
				fallthrough
			case '4':
				fallthrough
			case '5':
				fallthrough
			case '6':
				fallthrough
			case '7':
				fallthrough
			case '8':
				fallthrough
			case '9':
				if reRow.MatchString(line) {
					rank := int(line[1] - '0')
					for i := 0; i < 9; i++ {
						e := line[2+3*i : 2+3*(i+1)]
						if e != " * " {
							var turn shogi.Turn
							switch e[0] {
							case '+':
								turn = shogi.TurnFirst
							case '-':
								turn = shogi.TurnSecond
							}
							file := 9 - i
							state.SetBoardPiece(file, rank, turn, shogi.NewPiece(codeMap[e[1:]]))
						}
					}
				}
			case 'I':
				state.SetBoardPiece(1, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KY))
				state.SetBoardPiece(2, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KE))
				state.SetBoardPiece(3, 1, shogi.TurnSecond, shogi.NewPiece(shogi.GI))
				state.SetBoardPiece(4, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KI))
				state.SetBoardPiece(5, 1, shogi.TurnSecond, shogi.NewPiece(shogi.OU))
				state.SetBoardPiece(6, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KI))
				state.SetBoardPiece(7, 1, shogi.TurnSecond, shogi.NewPiece(shogi.GI))
				state.SetBoardPiece(8, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KE))
				state.SetBoardPiece(9, 1, shogi.TurnSecond, shogi.NewPiece(shogi.KY))
				state.SetBoardPiece(2, 2, shogi.TurnSecond, shogi.NewPiece(shogi.KA))
				state.SetBoardPiece(8, 2, shogi.TurnSecond, shogi.NewPiece(shogi.HI))
				state.SetBoardPiece(1, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(2, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(3, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(4, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(5, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(6, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(7, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(8, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(9, 3, shogi.TurnSecond, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(1, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(2, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(3, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(4, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(5, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(6, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(7, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(8, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(9, 7, shogi.TurnFirst, shogi.NewPiece(shogi.FU))
				state.SetBoardPiece(2, 8, shogi.TurnFirst, shogi.NewPiece(shogi.HI))
				state.SetBoardPiece(8, 8, shogi.TurnFirst, shogi.NewPiece(shogi.KA))
				state.SetBoardPiece(1, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KY))
				state.SetBoardPiece(2, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KE))
				state.SetBoardPiece(3, 9, shogi.TurnFirst, shogi.NewPiece(shogi.GI))
				state.SetBoardPiece(4, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KI))
				state.SetBoardPiece(5, 9, shogi.TurnFirst, shogi.NewPiece(shogi.OU))
				state.SetBoardPiece(6, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KI))
				state.SetBoardPiece(7, 9, shogi.TurnFirst, shogi.NewPiece(shogi.GI))
				state.SetBoardPiece(8, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KE))
				state.SetBoardPiece(9, 9, shogi.TurnFirst, shogi.NewPiece(shogi.KY))
				for _, e := range rePos.FindAllStringSubmatch(line, -1) {
					file, rank := int(e[0][0]-'0'), int(e[0][1]-'0')
					code := codeMap[e[0][2:]]
					if state.Board[rank-1][9-file].Piece.Code() == code {
						state.Board[rank-1][9-file] = nil
					}
				}
			}
		}
	}
	return state, nil
}
