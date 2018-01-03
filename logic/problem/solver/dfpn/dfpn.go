package dfpn

import (
	"github.com/sugyan/shogi"
)

const inf = uint32(1) << 12

type hash struct {
	turn   shogi.Turn
	pn, dn uint32
}

func (h *hash) p() uint32 {
	switch h.turn {
	case shogi.TurnBlack:
		return h.dn
	case shogi.TurnWhite:
		return h.pn
	}
	return 0
}

func (h *hash) d() uint32 {
	switch h.turn {
	case shogi.TurnBlack:
		return h.pn
	case shogi.TurnWhite:
		return h.dn
	}
	return 0
}

// Solver type
type Solver struct {
	hash     map[string]*hash
	maxDepth int
}

// NewSolver function
func NewSolver() *Solver {
	return &Solver{
		hash: map[string]*hash{},
	}
}

// Solve method
func (s *Solver) Solve(root *Node) {
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
func (s *Solver) SetMaxDepth(d int) {
	s.maxDepth = d
}

func (s *Solver) mid(n *Node) {
	h := s.lookUpHash(n)
	if n.pn <= h.pn || n.dn <= h.dn {
		n.pn = h.pn
		n.dn = h.dn
		if n.pn == 0 {
			n.Result = ResultT
		}
		if n.dn == 0 {
			n.Result = ResultF
		}
		return
	}
	if n.expanded {
		if len(n.Children) == 0 {
			n.setP(inf)
			n.setD(0)
			s.putInHash(n, n.pn, n.dn)
			return
		}
	} else {
		if !(s.maxDepth != 0 && n.depth > s.maxDepth && n.Move.Turn == shogi.TurnWhite) {
			for _, ms := range candidates(n.State, !n.Move.Turn) {
				n.Children = append(n.Children, &Node{
					Move:  ms.move,
					State: ms.state,
					depth: n.depth + 1,
				})
			}
		}
		n.expanded = true
	}
	switch n.Move.Turn {
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
				n.Result = ResultT
			}
			if n.dn == 0 {
				n.Result = ResultF
			}
			return
		}
		best, cp, cd, d2 := s.selectChild(n)
		c := n.Children[best]
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

func (s *Solver) lookUpHash(n *Node) *hash {
	if h, ok := s.hash[n.State.Hash()]; ok {
		return h
	}
	return &hash{n.Move.Turn, 1, 1}
}

func (s *Solver) putInHash(n *Node, pn, dn uint32) {
	s.hash[n.State.Hash()] = &hash{
		turn: n.Move.Turn,
		pn:   pn, dn: dn,
	}
}

func (s *Solver) minDelta(n *Node) uint32 {
	min := inf
	for _, c := range n.Children {
		h := s.lookUpHash(c)
		d := h.d()
		if d < min {
			min = d
		}
	}
	return min
}

func (s *Solver) sumPhi(n *Node) uint32 {
	sum := uint32(0)
	for _, c := range n.Children {
		h := s.lookUpHash(c)
		sum += h.p()
	}
	return sum
}

func (s *Solver) selectChild(n *Node) (int, uint32, uint32, uint32) {
	d2 := inf
	pn, dn := inf, inf
	best := 0
	for i, c := range n.Children {
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
