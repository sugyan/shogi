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
	"+NY": shogi.BNY, "-NY": shogi.WNY,
	"+NK": shogi.BNK, "-NK": shogi.WNK,
	"+NG": shogi.BNG, "-NG": shogi.WNG,
	"+UM": shogi.BUM, "-UM": shogi.WUM,
	"+RY": shogi.BRY, "-RY": shogi.WRY,
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
		if len(line) == 0 {
			continue
		}
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
				*record.State = *shogi.InitialState
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
					err := record.State.SetPiece(9-i, int(line[1]-'1')+1, pieceMap[line[i*3+2:i*3+5]])
					if err != nil {
						return nil, err
					}
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
									file, rank = 9-j, i+1
									piece, err := record.State.GetPiece(file, rank)
									if err != nil {
										return nil, err
									}
									switch piece {
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
			src := shogi.Position{File: int(line[1] - '0'), Rank: int(line[2] - '0')}
			dst := shogi.Position{File: int(line[3] - '0'), Rank: int(line[4] - '0')}
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
