package automata

import (
	"log"
	"strings"
)

type Transition struct {
	Source *Node
	Char   rune
	Target *Node
}

type Dawg struct {
	Root                 *Node
	previousTerm         []rune
	uncheckedTransitions []Transition
	minimizedNodes       map[string]*Node
}

func CreateDawg() *Dawg {
	return &Dawg{
		previousTerm:         []rune(""),
		Root:                 CreateNode([]rune{}, false),
		uncheckedTransitions: []Transition{},
		minimizedNodes:       map[string]*Node{},
	}
}

func (d *Dawg) minimize(lowerBound int) {
	for j := len(d.uncheckedTransitions) - 1; j >= lowerBound; j-- {
		t := d.uncheckedTransitions[j]

		if _, e := d.minimizedNodes[string(t.Target.Label)]; !e {
			d.minimizedNodes[string(t.Target.Label)] = t.Target
		}
		
		t.Source.AddEdge(t.Char, d.minimizedNodes[string(t.Target.Label)])
	}

	d.uncheckedTransitions = d.uncheckedTransitions[0:lowerBound]
}

func (d *Dawg) Finish() {
	d.minimize(0)
}

func (d *Dawg) AddPattern(term []rune) {
	if strings.Compare(string(term), string(d.previousTerm)) < 0 {
		log.Panic("Must be sorted")
	}

	upper := len(term)
	if len(d.previousTerm) < len(term) {
		upper = len(d.previousTerm)
	}

	i := 0
	for i < upper && term[i] == d.previousTerm[i] {
		i += 1
	}

	d.minimize(i)

	node := d.Root

	if len(d.uncheckedTransitions) > 0 {
		last := len(d.uncheckedTransitions) - 1
		node = d.uncheckedTransitions[last].Target
	}

	k := len(term) - 1

	for i <= k {
		char := rune(term[i])
		nextNode := CreateNode(term[:i+1], i == k)
		transition := Transition{Source: node, Char: char, Target: nextNode}
		d.uncheckedTransitions = append(d.uncheckedTransitions, transition)
		node = nextNode
		i += 1
	}

	d.previousTerm = term
}
