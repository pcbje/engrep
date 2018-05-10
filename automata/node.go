package automata

type Node struct {
	Label  []rune
	Final  bool
	Edges  map[rune]*Node
	Labels []rune
}

func (n *Node) Transition(char rune) *Node {
	return n.Edges[char]
}

func (n *Node) AddEdge(char rune, target *Node) {
	if _, e := n.Edges[char]; !e {
		n.Labels = append(n.Labels, char)
	}

	n.Edges[char] = target
}

func CreateNode(label []rune, final bool) *Node {
	return &Node{
		Label:  label,
		Final:  final,
		Edges:  map[rune]*Node{},
		Labels: []rune{},
	}
}
