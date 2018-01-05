package node

import (
	"fmt"
	"strings"

	"github.com/sugyan/shogi"
)

// TreeString function
func TreeString(n Node) string {
	return archy(n, "")[1:]
}

func archy(n Node, prefix string) string {
	var (
		line   string
		result rune
	)
	switch n.Result() {
	case ResultU:
		result = '?'
	case ResultT:
		result = 'T'
	case ResultF:
		result = 'F'
	}
	if prefix != "" {
		move := n.Move()
		sign := '+'
		if move.Turn == shogi.TurnWhite {
			sign = '-'
		}
		line = fmt.Sprintf("%c%d%d%d%d%s (%c)", sign, move.Src.File, move.Src.Rank, move.Dst.File, move.Dst.Rank, move.Piece, result)
	} else {
		line = fmt.Sprintf(" %c", result)
	}

	nodeLines := []string{}
	for i, c := range n.Children() {
		_prefix := prefix + " "
		if i < len(n.Children())-1 {
			_prefix = prefix + "│"
		}
		_prefix += " "

		l := "└"
		if i < len(n.Children())-1 {
			l = "├"
		}
		m := "─"
		if len(c.Children()) > 0 {
			m = "┬"
		}
		line := prefix +
			l + "─" +
			m + " " +
			string([]rune(archy(c, _prefix))[len([]rune(prefix))+2:])
		nodeLines = append(nodeLines, line)
	}
	return prefix +
		line + "\n" +
		strings.Join(nodeLines, "")
}
