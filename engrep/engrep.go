package engrep

import (
	"sort"
)

type State struct {
	Start   int
	Node    *Node
	Deletes int
	Inserts int
	Depth int
}

type Engrep struct {
	dawg      *Dawg
	root      *Node
	k         int
	backtrack bool
}

func CreateEngrep(k int, backtrack bool, dawg *Dawg) *Engrep {
	return &Engrep{
		dawg:      dawg,
		k:         k,
		backtrack: backtrack,
	}
}

func (t *Engrep) AddReferences(references []string) {
	sort.Strings(references)

	var prev string
	for _, reference := range references {
		if reference != prev {
			t.dawg.AddPattern(reference)
			prev = reference
		}
	}

	t.dawg.Commit(0)
	t.root = t.dawg.Iterator()
}

func (t *Engrep) Scan(text string, k int, callback func(int, int, string, string, string, int))  {
	var states [100000]State = [100000]State{}
	var up bool = true
	var counter int = 0

	if k > t.k {
		k = t.k
	}

	for offset, char := range []rune(text) {
		up = !up
		nx := 0
		counts := 0

		new := 0
		nxdir := 1
		if !up {
			nx = len(states) - 1
			nxdir = -1
		}

		for i := 0; i < counter; i++ {
			ii := i

			if up {
				ii = len(states) - i - 1
			}

			state := states[ii]
			node := state.Node.Transition(char)

			if node != nil && state.Inserts+node.Cost <= k {
				states[nx].Node = node
				states[nx].Deletes = state.Deletes
				states[nx].Inserts = state.Inserts + node.Cost
				states[nx].Start = state.Start
				states[nx].Depth = state.Depth+1

				nx += nxdir
				counts++
				new++

				if node.Final && (state.Deletes <= node.Remaining || state.Inserts <= node.Remaining) {
					if t.backtrack {
						actual := text[state.Start: offset+2]
						pre := string(text[state.Start-1])
						suf := string(text[offset+2])
						callback(state.Start, offset, actual, pre, suf, state.Deletes+state.Inserts)
					} else {
						callback(state.Start, offset, "", "", "", state.Deletes+state.Inserts)
					}
				}
			}

			if state.Deletes+1 <= k {
				states[nx].Node = state.Node
				states[nx].Deletes = state.Deletes + 1
				states[nx].Inserts = state.Inserts
				states[nx].Start = state.Start
				states[nx].Depth = state.Depth

				nx += nxdir
				counts++
			}
		}

		node := t.root.Transition(char)

		if node != nil {
			states[nx].Node = node
			states[nx].Inserts = node.Cost
			states[nx].Deletes = node.Cost
			states[nx].Start = offset
			states[nx].Depth = 1
			counts++
		}

		counter = counts
	}
}
