package csa

import (
	"bufio"
	"bytes"
	"io"
	"log"

	"github.com/sugyan/shogi"
)

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
	scanner := bufio.NewScanner(p.r)
	for scanner.Scan() {
		line := scanner.Text()
		switch line[0] {
		case '\'':
			// comment
			continue
		case 'V':
			// version
			continue
		case 'N':
			// player names
			switch line[1] {
			case '+':
				record.Players[0].Namee = line[2:]
			case '-':
				record.Players[1].Namee = line[2:]
			}
		case '$':
			// meta info
		case 'P':
			// initial positions
			switch line[1] {
			case 'I':
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			case '+', '-':
			}
		case '+', '-':
			// moves
		case 'T':
			// consumed times
		case '%':
			// special case
		default:
			log.Printf("%v", scanner.Text())
		}
	}
	return record, nil
}
