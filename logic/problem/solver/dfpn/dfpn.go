package dfpn

import (
	"github.com/sugyan/shogi"
)

const inf = uint32(1) << 24

// Solver type
type Solver struct {
	hash map[string]*hash
}

type hash struct {
	pn uint32
	dn uint32
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
	if root.getPhi() < inf && root.getDelta() < inf {
		root.setPhi(inf)
		root.setDelta(inf)
		s.mid(root)
	}
}

func (s *Solver) mid(n *Node) {
	// 1. look up hash
	h := s.lookUpHash(n)
	if n.getPhi() <= h.pn || n.getDelta() <= h.dn {
		n.setPhi(h.pn)
		n.setDelta(h.dn)
		return
	}
	// 2. generate legal moves
	if len(n.Children) == 0 {
		for _, ms := range candidates(n.State, !n.Move.Turn) {
			n.Children = append(n.Children, &Node{
				Move:  ms.move,
				State: ms.state,
				pn:    1, dn: 1,
				parent: n,
			})
		}
	}
	if len(n.Children) == 0 {
		n.setPhi(inf)
		n.setDelta(0)
		switch n.Move.Turn {
		case shogi.TurnBlack:
			n.setResult(ResultT)
		case shogi.TurnWhite:
			n.setResult(ResultF)
		}
		s.putInHash(n)
		return
	}
	// 3. put in hash
	s.putInHash(n)
	// 4. multiple-iterative deepening
	for {
		p, d := n.getPhi(), n.getDelta()
		sp, md := s.sumPhi(n), s.minDelta(n)
		if p <= md || d <= sp {
			if n.Result == ResultU && (sp >= inf || md >= inf) {
				switch n.Move.Turn {
				case shogi.TurnBlack:
					if md >= inf {
						n.setResult(ResultT)
					} else {
						n.setResult(ResultF)
					}
				case shogi.TurnWhite:
					if sp >= inf {
						n.setResult(ResultF)
					} else {
						n.setResult(ResultT)
					}
				}
			}
			n.setPhi(md)
			n.setDelta(sp)
			s.putInHash(n)
			return
		}

		c, h, d2 := s.selectChild(n)
		if h.pn == inf-1 {
			c.setPhi(inf)
		} else if d >= inf-1 {
			c.setPhi(inf - 1)
		} else {
			c.setPhi(d + h.pn - sp)
		}
		if h.dn == inf-1 {
			c.setDelta(inf)
		} else {
			if d2 < inf {
				d2++
			}
			min := d2
			if p < min {
				min = p
			}
			c.setDelta(min)
		}
		s.mid(c)
		// if n.Result != ResultU {
		// 	return
		// }
	}
}

func (s *Solver) selectChild(n *Node) (*Node, *hash, uint32) {
	h := &hash{pn: inf, dn: inf}
	delta2 := inf
	var best *Node
	for _, child := range n.Children {
		hash := s.lookUpHash(child)
		if hash.dn < h.dn {
			best = child
			delta2 = h.dn
			h.pn = hash.pn
			h.dn = hash.dn
		} else if hash.dn < delta2 {
			delta2 = hash.dn
		}
		if hash.pn >= inf {
			return best, h, delta2
		}
	}
	return best, h, delta2
}

func (s *Solver) putInHash(n *Node) {
	key := n.State.Hash()
	s.hash[key] = &hash{
		pn: n.getPhi(),
		dn: n.getDelta(),
	}
}

func (s *Solver) lookUpHash(n *Node) *hash {
	key := n.State.Hash()
	if v, exist := s.hash[key]; exist {
		return v
	}
	return &hash{pn: 1, dn: 1}
}

func (s *Solver) sumPhi(n *Node) uint32 {
	sum := uint32(0)
	for _, child := range n.Children {
		h := s.lookUpHash(child)
		sum += h.pn
	}
	return sum
}

func (s *Solver) minDelta(n *Node) uint32 {
	min := inf
	for _, child := range n.Children {
		h := s.lookUpHash(child)
		if h.dn < min {
			min = h.dn
		}
	}
	return min
}
