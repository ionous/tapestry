package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type Sequence struct {
	Seq   string `if:"internal"` // generated at import time to provide a unique counter for each sequence
	Parts []rt.TextEval
}

// CycleText when called multiple times returns each of its elements in turn.
type CycleText struct {
	Sequence
}

// ShuffleText walks randomly through the list of elements without repeating any until they all have been shown.
// Unlike the other sequences, ShuffleText is not persistent.
type ShuffleText struct {
	Sequence
	indices []int
}

// StoppingText options returns its values in-order, once per evaluation, until the options are exhausted, then the last one is kept.
// As a special case, if there is only one option: it gets returned once, followed by the empty string forever after.
type StoppingText struct {
	Sequence
}

// sequences store fields into autogenerated targets
// if persistence wasn't necessary, we could use in-memory private fields of the commands
func (op *Sequence) updateCounter(run rt.Runtime, inc func(int, int) int) (ret int, err error) {
	if max := len(op.Parts); max > 0 {
		if p, e := run.GetField(object.Counter, op.Seq); e != nil {
			err = e
		} else if num, e := p.GetNumber(run); e != nil {
			err = e
		} else {
			curr := int(num)
			next := &generic.Int{Value: inc(curr, max)}
			if e := run.SetField(object.Counter, op.Seq, next); e != nil {
				err = e
			} else {
				ret = curr
			}
		}
	}
	return
}

func (*CycleText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "cycle_text",
		Group: "cycle",
		Desc:  "Cycle Text: When called multiple times, returns each of its inputs in turn.",
	}
}

func (op *CycleText) GetText(run rt.Runtime) (ret string, err error) {
	if curr, e := op.updateCounter(run, wrap); e != nil {
		err = e
	} else if max := len(op.Parts); curr < max {
		ret, err = rt.GetText(run, op.Parts[curr])
	}
	return
}

func (*ShuffleText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "shuffle_text",
		Group: "format",
		Desc:  "Shuffle Text: When called multiple times returns its inputs at random.",
	}
}

func (op *ShuffleText) GetText(run rt.Runtime) (ret string, err error) {
	if curr, e := op.updateCounter(run, wrap); e != nil {
		err = e
	} else if max := len(op.Parts); curr < max {
		// uses the Fisher–Yates algorithm to sort indices
		if len(op.indices) == 0 {
			op.indices = makeIndices(max)
		}
		// we pick anything b/t the start and end of our list
		j := run.Random(curr, max)
		if curr != j { // switch if they are different locations.
			op.indices[curr], op.indices[j] = op.indices[j], op.indices[curr]
		}
		// and that's the real index we use.
		sel := op.indices[curr]
		ret, err = rt.GetText(run, op.Parts[sel])
	}
	return
}

func (*StoppingText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "stopping_text",
		Group: "format",
		Desc:  "Stopping Text: When called multiple times returns each of its inputs in turn, sticking to the last one.",
	}
}

func (op *StoppingText) GetText(run rt.Runtime) (ret string, err error) {
	switch max := len(op.Parts); max {
	case 0:
		// do nothing.
	case 1:
		// when one element, return it once then the empty string after.
		if curr, e := op.updateCounter(run, saturate); e != nil {
			err = e
		} else if curr == 0 {
			ret, err = rt.GetText(run, op.Parts[curr])
		}
	default:
		if curr, e := op.updateCounter(run, cap); e != nil {
			err = e
		} else if curr < max {
			ret, err = rt.GetText(run, op.Parts[curr])
		}
	}
	return
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
