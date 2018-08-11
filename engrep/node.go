package engrep

import "strconv"

type Node struct {
	ID         int
	K          int
	Final      bool
	frontier   bool
	Edges      map[rune]*Node
	References map[int]Reference
	Suffix     []rune
	char       rune
	Cost       int
	Remaining  int
}

func (n *Node) IsFinal() bool {
	return n.Final
}

func (n *Node) GetHash() string {
	hash := ""

	if n.Final {
		hash += "f"
	}

	for char, target := range n.Edges {
		hash += string(char)
		hash += ":"
		hash += strconv.Itoa(target.ID)
		hash += ","
	}

	hash += string(n.Suffix)

	return hash
}

var index = 0

func (n *Node) ExploreRec(remaining int, node *Node, round int) {
	node.Split(node.ID + 1)

	for char, target := range node.Edges {
		n.AddEdgeRef(char, round, Reference{Node: target, RemainingError: remaining - round})
		index++
		if remaining > round {
			n.ExploreRec(remaining, target, round+1)
		}
	}
}

func (n *Node) Split(i int) (rune, *Node) {
	if len(n.Suffix) == 0 {
		// Already splitted?
		return 0, nil
	}
	char, newSuffix := n.Suffix[0], n.Suffix[1:]

	next := CreateNode(i, newSuffix, n.K, true, len(newSuffix) <= n.K, char, n.K)
	next.Edges = n.Edges

	n.Suffix = []rune{}
	n.Edges = map[rune]*Node{char: next}

	return char, next
}

func (n *Node) Explore() {
	for _, reference := range n.References {
		n.ExploreRec(reference.RemainingError, reference.Node, 0)
	}

	//n.References = nil
}

func (n *Node) Transition(char rune) *Node {
	if n.frontier {
		n.Explore()
		n.frontier = false
	}

	return n.Edges[char]
}

func (n *Node) AddEdge(char rune, target *Node) {
	n.frontier = false

	n.Edges[char] = target
}

func (n *Node) SetFinal() {
	n.Final = true
}

func (n *Node) IsFrontier() bool {
	return n.frontier
}

func (n *Node) SetFrontier(f bool) {
	n.frontier = f
}

func (n *Node) AddReference(ref Reference) {
	if ref.Node != nil {
		if n.References[ref.Node.ID].Node == nil || n.References[ref.Node.ID].RemainingError < ref.RemainingError {
			n.References[ref.Node.ID] = ref
		}
	}
}

func (n *Node) AddEdgeRef(char rune, cost int, ref Reference) {
	if n.Edges[char] == nil {
		n.Edges[char] = CreateNode(ref.Node.ID, []rune{}, ref.Node.K, true, ref.Node.Final, char, ref.RemainingError)
		n.Edges[char].Cost = cost
	}

	if cost < n.Edges[char].Cost {
		n.Edges[char].Cost = cost
	}

	n.Edges[char].AddReference(ref)
}

func CreateNode(id int, suffix []rune, k int, frontier bool, final bool, char rune, remaining int) *Node {
	return &Node{
		ID:         id,
		K:          k,
		Edges:      map[rune]*Node{},
		References: map[int]Reference{},
		Final:      final,
		frontier:   frontier,
		Suffix:     suffix,
		char:       char,
		Remaining:  k - remaining,
	}
}
