package node

import (
	"github.com/sugyan/shogi"
)

// Result type
type Result int

// Resulst constants
const (
	ResultU = iota // Unknown
	ResultT        // True
	ResultF        // False
)

// Node interface
type Node interface {
	Children() []Node
	Result() Result
	Move() *shogi.Move
	State() *shogi.State
	Hash() string
}
