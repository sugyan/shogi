package csa

import (
	"io"
	"regexp"

	"github.com/sugyan/shogi"
)

const piecesPattern = `(?:FU|KY|KE|GI|KI|KA|HI|OU|TO|NY|NK|NG|UM|RY)`

var (
	rePos  = regexp.MustCompile(`\d\d(?:` + piecesPattern + `|AL)`)
	reRow  = regexp.MustCompile(`^P[1-9](?:[\+\-]` + piecesPattern + `| \* ){9}$`)
	reMove = regexp.MustCompile(`^[\+\-]\d{4}` + piecesPattern + `$`)
)

// Parse function
func Parse(r io.Reader) (*shogi.Record, error) {
	record := &shogi.Record{
		State: &shogi.State{},
		Moves: []*shogi.Move{},
	}
	// TODO
	return record, nil
}
