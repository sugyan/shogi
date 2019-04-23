package kif

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/record"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const stateSeparator = "+---------------------------+"
const movesSeparator = "手数----指手---------消費時間--"

var pieceMap = map[rune]shogi.Piece{
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
	"１":  1,
	"２":  2,
	"３":  3,
	"４":  4,
	"５":  5,
	"６":  6,
	"７":  7,
	"８":  8,
	"９":  9,
}

// Parse function
func Parse(r io.Reader) (*record.Record, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	reader, err := reader(b)
	if err != nil {
		return nil, err
	}

	type movesContent struct {
		index   int
		move    *shogi.Move
		comment *string
	}
	state := shogi.NewState()
	movesContents := []*movesContent{}
	var (
		prevMove  *shogi.Move
		boardFlag = false
		movesFlag = false
		row       = 0
		handsRE   = regexp.MustCompile(`(先|後)手の持駒：(.*)`)
		movesRE   = regexp.MustCompile(`(\d+)\s([１２３４５６７８９同]\S+)`)
		pieceRE   = regexp.MustCompile(`(?:[歩と金角馬飛龍玉]|成?[香桂銀])`)
		mvSrcRE   = regexp.MustCompile(`\(([1-9])([1-9])\)`)
	)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case stateSeparator:
			boardFlag = !boardFlag
			continue
		case movesSeparator:
			movesFlag = !movesFlag
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
						Piece: pieceMap[runes[i]],
					}
				}
				col++
				if col >= 9 {
					break
				}
			}
			row++
		}
		if movesFlag {
			if !movesRE.MatchString(line) {
				movesContents = append(movesContents, &movesContent{
					comment: &line,
				})
				continue
			}
			submatch := movesRE.FindStringSubmatch(line)
			index, _ := strconv.Atoi(submatch[1])
			runes := []rune(submatch[2])
			move := &shogi.Move{}
			move.Turn = shogi.TurnBlack
			if index%2 == 0 {
				move.Turn = shogi.TurnWhite
			}
			if mvSrcRE.MatchString(submatch[2]) {
				srcSubmatch := mvSrcRE.FindStringSubmatch(submatch[2])
				file, _ := strconv.Atoi(srcSubmatch[1])
				rank, _ := strconv.Atoi(srcSubmatch[2])
				move.Src = shogi.Pos(file, rank)
			}
			switch runes[0] {
			case '同':
				move.Dst = prevMove.Dst
			default:
				file := numberMap[string(runes[0])]
				rank := numberMap[string(runes[1])]
				move.Dst = shogi.Pos(file, rank)
			}
			piece := pieceRE.FindString(line)
			if piece == "" {
				return nil, errors.New("failed to parse moves")
			}
			move.Piece = pieceMap[[]rune(piece)[0]]
			if strings.Index(line[pieceRE.FindStringIndex(line)[0]:], "成") != -1 {
				move.Piece = map[shogi.Piece]shogi.Piece{
					shogi.FU: shogi.TO,
					shogi.KY: shogi.NY,
					shogi.KE: shogi.NK,
					shogi.GI: shogi.NG,
					shogi.KA: shogi.UM,
					shogi.HI: shogi.RY,
				}[move.Piece]
			}
			movesContents = append(movesContents, &movesContent{
				index: index,
				move:  move,
			})
			prevMove = move
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
					state.Captured[turn].Add(pieceMap[runes[0]])
				}
			}
		}
	}
	moves := []*shogi.Move{}
	finishCommentRE := regexp.MustCompile(`(?:解説|解答図)`)
	for _, mc := range movesContents {
		if mc.index != 0 && mc.move != nil {
			if mc.index > len(moves) {
				moves = append(moves, mc.move)
			} else {
				if mc.index == 1 {
					moves = []*shogi.Move{mc.move}
				} else {
					moves[mc.index-1] = mc.move
				}
			}
		} else {
			if finishCommentRE.MatchString(*mc.comment) {
				break
			}
		}
	}

	return &record.Record{
		State: state,
		Moves: moves,
	}, nil
}

func reader(b []byte) (io.Reader, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	if ok := scanner.Scan(); !ok {
		return nil, errors.New("failed to read header")
	}
	re, err := regexp.Compile(`^.*?#KIF version=([\d\.]+) encoding=(\S+).*?$`)
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
