package automata

type Automata struct {
	Dawg *Dawg
}

func (a *Automata) distance(state *State, queryLength int) int {
	minDistance := queryLength

	itr := state.Iterator()
	for itr.HasNext() {
		position := itr.Next()
		i := position.TermIndex
		e := position.ErrorCount
		distance := queryLength - i + e
		if distance < minDistance {
			minDistance = distance
		}
	}

	return minDistance
}

func (auto *Automata) characteristicVector(x rune, term string, k int, i int) []bool {
	characteristicVector := make([]bool, k)

	for j := 0; j < k; j++ {
		characteristicVector[j] = x == rune(term[i+j])
	}

	return characteristicVector
}

type Result struct {
	Match string
	Error int
}


func (auto *Automata) FindAll(term string, maxDistance int) []Result {
	var labels []rune = []rune{}
	var label rune
	var intersection *Intersection

	i := 0
	k := 0

	result := []Result{}

	stateTransition := CreateStateTransition(maxDistance, len(term))
	a := (maxDistance << 1) + 1

	initialState := CreateState([]*Position{CreatePosition(0, 0)})
	pendingQueue := []*Intersection{CreateIntersection(auto.Dawg.Root, initialState, nil, 0)}
	
	for len(labels) > 0 || len(pendingQueue) > 0 {
		if len(labels) > 0 {
			dictionaryNode := intersection.Node
			levenshteinState := intersection.State
			label, labels = labels[0], labels[1:]

			nextDictionaryNode := dictionaryNode.Transition(label)
			characteristicVector := auto.characteristicVector(label, term, k, i)
			nextLevenshteinState := stateTransition.Next(levenshteinState, characteristicVector)
			if nextLevenshteinState != nil {
				nextIntersection := CreateIntersection(nextDictionaryNode, nextLevenshteinState, intersection, label)
				pendingQueue = append(pendingQueue, nextIntersection)
				if nextDictionaryNode.Final {
					distance := auto.distance(nextLevenshteinState, len(term))
					if distance <= maxDistance {
						result = append(result, Result{Match: nextIntersection.Candidate(), Error: distance})
					}

					if distance == 0 {
						break
					}
				}
			}
		} else {
			intersection, pendingQueue = pendingQueue[0], pendingQueue[1:]
			dictionaryNode := intersection.Node
			levenshteinState := intersection.State

			i = levenshteinState.Head.TermIndex
			b := len(term) - i
			k = b
			if a < b {
				k = a
			}
			labels = dictionaryNode.Labels
		}
	}

	return result
}

func CreateAutomata(patterns []string) *Automata {
	auto := &Automata{
		Dawg: CreateDawg(),
	}

	for _, pattern := range patterns {
		auto.Dawg.AddPattern(pattern)
	}

	auto.Dawg.Finish()

	return auto
}
