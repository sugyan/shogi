package csa

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/sugyan/shogi"
)

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
					code := shogi.PieceCode(e[0][2:])
					if !(file == 0 && rank == 0) {
						piece := shogi.NewPiece(turn, code)
						state.SetPiece(file, rank, piece)
					} else {
						if e[0][2:] != "AL" {
							piece := shogi.NewPiece(turn, code)
							state.AddCapturedPieces(piece)
						} else {
							state.Captured[turn] = &shogi.CapturedPieces{
								Fu: 18 - state.Captured[!turn].Fu,
								Ky: 4 - state.Captured[!turn].Ky,
								Ke: 4 - state.Captured[!turn].Ke,
								Gi: 4 - state.Captured[!turn].Gi,
								Ki: 4 - state.Captured[!turn].Ki,
								Ka: 2 - state.Captured[!turn].Ka,
								Hi: 2 - state.Captured[!turn].Hi,
							}
							for i := 0; i < 9; i++ {
								for j := 0; j < 9; j++ {
									p := state.Board[i][j]
									if p != nil {
										switch shogi.PieceCode(p.Code()) {
										case shogi.FU:
											fallthrough
										case shogi.TO:
											state.Captured[turn].Fu--
										case shogi.KY:
											fallthrough
										case shogi.NY:
											state.Captured[turn].Ky--
										case shogi.KE:
											fallthrough
										case shogi.NK:
											state.Captured[turn].Ke--
										case shogi.GI:
											fallthrough
										case shogi.NG:
											state.Captured[turn].Gi--
										case shogi.KI:
											state.Captured[turn].Ki--
										case shogi.KA:
											fallthrough
										case shogi.UM:
											state.Captured[turn].Ka--
										case shogi.HI:
											fallthrough
										case shogi.RY:
											state.Captured[turn].Hi--
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
							piece := shogi.NewPiece(turn, shogi.PieceCode(e[1:]))
							state.SetPiece(9-i, rank, piece)
						}
					}
				}
			case 'I':
				// TODO
			}
		}
	}
	return state, nil
}
