package dfpn

import (
	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

const inf = uint32(10000)

type solver struct {
	hash     map[string]*hash
	maxDepth int
}

// Solve function
func Solve(root *Node, maxDepth int) {
	s := &solver{
		hash:     map[string]*hash{},
		maxDepth: maxDepth,
	}
	root.pn = inf - 1
	root.dn = inf - 1
	s.mid(root)
	if root.pn < inf && root.dn < inf {
		root.pn = inf
		root.dn = inf
		s.mid(root)
	}
	// could not solved perfectly...
	if root.pn > 0 && root.dn > 0 {
		if root.dn == 0 {
			root.result = node.ResultT
		}
		if root.pn == 0 {
			root.result = node.ResultF
		}
	}

	if root.result == node.ResultU {
		// ??
		root.result = node.ResultF
	}
}

func (s *solver) mid(n *Node) {
	h := s.lookUpHash(n)
	if n.pn <= h.pn || n.dn <= h.dn {
		n.pn = h.pn
		n.dn = h.dn
		if n.pn == 0 {
			n.result = node.ResultT
		}
		if n.dn == 0 {
			n.result = node.ResultF
		}
		return
	}
	if n.expanded {
		if len(n.children) == 0 {
			n.setP(inf)
			n.setD(0)
			s.putInHash(n, n.pn, n.dn)
			return
		}
	} else {
		cutoff := false
		if s.maxDepth != 0 && n.depth+1 > s.maxDepth {
			if n.move.Turn == shogi.TurnWhite && !n.move.Src.IsCaptured() {
				cutoff = true
			}
		}
		if !cutoff {
			for _, ms := range problem.Candidates(n.state, !n.move.Turn) {
				// checkmating with dropping FU
				if ms.Move.Turn == shogi.TurnBlack && ms.Move.Src.IsCaptured() && ms.Move.Piece == shogi.FU {
					if len(problem.Candidates(ms.State, shogi.TurnWhite)) == 0 {
						continue
					}
				}
				depth := n.depth + 1
				if n.move.Turn == shogi.TurnWhite && n.move.Src.IsCaptured() && ms.Move.Dst == n.move.Dst {
					depth -= 2
				}
				n.children = append(n.children, &Node{
					move:  ms.Move,
					state: ms.State,
					depth: depth,
					hash:  ms.State.Hash(),
				})
			}
		}
		n.expanded = true
	}
	s.putInHash(n, n.dn, n.pn)
	for {
		minD := s.minDelta(n)
		sumP := s.sumPhi(n)
		if n.getP() <= minD || n.getD() <= sumP {
			n.setP(minD)
			n.setD(sumP)
			s.putInHash(n, n.pn, n.dn)
			if n.pn == 0 {
				n.result = node.ResultT
			}
			if n.dn == 0 {
				n.result = node.ResultF
			}
			return
		}
		best, cp, cd, d2 := s.selectChild(n)
		c := n.children[best]
		if cp == inf-1 {
			c.setP(inf)
		} else if n.getD() >= inf-1 {
			c.setP(inf - 1)
		} else {
			c.setP(n.getD() + cp - sumP)
		}
		if cd == inf-1 {
			c.setD(inf)
		} else {
			_dn := d2 + 1
			if n.getP() < _dn {
				_dn = n.getP()
			}
			c.setD(_dn)
		}
		s.mid(c)
	}
}

func (s *solver) lookUpHash(n *Node) *hash {
	if h, ok := s.hash[n.hash]; ok {
		return h
	}
	return &hash{n.move.Turn, 1, 1}
}

func (s *solver) putInHash(n *Node, pn, dn uint32) {
	s.hash[n.hash] = &hash{
		turn: n.move.Turn,
		pn:   pn, dn: dn,
	}
}

func (s *solver) minDelta(n *Node) uint32 {
	min := inf
	for _, c := range n.children {
		h := s.lookUpHash(c)
		d := h.d()
		if d < min {
			min = d
		}
	}
	return min
}

func (s *solver) sumPhi(n *Node) uint32 {
	sum := uint32(0)
	for _, c := range n.children {
		h := s.lookUpHash(c)
		sum += h.p()
	}
	return sum
}

func (s *solver) selectChild(n *Node) (int, uint32, uint32, uint32) {
	d2 := inf
	pn, dn := inf, inf
	best := 0
	for i, c := range n.children {
		h := s.lookUpHash(c)
		if h.d() < dn {
			best = i
			d2 = dn
			pn = h.p()
			dn = h.d()
		} else if h.d() < d2 {
			d2 = h.d()
		}
		if h.p() == inf {
			return best, pn, dn, d2
		}
	}
	return best, pn, dn, d2
}
