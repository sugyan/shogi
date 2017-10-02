package kif

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/sugyan/shogi"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const separator = "+---------------------------+"

// Parse function
func Parse(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	reader, err := reader(b)
	if err != nil {
		return err
	}

	board := shogi.NewBoard()
	var (
		boardFlag = false
		row       = 0
	)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == separator {
			boardFlag = !boardFlag
			continue
		}
		runes := []rune(line)
		if boardFlag {
			col := 0
			for i := 0; i < len(runes); i++ {
				c := runes[i]
				if c == ' ' || c == '|' {
					continue
				}
				if c == '・' {
					board[row][col] = nil
				} else {
					p := shogi.Piece{First: true}
					if c == 'v' {
						p.First = false
						i++
					}
					switch runes[i] {
					case '歩':
						p.Type = shogi.FU
					case '香':
						p.Type = shogi.KY
					case '桂':
						p.Type = shogi.KE
					case '銀':
						p.Type = shogi.GI
					case '金':
						p.Type = shogi.KI
					case '角':
						p.Type = shogi.KA
					case '馬':
						p.Type = shogi.UM
					case '飛':
						p.Type = shogi.HI
					case '龍':
						p.Type = shogi.RY
					case '王':
						fallthrough
					case '玉':
						p.Type = shogi.OU
					}
					board[row][col] = &p
				}
				col++
				if col >= 9 {
					break
				}
			}
			row++
		}
	}
	print(board.String())

	return nil
}

func reader(b []byte) (io.Reader, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	if ok := scanner.Scan(); !ok {
		return nil, errors.New("failed to read header")
	}
	re, err := regexp.Compile(`^#KIF version=([\d\.]+) encoding=(\S+)$`)
	if err != nil {
		return nil, err
	}
	submatch := re.FindStringSubmatch(scanner.Text())
	if len(submatch) < 3 {
		return nil, errors.New("failed to parse header")
	}
	switch submatch[2] {
	case "Shift_JIS":
		return transform.NewReader(bytes.NewBuffer(b), japanese.ShiftJIS.NewDecoder()), nil
	default:
		return bytes.NewBuffer(b), nil
	}
}
