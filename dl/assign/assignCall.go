package assign

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *CallPattern) Execute(run rt.Runtime) error {
	_, err := op.determine(run, affine.None)
	return err
}

func (op *CallPattern) GetBool(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallPattern) GetNumber(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Number)
}

func (op *CallPattern) GetText(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Text)
}

func (op *CallPattern) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.Record)
}

func (op *CallPattern) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.NumList)
}

func (op *CallPattern) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.TextList)
}

func (op *CallPattern) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.determine(run, affine.RecordList)
}

// note: at one point this would unwrap errors so that callers couldn't see them
// i no longer am sure why. doing stops game.Signals(s) ( ex SignalQuit ) from reaching the parser.
func (op *CallPattern) determine(run rt.Runtime, aff affine.Affinity) (ret rt.Value, err error) {
	name := inflect.Normalize(op.PatternName)
	if k, v, e := ExpandArgs(run, op.Arguments); e != nil {
		err = CmdErrorCtx(op, name, e)
	} else if v, e := run.Call(name, aff, k, v); e != nil {
		err = CmdErrorCtx(op, name, e)
	} else {
		ret = v
	}
	return
}

func ExpandArgs(run rt.Runtime, args []Arg) (retKeys []string, retVals []rt.Value, err error) {
	if len(args) > 0 {
		keys, vals := make([]string, 0, len(args)), make([]rt.Value, len(args))
		for i, a := range args {
			if val, e := safe.GetAssignment(run, a.Value); e != nil {
				err = fmt.Errorf("%w while reading arg %d(%s)", e, i, a.Name)
				break
			} else if n := inflect.Normalize(a.Name); len(n) > 0 {
				keys = append(keys, n)
				vals[i] = val
			} else if len(keys) > 0 {
				err = fmt.Errorf("only named arguments can follow named arguments %d", i)
			} else {
				vals[i] = val
			}
		}
		if err == nil {
			retKeys, retVals = keys, vals
		}
	}
	return
}
