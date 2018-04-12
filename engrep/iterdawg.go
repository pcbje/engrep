package engrep

import (
	"log"
	"strings"
)

type Reference struct {
	Node           *Node
	RemainingError int
}

type Transition struct {
	Source *Node
	Char   rune
	Target *Node
}

type Dawg struct {
	index                int
	Root                 *Node
	previousTerm         string
	uncheckedTransitions []Transition
	minimizedNodes       map[string]*Node
	K                    int
}

func CreateDawg(k int) *Dawg {
	return &Dawg{
		K:                    k,
		index:                1,
		previousTerm:         "",
		Root:                 CreateNode(0, []rune{}, k, true, false, 0, k),
		uncheckedTransitions: []Transition{},
		minimizedNodes:       map[string]*Node{},
	}
}

func (d *Dawg) Iterator() *Node {
	d.Commit(0)
	itr := CreateNode(0, []rune{}, d.K, false, false, 0, d.K)
	itr.AddReference(Reference{Node: d.Root, RemainingError: d.K})
	itr.Explore()
	return itr
}

func (d *Dawg) GetRoot() *Node {
	return d.Root
}

func (d *Dawg) Commit(lower int) {
	for j := len(d.uncheckedTransitions) - 1; j >= lower; j-- {
		tr := d.uncheckedTransitions[j]
		node := tr.Target
		hash := node.GetHash()

		if _, ok := d.minimizedNodes[hash]; ok {
			tr.Source.AddEdge(tr.Char, d.minimizedNodes[hash])
		} else {
			d.minimizedNodes[hash] = node
		}
	}

	d.uncheckedTransitions = d.uncheckedTransitions[:lower]
}

func (d *Dawg) AddPattern(term string) {
	if strings.Compare(term, d.previousTerm) < 0 {
		log.Panic("Must be sorted")
	}

	upper := len(term)
	if len(d.previousTerm) < len(term) {
		upper = len(d.previousTerm)
	}

	node := d.Root
	runes := []rune(term)

	i := 0

	for i < upper && term[i] == d.previousTerm[i] {
		node = d.uncheckedTransitions[i].Target
		suffixLength := len(node.Suffix)
		if suffixLength > 0 {
			char, next := node.Split(d.index)
			transition := Transition{Source: node, Char: char, Target: next}
			d.uncheckedTransitions = append(d.uncheckedTransitions, transition)
			d.index += suffixLength + 1
		}

		i += 1
	}

	d.Commit(i)

  if i == len(runes) {
    return
  }

	char, suffix := runes[i], runes[i+1:]
	nextNode := CreateNode(d.index, suffix, d.K, true, len(suffix) == 0, char, d.K)
	transition := Transition{Source: node, Char: char, Target: nextNode}
	d.uncheckedTransitions = append(d.uncheckedTransitions, transition)
	node.AddEdge(char, nextNode)
	d.index += len(suffix) + 1

	d.previousTerm = term
}
