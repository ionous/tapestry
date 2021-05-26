package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"

	g "git.sr.ht/~ionous/iffy/rt/generic"
)

func (op *Determine) Execute(run rt.Runtime) error {
	_, err := op.determine(run, "")
	return err
}

// GetNumber returns the first matching num evaluation.
func (op *Determine) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

// GetText returns the first matching text evaluation.
func (op *Determine) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
}

// GetBool returns the first matching bool evaluation.
func (op *Determine) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *Determine) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *Determine) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *Determine) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *Determine) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.RecordList)
}

// backdoor for RenderPattern
// could maybe be GetAssignedValue but that would expose it to the composer that way too
func (op *Determine) DetermineValue(run rt.Runtime) (ret g.Value, err error) {
	return op.determine(run, "")
}

func (op *Determine) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	var args []rt.Arg
	// FIX THIS COPY!
	for _, a := range op.Arguments.Args {
		args = append(args, rt.Arg{a.Name.Value(), a.From})
	}
	name := op.Pattern.Value()
	if v, e := run.Call(name, aff, args); e != nil && !errors.Is(e, rt.NoResult{}) {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
