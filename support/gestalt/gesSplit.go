package gestalt

// tries each branch; the result is the union of all separate successes.
type Branch struct {
	Options []Interpreter
}

func (op *Branch) Match(q Query, cs []InputState) (ret []InputState) {
	for _, m := range op.Options {
		if next := m.Match(q, cs); len(next) > 0 {
			ret = append(ret, next...)
		}
	}
	return
}

// matches each of its sub-expressions in turn
// filter the results with each new
type Sequence struct {
	Series []Interpreter
}

func (op *Sequence) Match(q Query, cs []InputState) []InputState {
	return runSequence(q, cs, op.Series)
}
func runSequence(q Query, cs []InputState, series []Interpreter) []InputState {
	next := cs
	for _, m := range series {
		if next = m.Match(q, next); len(next) == 0 {
			break
		}
	}
	return next
}

// target a specific group in the output
// run the group as a sequence.
type Target struct {
	Primary bool
	Group   []Interpreter
}

func (op *Target) Match(q Query, cs []InputState) (ret []InputState) {
	for i := range cs {
		// needs to modify the elements in the list
		cs[i].AddResult(MatchedTarget(op.Primary))
	}
	return runSequence(q, cs, op.Group)
}

type MatchedTarget bool

func (m MatchedTarget) String() string { return "target" }
func (m MatchedTarget) NumWords() int  { return 0 }
