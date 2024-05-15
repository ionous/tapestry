package render

import (
	"bytes"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RenderPattern) Execute(run rt.Runtime) error {
	_, err := op.RenderEval(run, affine.None)
	return err
}

func (op *RenderPattern) GetBool(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.Bool)
}

func (op *RenderPattern) GetNumber(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.Number)
}

// expressions are text patterns... so for now adapt via text
// ideally could generate the buffer based on the pattern type at assembly type
func (op *RenderPattern) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.getText(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) getText(run rt.Runtime) (ret rt.Value, err error) {
	var buf bytes.Buffer
	prev := run.SetWriter(&buf)
	if _, e := op.render(run, affine.None); e != nil {
		err = e
	} else if str := buf.String(); len(str) > 0 {
		ret = rt.StringOf(str)
	} else {
		ret = safe.GetTemplateText()
	}
	run.SetWriter(prev)
	return
}

func (op *RenderPattern) GetRecord(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.Record)
}

func (op *RenderPattern) GetNumList(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.NumList)
}

func (op *RenderPattern) GetTextList(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.TextList)
}

func (op *RenderPattern) GetRecordList(run rt.Runtime) (rt.Value, error) {
	return op.RenderEval(run, affine.RecordList)
}

// one of the above evals might be called, or this might be called directly from a different pattern
// the hint tells us what return value type is expected.
func (op *RenderPattern) RenderEval(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if v, e := op.render(run, hint); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) render(run rt.Runtime, hint affine.Affinity) (ret rt.Value, err error) {
	if k, e := run.GetKindByName(op.PatternName); e != nil {
		err = e
	} else {
		name := k.Name()
		vals := make([]rt.Value, len(op.Render))
		for i, el := range op.Render { // use the targeted field to know how to read the value
			if v, e := el.RenderEval(run, k.Field(i).Affinity); e != nil {
				err = fmt.Errorf("%w rendering %s arg %d", e, name, i)
				break
			} else {
				vals[i] = v
			}
		}
		if err == nil {
			if v, e := run.Call(name, hint, nil, vals); e != nil {
				err = fmt.Errorf("%w calling %s", e, name)
			} else {
				ret = v
			}
		}
	}
	return
}
