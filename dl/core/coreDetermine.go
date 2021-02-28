package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/pattern"

	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// Determine helps run a pattern.
// It implements every core evaluation,
// erroring if the value requested doesnt support the error returned.
type Determine struct {
	Pattern   pattern.PatternName // a text eval here would be like a function pointer maybe...
	Arguments *Arguments          // pattern args kept as a pointer for composer formatting...
}

func (*Determine) Compose() composer.Spec {
	return composer.Spec{
		Spec:  "{pattern%name:pattern_name}{?arguments}",
		Group: "patterns",
		Desc:  "Runs a pattern, and potentially returns a value.",
		Stub:  true,
	}
}

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

func (op *Determine) determine(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	var args []rt.Arg
	if op.Arguments != nil { // FIX!!!!!!!!
		for _, a := range op.Arguments.Args {
			args = append(args, rt.Arg{a.Name, a.From})
		}
	}
	name := op.Pattern.String()
	if v, e := run.Call(name, aff, args); e != nil {
		err = cmdErrorCtx(op, name, e)
	} else {
		ret = v
	}
	return
}
