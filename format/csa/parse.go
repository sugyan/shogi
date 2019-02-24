package csa

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/sugyan/shogi"
)

// ErrInvalidLine is error
var ErrInvalidLine = errors.New("invalid line")

type phase int

const (
	phase1 phase = iota
	phase2
	phase3_1
	phase3_2
	phase4
)

var pieceMap = map[string]shogi.Piece{
	" * ": shogi.EMP,
	"+FU": shogi.BFU, "-FU": shogi.WFU,
	"+KY": shogi.BKY, "-KY": shogi.WKY,
	"+KE": shogi.BKE, "-KE": shogi.WKE,
	"+GI": shogi.BGI, "-GI": shogi.WGI,
	"+KI": shogi.BKI, "-KI": shogi.WKI,
	"+KA": shogi.BKA, "-KA": shogi.WKA,
	"+HI": shogi.BHI, "-HI": shogi.WHI,
	"+OU": shogi.BOU, "-OU": shogi.WOU,
	"+TO": shogi.BTO, "-TO": shogi.WTO,
	"+NY": shogi.BKY, "-NY": shogi.WKY,
	"+NK": shogi.BKY, "-NK": shogi.WKY,
	"+NG": shogi.BKY, "-NG": shogi.WKY,
	"+UM": shogi.BKY, "-UM": shogi.WKY,
	"+RY": shogi.BKY, "-RY": shogi.WKY,
}

type parser struct {
	r io.Reader
}

// Parse function
func Parse(r io.Reader) (*shogi.Record, error) {
	p := parser{r: r}
	return p.parse()
}

// ParseString function
func ParseString(s string) (*shogi.Record, error) {
	return Parse(bytes.NewBufferString(s))
}

func (p *parser) parse() (*shogi.Record, error) {
	record := &shogi.Record{
		Players: [2]*shogi.Player{},
		State:   &shogi.State{},
		Moves:   []*shogi.Move{},
	}
	phase := phase1
	scanner := bufio.NewScanner(p.r)
	for scanner.Scan() {
		line := scanner.Text()
		switch line[0] {
		case '\'': // comment
			continue
		case 'V': // version
			if phase != phase1 {
				continue
			}
			phase = phase2
			continue
		case 'N': // player names
			if phase > phase2 {
				continue
			}
			switch line[1] {
			case '+':
				record.Players[0] = &shogi.Player{Name: line[2:]}
			case '-':
				record.Players[1] = &shogi.Player{Name: line[2:]}
			}
		case '$': // meta info
			if phase != phase2 {
				continue
			}
		case 'P': // initial positions
			switch line[1] {
			case 'I':
				if phase == phase3_2 {
					continue
				}
				phase = phase3_1
				record.State.SetPiece(1, 1, shogi.WKY)
				record.State.SetPiece(2, 1, shogi.WKE)
				record.State.SetPiece(3, 1, shogi.WGI)
				record.State.SetPiece(4, 1, shogi.WKI)
				record.State.SetPiece(5, 1, shogi.WOU)
				record.State.SetPiece(6, 1, shogi.WKI)
				record.State.SetPiece(7, 1, shogi.WGI)
				record.State.SetPiece(8, 1, shogi.WKE)
				record.State.SetPiece(9, 1, shogi.WKY)
				record.State.SetPiece(2, 2, shogi.WKA)
				record.State.SetPiece(8, 2, shogi.WHI)
				for i := 1; i < 10; i++ {
					record.State.SetPiece(i, 3, shogi.WFU)
					record.State.SetPiece(i, 7, shogi.BFU)
				}
				record.State.SetPiece(2, 8, shogi.BHI)
				record.State.SetPiece(8, 8, shogi.BKA)
				record.State.SetPiece(1, 9, shogi.BKY)
				record.State.SetPiece(2, 9, shogi.BKE)
				record.State.SetPiece(3, 9, shogi.BGI)
				record.State.SetPiece(4, 9, shogi.BKI)
				record.State.SetPiece(5, 9, shogi.BOU)
				record.State.SetPiece(6, 9, shogi.BKI)
				record.State.SetPiece(7, 9, shogi.BGI)
				record.State.SetPiece(8, 9, shogi.BKE)
				record.State.SetPiece(9, 9, shogi.BKY)
				for i := 0; i+2 < len(line); i += 4 {
					file, rank := int(line[i+2]-'0'), int(line[i+3]-'0')
					// TODO: check piece?
					record.State.SetPiece(file, rank, shogi.EMP)
				}
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				if phase == phase3_1 {
					continue
				}
				phase = phase3_2
				for i := 0; i < 9; i++ {
					if i*3+5 > len(line) {
						return nil, ErrInvalidLine
					}
					record.State.Board[line[1]-'1'][i] = pieceMap[line[i*3+2:i*3+5]]
				}
			case '+', '-':
				if phase == phase3_1 {
					continue
				}
				phase = phase3_2
				capturedIndex := 0
				if line[1] == '-' {
					capturedIndex = 1
				}
				for i := 0; i+2 < len(line); i += 4 {
					if i+6 > len(line) {
						return nil, ErrInvalidLine
					}
					file, rank := int(line[i+2]-'0'), int(line[i+3]-'0')
					if file == 0 && rank == 0 {
						switch line[i+4 : i+6] {
						case "FU":
							record.State.Captured[capturedIndex].FU++
						case "KY":
							record.State.Captured[capturedIndex].KY++
						case "KE":
							record.State.Captured[capturedIndex].KE++
						case "GI":
							record.State.Captured[capturedIndex].GI++
						case "KI":
							record.State.Captured[capturedIndex].KI++
						case "KA":
							record.State.Captured[capturedIndex].KA++
						case "HI":
							record.State.Captured[capturedIndex].HI++
						case "AL":
							record.State.Captured[capturedIndex] = shogi.Captured{
								FU: 18 - record.State.Captured[1-capturedIndex].FU,
								KY: 4 - record.State.Captured[1-capturedIndex].KY,
								KE: 4 - record.State.Captured[1-capturedIndex].KE,
								GI: 4 - record.State.Captured[1-capturedIndex].GI,
								KI: 4 - record.State.Captured[1-capturedIndex].KI,
								KA: 2 - record.State.Captured[1-capturedIndex].KA,
								HI: 2 - record.State.Captured[1-capturedIndex].HI,
							}
							for i := 0; i < 9; i++ {
								for j := 0; j < 9; j++ {
									switch record.State.Board[i][j] {
									case shogi.BFU, shogi.WFU, shogi.BTO, shogi.WTO:
										record.State.Captured[capturedIndex].FU--
									case shogi.BKY, shogi.WKY, shogi.BNY, shogi.WNY:
										record.State.Captured[capturedIndex].KY--
									case shogi.BKE, shogi.WKE, shogi.BNK, shogi.WNK:
										record.State.Captured[capturedIndex].KE--
									case shogi.BGI, shogi.WGI, shogi.BNG, shogi.WNG:
										record.State.Captured[capturedIndex].GI--
									case shogi.BKI, shogi.WKI:
										record.State.Captured[capturedIndex].KI--
									case shogi.BKA, shogi.WKA, shogi.BUM, shogi.WUM:
										record.State.Captured[capturedIndex].KA--
									case shogi.BHI, shogi.WHI, shogi.BRY, shogi.WRY:
										record.State.Captured[capturedIndex].HI--
									}
								}
							}
						}
					} else {
						piece := pieceMap[string(line[1])+line[i+4:i+6]]
						record.State.SetPiece(file, rank, piece)
					}
				}
			}
		case '+', '-': // moves
			if len(line) == 1 {
				// first move
				phase = phase4
				continue
			}
			if phase != phase4 {
				continue
			}
			src := shogi.Position{File: line[1] - '0', Rank: line[2] - '0'}
			dst := shogi.Position{File: line[3] - '0', Rank: line[4] - '0'}
			piece := pieceMap[string(line[0])+line[5:7]]
			record.Moves = append(record.Moves, &shogi.Move{
				Src:   src,
				Dst:   dst,
				Piece: piece,
			})
		case 'T': // consumed times
		case '%': // special case
		default:
			return nil, ErrInvalidLine
		}
	}
	return record, nil
}
