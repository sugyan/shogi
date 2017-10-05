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
					state.Board[row][col] = newPiece(runes[i], turn)
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
					turn = shogi.TurnFirst
				case "後":
					turn = shogi.TurnSecond
				}
				n := 1
				if len(runes[1:]) > 0 {
					n = numberMap[string(runes[1:])]
				}
				for i := 0; i < n; i++ {
					state.AddCapturedPieces(newPiece(runes[0], turn))
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

func newPiece(c rune, turn shogi.Turn) shogi.Piece {
	switch c {
	case '歩':
		return shogi.NewPiece(turn, shogi.FU)
	case 'と':
		return shogi.NewPiece(turn, shogi.TO)
	case '香':
		return shogi.NewPiece(turn, shogi.KY)
	case '桂':
		return shogi.NewPiece(turn, shogi.KE)
	case '銀':
		return shogi.NewPiece(turn, shogi.GI)
	case '金':
		return shogi.NewPiece(turn, shogi.KI)
	case '角':
		return shogi.NewPiece(turn, shogi.KA)
	case '馬':
		return shogi.NewPiece(turn, shogi.UM)
	case '飛':
		return shogi.NewPiece(turn, shogi.HI)
	case '龍':
		return shogi.NewPiece(turn, shogi.RY)
	case '王':
		return shogi.NewPiece(turn, shogi.OU)
	case '玉':
		return shogi.NewPiece(turn, shogi.OU)
	}
	return nil
}
