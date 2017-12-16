package dfpn

import (
	"log"

	"github.com/sugyan/shogi"
)

const inf = uint32(1) << 31

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
func (s *Solver) Solve(state *shogi.State) *Node {
	root := &Node{
		pn: inf - 1,
		dn: inf - 1,
		move: &shogi.Move{
			Turn: shogi.TurnWhite,
		},
		state: state,
	}
	// s.root = root
	s.mid(root)
	if root.getPhi() != inf && root.getDelta() != inf {
		root.setPhi(inf)
		root.setDelta(inf)
		s.mid(root)
	}
	return root
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
		for _, ms := range candidates(n.state, !n.move.Turn) {
			n.Children = append(n.Children, &Node{
				pn: 1, dn: 1,
				move:   ms.move,
				state:  ms.state,
				parent: n,
			})
		}
	}
	if len(n.Children) == 0 {
		switch n.move.Turn {
		case shogi.TurnBlack:
			n.setResult(resultT)
		case shogi.TurnWhite:
			n.setResult(resultF)
		}
		s.putInHash(n)
		return
	}
	// 3. put in hash
	s.putInHash(n)
	// 4. multiple-iterative deepening
	for {
		p, d := n.getPhi(), n.getDelta()
		md, sp := s.minDelta(n), s.sumPhi(n)
		if p <= md || d <= sp {
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
			if d2 != inf {
				d2++
			}
			min := d2
			if p < min {
				min = p
			}
			c.setDelta(min)
		}
		s.mid(c)
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
		if hash.pn == inf {
			log.Printf("%v %v %v", best, h, delta2)
			return best, h, delta2
		}
	}
	return best, h, delta2
}

func (s *Solver) putInHash(n *Node) {
	key := n.state.Hash()
	s.hash[key] = &hash{
		pn: n.getPhi(),
		dn: n.getDelta(),
	}
}

func (s *Solver) lookUpHash(n *Node) *hash {
	key := n.state.Hash()
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
