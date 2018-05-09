package automata

type StateIterator struct {
	state     *State
	lookAhead *Position
	curr      *Position
	prev      *Position
	next      *Position
}

func (si *StateIterator) Peek() *Position {
	si.advance()
	return si.next
}

func (si *StateIterator) HasNext() bool {
	si.advance()
	return si.next != nil
}

func (si *StateIterator) Next() *Position {
	si.advance()
	nextLocal := si.next
	si.next = nil
	return nextLocal
}

func (si *StateIterator) Copy() *StateIterator {
	copy := CreateStateIterator(si.state, si.lookAhead)
	copy.next = si.next
	return copy
}

func (si *StateIterator) Remove() {
	if si.curr != nil {
		si.state.Remove(si.prev, si.curr)
		si.curr = nil
	}
}

func (si *StateIterator) advance() {
	if si.next == nil && si.lookAhead != nil {
		si.next = si.lookAhead
		if si.curr != nil {
			si.prev = si.curr
		}
		si.curr = si.next
		si.lookAhead = si.lookAhead.Next
	}
}

func CreateStateIterator(state *State, head *Position) *StateIterator {
	return &StateIterator{
		state:     state,
		lookAhead: head,
	}
}
