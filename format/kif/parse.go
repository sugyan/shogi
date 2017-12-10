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

var codeMap = map[rune]shogi.Piece{
	'歩': shogi.FU,
	'と': shogi.TO,
	'香': shogi.KY,
	'桂': shogi.KE,
	'銀': shogi.GI,
	'金': shogi.KI,
	'角': shogi.KA,
	'馬': shogi.UM,
	'飛': shogi.HI,
	'龍': shogi.RY,
	'王': shogi.OU,
	'玉': shogi.OU,
}

var numberMap = map[string]int{
	"一":  1,
	"二":  2,
	"三":  3,
	"四":  4,
	"五":  5,
	"六":  6,
	"七":  7,
	"八":  8,
	"九":  9,
	"十":  10,
	"十一": 11,
	"十二": 12,
	"十三": 13,
	"十四": 14,
	"十五": 15,
	"十六": 16,
	"十七": 17,
	"十八": 18,
}

// Parse function
func Parse(r io.Reader) (*shogi.State, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	reader, err := reader(b)
	if err != nil {
		return nil, err
	}

	state := shogi.NewState()
	var (
		boardFlag = false
		row       = 0
	)
	handsRE := regexp.MustCompile(`(先|後)手の持駒：(.*)`)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if line == separator {
			boardFlag = !boardFlag
			continue
		}
		if boardFlag {
			runes := []rune(line)
			col := 0
			for i := 0; i < len(runes); i++ {
				c := runes[i]
				if c == ' ' || c == '|' {
					continue
				}
				if c == '・' {
					state.Board[row][col] = nil
				} else {
					turn := shogi.TurnBlack
					if c == 'v' {
						turn = shogi.TurnWhite
						i++
					}
					state.Board[row][col] = &shogi.BoardPiece{
						Turn:  turn,
						Piece: codeMap[runes[i]],
					}
				}
				col++
				if col >= 9 {
					break
				}
			}
			row++
		}
		if handsRE.MatchString(line) {
			spacesRE := regexp.MustCompile(`[ 　]+`)
			submatch := handsRE.FindStringSubmatch(line)
			for _, s := range spacesRE.Split(submatch[2], -1) {
				if len(s) == 0 {
					continue
				}
				runes := []rune(s)
				var turn shogi.Turn
				switch submatch[1] {
				case "先":
					turn = shogi.TurnBlack
				case "後":
					turn = shogi.TurnWhite
				}
				n := 1
				if len(runes[1:]) > 0 {
					n = numberMap[string(runes[1:])]
				}
				for i := 0; i < n; i++ {
					state.Captured[turn].Add(codeMap[runes[0]])
				}
			}
		}
	}
	return state, nil
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
	if len(submatch) > 2 {
		switch submatch[2] {
		case "Shift_JIS":
			return transform.NewReader(bytes.NewBuffer(b), japanese.ShiftJIS.NewDecoder()), nil
		default:
			return bytes.NewBuffer(b), nil
		}
	}
	// default...?
	return transform.NewReader(bytes.NewBuffer(b), japanese.ShiftJIS.NewDecoder()), nil
}
