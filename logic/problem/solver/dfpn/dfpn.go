package dfpn

import (
	"github.com/sugyan/shogi"
	"github.com/sugyan/shogi/logic/problem/solver/node"
)

const inf = uint32(1) << 12

// Searcher type
type Searcher struct {
	hash     map[string]*hash
	maxDepth int
}

// NewSearcher function
func NewSearcher() *Searcher {
	return &Searcher{
		hash: map[string]*hash{},
	}
}

// Search method
func (s *Searcher) Search(root *Node) {
	root.pn = inf - 1
	root.dn = inf - 1
	s.mid(root)
	if root.getP() < inf && root.getD() < inf {
		root.pn = inf
		root.dn = inf
		s.mid(root)
	}
}

// SetMaxDepth method
func (s *Searcher) SetMaxDepth(d int) {
	s.maxDepth = d
}

func (s *Searcher) mid(n *Node) {
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
		if !(s.maxDepth != 0 && n.move.Turn == shogi.TurnWhite && n.depth+1 > s.maxDepth) {
			for _, ms := range candidates(n.state, !n.move.Turn) {
				// checkmating with dropping FU
				if ms.move.Turn == shogi.TurnBlack && ms.move.Src.IsCaptured() && ms.move.Piece == shogi.FU {
					if len(candidates(ms.state, shogi.TurnWhite)) == 0 {
						continue
					}
				}
				n.children = append(n.children, &Node{
					move:  ms.move,
					state: ms.state,
					depth: n.depth + 1,
					hash:  ms.state.Hash(),
				})
			}
		}
		n.expanded = true
	}
	switch n.move.Turn {
	case shogi.TurnBlack:
		s.putInHash(n, inf, 0)
	case shogi.TurnWhite:
		s.putInHash(n, 0, inf)
	}
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

func (s *Searcher) lookUpHash(n *Node) *hash {
	if h, ok := s.hash[n.hash]; ok {
		return h
	}
	return &hash{n.move.Turn, 1, 1}
}

func (s *Searcher) putInHash(n *Node, pn, dn uint32) {
	s.hash[n.hash] = &hash{
		turn: n.move.Turn,
		pn:   pn, dn: dn,
	}
}

func (s *Searcher) minDelta(n *Node) uint32 {
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

func (s *Searcher) sumPhi(n *Node) uint32 {
	sum := uint32(0)
	for _, c := range n.children {
		h := s.lookUpHash(c)
		sum += h.p()
	}
	return sum
}

func (s *Searcher) selectChild(n *Node) (int, uint32, uint32, uint32) {
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
		} else if h.pn < d2 {
			d2 = h.d()
		}
		if h.p() == inf {
			return best, pn, dn, d2
		}
	}
	return best, pn, dn, d2
}
