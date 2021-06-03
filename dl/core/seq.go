package core

import (
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// sequences store fields into autogenerated targets
// if persistence wasn't necessary, we could use in-memory private fields of the commands
func updateCounter(run rt.Runtime, seq string, parts []rt.TextEval, inc func(int, int) int) (ret int, err error) {
	if max := len(parts); max > 0 {
		if p, e := run.GetField(object.Counter, seq); e != nil {
			err = e
		} else {
			curr := p.Int()
			next := g.IntOf(inc(curr, max))
			if e := run.SetField(object.Counter, seq, next); e != nil {
				err = e
			} else {
				ret = curr + 1
			}
		}
	}
	return
}

func getNextText(run rt.Runtime, parts []rt.TextEval, onedex int) (ret g.Value, err error) {
	if i, max := onedex-1, len(parts); i >= 0 && i < max {
		ret, err = safe.GetText(run, parts[i])
	} else {
		ret = g.Empty
	}
	return
}

func (op *CallCycle) GetText(run rt.Runtime) (ret g.Value, err error) {
	if onedex, e := updateCounter(run, op.At.String(), op.Parts, wrap); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = getNextText(run, op.Parts, onedex)
	}
	return
}

func (op *CallShuffle) GetText(run rt.Runtime) (ret g.Value, err error) {
	if curr, e := updateCounter(run, op.At.String(), op.Parts, wrap); e != nil {
		err = cmdError(op, e)
	} else if curr, max := curr-1, len(op.Parts); curr < max {
		onedex := op.Indices.shuffle(run, curr, max)
		ret, err = getNextText(run, op.Parts, onedex)
	}
	return
}

func (op *CallTerminal) GetText(run rt.Runtime) (ret g.Value, err error) {
	if onedex, e := op.stopping(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret, err = getNextText(run, op.Parts, onedex)
	}
	return
}

func (op *CallTerminal) stopping(run rt.Runtime) (ret int, err error) {
	switch max := len(op.Parts); max {
	case 0:
		// no elements, nothing to do.
	case 1:
		// when one element, return it once then the empty string after.
		ret, err = updateCounter(run, op.At.String(), op.Parts, saturate)
	default:
		ret, err = updateCounter(run, op.At.String(), op.Parts, cap)
	}
	return
}

type Shuffler struct {
	indices []int
}

func (x *Shuffler) shuffle(run rt.Runtime, curr, max int) int {
	// uses the Fisher–Yates algorithm to sort indices
	if len(x.indices) == 0 {
		x.indices = makeIndices(max)
	}
	// we pick anything b/t the start and end of our list
	j := run.Random(curr, max)
	if curr != j { // switch if they are different locations.
		x.indices[curr], x.indices[j] = x.indices[j], x.indices[curr]
	}
	// and that's the real index we use.
	return 1 + x.indices[curr]
}

func makeIndices(max int) []int {
	indices := make([]int, max)
	for i := 0; i < max; i++ {
		indices[i] = i
	}
	return indices
}

func wrap(curr, max int) int {
	return (curr + 1) % max
}

func cap(curr, max int) int {
	if next := curr + 1; next < max {
		curr = next
	}
	return curr
}

func saturate(curr, max int) int {
	if curr < max {
		curr++
	}
	return curr
}
