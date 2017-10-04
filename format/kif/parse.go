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

var pieceMap = map[rune]shogi.Piece{
	'歩': shogi.NewPiece(shogi.TurnFirst, shogi.FU),
	'香': shogi.NewPiece(shogi.TurnFirst, shogi.KY),
	'桂': shogi.NewPiece(shogi.TurnFirst, shogi.KE),
	'銀': shogi.NewPiece(shogi.TurnFirst, shogi.GI),
	'金': shogi.NewPiece(shogi.TurnFirst, shogi.KI),
	'角': shogi.NewPiece(shogi.TurnFirst, shogi.KA),
	'馬': shogi.NewPiece(shogi.TurnFirst, shogi.UM),
	'飛': shogi.NewPiece(shogi.TurnFirst, shogi.HI),
	'龍': shogi.NewPiece(shogi.TurnFirst, shogi.RY),
	'王': shogi.NewPiece(shogi.TurnFirst, shogi.OU),
	'玉': shogi.NewPiece(shogi.TurnFirst, shogi.OU),
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
					turn := shogi.TurnFirst
					if c == 'v' {
						turn = shogi.TurnSecond
						i++
					}
					p := pieceMap[runes[i]]
					p.SetTurn(turn)
					state.Board[row][col] = p
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
				p := pieceMap[runes[0]]
				switch submatch[1] {
				case "先":
					p.SetTurn(shogi.TurnFirst)
				case "後":
					p.SetTurn(shogi.TurnSecond)
				}
				n := 1
				if len(runes[1:]) > 0 {
					n = numberMap[string(runes[1:])]
				}
				for i := 0; i < n; i++ {
					state.AddCapturedPieces(p)
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
