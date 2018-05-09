package automata

type Intersection struct {
	prevIntersection *Intersection
	Label            rune
	Node             *Node
	State            *State
}

func (i *Intersection) Candidate() string {
	if i.prevIntersection != nil {
		return i.prevIntersection.Candidate() + string(i.Label)
	} else {
		return ""
	}
}

func CreateIntersection(node *Node, state *State, prev *Intersection, label rune) *Intersection {
	return &Intersection{
		prevIntersection: prev,
		Node:             node,
		State:            state,
		Label:            label,
	}
}
