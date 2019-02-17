package csa

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"

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
	" * ": shogi.BLANK,
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
		Players: [2]*shogi.Player{{}, {}},
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
			if phase != phase2 {
				continue
			}
			switch line[1] {
			case '+':
				record.Players[0].Name = line[2:]
			case '-':
				record.Players[1].Name = line[2:]
			}
		case '$': // meta info
			if phase != phase2 {
				continue
			}
		case 'P': // initial positions
			switch line[1] {
			case 'I':
				phase = phase3_1
				// TODO
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				phase = phase3_2
				for i := 0; i < 9; i++ {
					if i*3+5 > len(line) {
						return nil, ErrInvalidLine
					}
					record.State.Board[line[1]-'1'][i] = pieceMap[line[i*3+2:i*3+5]]
				}
			case '+', '-':
				// TODO
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
			log.Printf("%v", scanner.Text())
		}
	}
	return record, nil
}
