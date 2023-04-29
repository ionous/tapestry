package render

import (
	"bytes"
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RenderPattern) Execute(run rt.Runtime) error {
	_, err := op.RenderEval(run, affine.None)
	return err
}

func (op *RenderPattern) GetBool(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.Bool)
}

func (op *RenderPattern) GetNumber(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.Number)
}

// expressions are text patterns... so for now adapt via text
// ideally could generate the buffer based on the pattern type at assembly type
func (op *RenderPattern) GetText(run rt.Runtime) (ret g.Value, err error) {
	var buf bytes.Buffer
	if v, e := core.WriteSpan(run, &buf, &buf, func() error {
		_, e := op.render(run, affine.None)
		return e
	}); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) GetRecord(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.Record)
}

func (op *RenderPattern) GetNumList(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.NumList)
}

func (op *RenderPattern) GetTextList(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.TextList)
}

func (op *RenderPattern) GetRecordList(run rt.Runtime) (g.Value, error) {
	return op.RenderEval(run, affine.RecordList)
}

// one of the above evals might be called, or this might be called directly from a different pattern
// the hint tells us what return value type is expected.
func (op *RenderPattern) RenderEval(run rt.Runtime, hint affine.Affinity) (ret g.Value, err error) {
	if v, e := op.render(run, hint); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) render(run rt.Runtime, hint affine.Affinity) (ret g.Value, err error) {
	name := op.PatternName
	if rec, e := assign.MakeRecord(run, name); e != nil {
		err = e
	} else {
		k := rec.Kind()
		if have, want := len(op.Render), k.NumField(); have > want {
			err = errutil.Fmt("too many arguments for %s have %d want %d", name, have, want)
		} else {
			for i, el := range op.Render {
				field := k.Field(i) // use render value to ask for the right type of value
				if v, e := el.RenderEval(run, field.Affinity); e != nil {
					err = errutil.New("rendering", name, "arg", i, e)
					break
				} else if val, e := safe.AutoConvert(run, k.Field(i), v); e != nil {
					err = errutil.New("converting", name, "arg", i, e)
					break
				} else if e := rec.SetIndexedField(i, val); e != nil {
					err = errutil.New("setting name", name, "arg", i, e)
					break
				}
			}
			if err == nil {
				if v, e := run.Call(rec, hint); e != nil && !errors.Is(e, rt.NoResult{}) {
					err = errutil.New("calling", name, e)
				} else {
					ret = v
				}
			}
		}
	}
	return
}
