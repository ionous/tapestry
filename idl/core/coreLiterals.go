package core

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// GetBool implements BoolEval; providing the dl with a boolean literal.
func (op *BoolValue) GetBool(rt.Runtime) (ret g.Value, _ error) {
	ret = g.BoolOf(op.Bool())
	return
}

// GetNumber implements NumberEval providing the dl with a number literal.
func (op *NumValue) GetNumber(rt.Runtime) (ret g.Value, _ error) {
	ret = g.FloatOf(op.Num())
	return
}

// // Int converts to native int.
func (op *NumValue) Int() int {
	return int(op.Num())
}

// Float converts to native float.
func (op *NumValue) Float() float64 {
	return op.Num()
}

// GetText implements interface TextEval providing the dl with a text literal.
func (op *TextValue) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.Text(); e != nil {
		err = cmdError(op, op.Struct, e)
	} else {
		ret = g.StringOf(v)
	}
	return
}

func (op *NumList) GetNumList(rt.Runtime) (ret g.Value, err error) {
	if vs, e := op.Values(); e != nil {
		err = cmdError(op, op.Struct, e)
	} else {
		// fix: add a generic r/o wrapper for capn float list?
		cnt := vs.Len()
		out := make([]float64, cnt)
		for i := 0; i < cnt; i++ {
			out[i] = vs.At(i)
		}
		ret = g.FloatsOf(out)
	}
	return
}

func (op *TextList) GetTextList(rt.Runtime) (ret g.Value, err error) {
	if vs, e := op.Values(); e != nil {
		err = cmdError(op, op.Struct, e)
	} else {
		// fix: add a generic r/o wrapper for capn text list?
		cnt := vs.Len()
		out := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			if v, e := vs.At(i); e != nil {
				err = cmdError(op, op.Struct, e)
				break
			} else {
				out[i] = v
			}
		}
		ret = g.StringsOf(out)
	}
	return
}
