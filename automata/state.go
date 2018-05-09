package automata

type State struct {
	Head *Position
}

func (s *State) SetHead(head *Position) {
	head.Next = s.Head
	s.Head = head
}

func (s *State) Add(next *Position) {
	if s.Head == nil {
		s.Head = next
	} else {
		curr := s.Head
		prev := curr
		for curr != nil {
			prev = curr
			curr = curr.Next
		}

		prev.Next = next
	}
}

func (s *State) insertAfter(curr *Position, next *Position) {
	if curr != nil {
		next.Next = curr.Next
		curr.Next = next
	} else {
		s.Add(next)
	}
}

func (s *State) Insert(curr *Position, next *Position) {
	if curr != nil {
		s.insertAfter(curr, next)
	} else {
		s.SetHead(next)
	}
}

func (s *State) mergePositions(other *State) {
	itr := other.Iterator()
	for itr.HasNext() {
		a := itr.Next()
		i := a.TermIndex
		e := a.ErrorCount

		iter := s.Iterator()
		var prevB *Position = nil
		for iter.HasNext() {
			b := iter.Peek()
			j := b.TermIndex
			f := b.ErrorCount

			if e < f || (e == f && i < j) {
				prevB = b
				iter.Next()
			} else {
				break
			}
		}

		if iter.HasNext() {
			b := iter.Peek()
			j := b.TermIndex
			f := b.ErrorCount

			if j != i || f != e {
				s.Insert(prevB, a)
			}
		} else {
			s.Insert(prevB, a)
		}
	}
}

func (s *State) Remove(prev *Position, curr *Position) {
	if prev != nil {
		prev.Next = curr.Next
	} else {
		s.Head = s.Head.Next
	}
}

func (s *State) Iterator() *StateIterator {
	return CreateStateIterator(s, s.Head)
}

func (s *State) compare(lhs *Position, rhs *Position) int {
	c := lhs.TermIndex - rhs.TermIndex

	if c != 0 {
		return c
	}

	return lhs.ErrorCount - rhs.ErrorCount
}

func (s *State) Sort() {
	s.Head = s.mergeSort(s.Head)
}

func (s *State) mergeSort(lhsHead *Position) *Position {
	if lhsHead == nil || lhsHead.Next == nil {
		return lhsHead
	}

	middle := s.middle(lhsHead)
	rhsHead := middle.Next
	middle.Next = nil

	return s.merge(s.mergeSort(lhsHead), s.mergeSort(rhsHead))
}

func (s *State) middle(head *Position) *Position {
	slow := head
	fast := head

	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	return slow
}

func (s *State) merge(lhsHead *Position, rhsHead *Position) *Position {
	next := CreatePosition(-1, -1)
	curr := next

	for lhsHead != nil && rhsHead != nil {
		if s.compare(lhsHead, rhsHead) <= 0 {
			curr.Next = lhsHead
			lhsHead = lhsHead.Next
		} else {
			curr.Next = rhsHead
			rhsHead = rhsHead.Next
		}

		curr = curr.Next
	}

	if rhsHead != nil {
		curr.Next = rhsHead
	} else if lhsHead != nil {
		curr.Next = lhsHead
	}

	curr = next.Next
	return curr
}

func CreateState(positions []*Position) *State {
	state := &State{}

	var prev *Position = nil

	for _, curr := range positions {
		state.insertAfter(prev, curr)
		prev = curr
	}

	return state
}
