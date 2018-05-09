package automata

type StateTransition struct {
	maxDistance int
	queryLength int
}

func CreateStateTransition(maxDistance int, queryLength int) StateTransition {
	return StateTransition{
		maxDistance: maxDistance,
		queryLength: queryLength,
	}
}

func (st StateTransition) subsumes(lhs *Position, rhs *Position) bool {
	i := lhs.TermIndex
	e := lhs.ErrorCount
	j := rhs.TermIndex
	f := rhs.ErrorCount
	if i < j {
		return j-i <= f-e
	} else {
		return i-j <= f-e
	}
}

func (st StateTransition) unsubsume(state *State) {
	outerIter := state.Iterator()

	for outerIter.HasNext() {
		outer := outerIter.Next()
		outerErrors := outer.ErrorCount
		innerIter := outerIter.Copy()

		for innerIter.HasNext() {
			inner := innerIter.Peek()
			if outerErrors < inner.ErrorCount {
				break
			}

			innerIter.Next()
		}

		for innerIter.HasNext() {
			inner := innerIter.Next()
			if st.subsumes(outer, inner) {
				innerIter.Remove()
			}
		}
	}
}

func (st StateTransition) indexOf(characteristicVector []bool, k int, i int) int {
	for j := 0; j < k; j++ {
		if characteristicVector[i+j] {
			return j
		}
	}

	return -1
}

func (st StateTransition) transition(pos *Position, characteristicVector []bool, offset int) *State {
	n := st.maxDistance
	i := pos.TermIndex
	e := pos.ErrorCount
	h := i - offset
	w := len(characteristicVector)

	if e < n {
		// Consider any character before the last one of the spelling candidate
		if h <= w-2 {
			a := n - e + 1
			b := w - h
			k := b

			if a < b {
				k = a
			}

			j := st.indexOf(characteristicVector, k, h)

			if j == 0 {
				// [No Error]: Increment the index by one; leave the error alone.
				return CreateState([]*Position{CreatePosition(1+i, e)})
			} else if j > 0 {
				// [Insertion]: Leave the index alone; increment the error by one.
				// [Substitution]: Increment both the index and error by one.
				// [Deletion]: Increment the index by one-more than the number of
				// deletions; increment the error by the number of deletions.
				return CreateState([]*Position{
					CreatePosition(i, e+1),
					CreatePosition(i+1, e+1),
					CreatePosition(i+j+1, e+j)})
			}
			//else, j < 0
			// [Insertion]: Leave the index alone; increment the error by one.
			// [Substitution]: Increment both the index and error by one.
			return CreateState([]*Position{
				CreatePosition(i, e+1),
				CreatePosition(i+1, e+1)})
		}
		// Consider the last character of the spelling candidate
		if h == w-1 {
			if characteristicVector[h] {
				// [No Error]: Increment the index by one; leave the error alone.
				return CreateState([]*Position{CreatePosition(i+1, e)})
			}

			// [Insertion]: Leave the index alone; increment the error by one.
			// [Substitution]: Increment both the index and error by one.
			return CreateState([]*Position{
				CreatePosition(i, e+1),
				CreatePosition(i+1, e+1)})

		}
		// else, h == w
		// [Insertion]: Leave the index alone; increment the error by one.		
		return CreateState([]*Position{CreatePosition(i, e+1)})
	}
	// The edit distance is at its maximum, allowed value.  Only consider this
	// spelling candidate if there is no error at the index of its current term.
	if e == n && h <= w-1 && characteristicVector[h] {
		// [No Error]: Increment the index by one; leave the error alone.
		return CreateState([]*Position{CreatePosition(i+1, n)})
	}

	// [Too Many Errors]: The edit distance has exceeded the max distance for
	// the candidate term.
	return nil
}

func (st StateTransition) Next(currState *State, characteristicVector []bool) *State {
	offset := currState.Head.TermIndex
	nextState := CreateState([]*Position{})

	itr := currState.Iterator()

	for itr.HasNext() {
		position := itr.Next()
		positions := st.transition(position, characteristicVector, offset)

		if positions == nil {
			continue
		}

		nextState.mergePositions(positions)
	}

	st.unsubsume(nextState)

	if nextState.Head != nil {
		nextState.Sort()
		return nextState
	}

	return nil
}
