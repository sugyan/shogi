package record

import "github.com/sugyan/shogi"

// Record type
type Record struct {
	State *shogi.State
	Moves []*shogi.Move
}

// Converter interface
type Converter interface {
	ConvertToString(*Record) string
}

// ConvertToString method
func (r *Record) ConvertToString(c Converter) string {
	return c.ConvertToString(r)
}
