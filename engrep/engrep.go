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

func (t *Engrep) Scan(text string, maxPatternLength int, callback func(int, int, string, string, int)) ([]int, map[int][]int, []int) {
	var states [100000]State = [100000]State{}
	var up bool = true
	var counter int = 0

  var count []int = make([]int, len(text))
	var hits map[int][]int = map[int][]int{}
	hits[0] = []int{0, 0}

	var errors []int = make([]int, 3)

	for offset, char := range []rune(text) {
		up = !up
		nx := 0
		counts := 0

		out := 0
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

			if _, ok := hits[state.Depth]; !ok {
				hits[state.Depth] = []int{0, 0}
			}

			if node != nil && state.Inserts+node.Cost <= t.k {
				hits[state.Depth][0]++

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
						actual := text[state.Start : offset+1]
						for _, reference := range node.Backtrack() {
							callback(state.Start, offset, reference, actual, state.Deletes+state.Inserts)
						}
					} else {
						callback(state.Start, offset, "", "", state.Deletes+state.Inserts)
					}
				}
			}

			if node == nil {
				hits[state.Depth][1]++
			}

			if state.Deletes+1 <= t.k {
				states[nx].Node = state.Node
				states[nx].Deletes = state.Deletes + 1
				states[nx].Inserts = state.Inserts
				states[nx].Start = state.Start
				states[nx].Depth = state.Depth

				nx += nxdir
				counts++
			} else {
				out++
			}
		}

		node := t.root.Transition(char)

		if node != nil {
			hits[0][0]++
			states[nx].Node = node
			states[nx].Inserts = node.Cost
			states[nx].Deletes = node.Cost
			states[nx].Start = offset
			states[nx].Depth = 1
			counts++
		} else {
			hits[0][1]++
		}

		count[offset] = counts
		errors[0] += new
		errors[1] += counts
		errors[2] += out

		counter = counts
	}

	return count, hits, errors
}
