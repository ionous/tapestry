package assign

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func (op *CallPattern) Execute(run rt.Runtime) error {
	_, err := op.determine(run, affine.None)
	return err
}

func (op *CallPattern) GetBool(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Bool)
}

func (op *CallPattern) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Number)
}

func (op *CallPattern) GetText(run rt.Runtime) (g.Value, error) {
	return op.determine(run, affine.Text)
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

func (op *CallPattern) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	pat := op.PatternName
	if rec, e := MakeRecord(run, pat, op.Arguments...); e != nil {
		err = cmdError(op, e)
	} else if k := rec.Kind(); !k.Implements(kindsOf.Record.String()) {
		// note: this doesnt positively affirm kindsOf.Pattern:
		// some tests use golang structs as faux patterns.
		// ( instead the internals verify there are labels )
		if v, e := run.Call(rec, aff); e != nil && !errors.Is(e, rt.NoResult) {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	} else {
		// fix? CallPattern is being used as a side effect to initialize records
		// tbd: at the very least maybe use error return checking instead of testing kind record hierarchy?
		// ( originally this was inside of Runner.Call; qnaCall.go )
		if aff != affine.Record {
			err = errutil.Fmt("attempting to call a record %q with affine %q", pat, aff)
		} else {
			ret = g.RecordOf(rec)
		}
	}
	return
}
