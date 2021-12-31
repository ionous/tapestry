package core

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

type Parameterizer interface{ Pack() []rt.Arg }

func (op *CallArgs) Pack() (args []rt.Arg) {
	// FIX THIS COPY!
	for _, a := range op.Args {
		args = append(args, rt.Arg{a.Name, a.From})
	}
	return
}

func (op *CallPattern) Execute(run rt.Runtime) error {
	_, err := op.determine(run, "")
	return err
}

// GetNumber returns the first matching num evaluation.
func (op *CallPattern) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

// GetText returns the first matching text evaluation.
func (op *CallPattern) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
}

// GetBool returns the first matching bool evaluation.
func (op *CallPattern) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallPattern) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *CallPattern) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *CallPattern) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *CallPattern) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.RecordList)
}

// backdoor for RenderPattern
// could maybe be GetAssignedValue but that would expose it to the composer that way too
func (op *CallPattern) DetermineValue(run rt.Runtime) (ret g.Value, err error) {
	return op.determine(run, "")
}

func (op *CallPattern) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := run.Call(op.Pattern.String(), aff, op.Arguments.Pack()); e != nil && !errors.Is(e, rt.NoResult{}) {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
