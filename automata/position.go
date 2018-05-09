package automata

type Position struct {
	TermIndex  int
	ErrorCount int
	Next       *Position
}

func CreatePosition(termIndex int, errorCount int) *Position {
	return &Position{
		TermIndex:  termIndex,
		ErrorCount: errorCount,
		Next:       nil,
	}
}
