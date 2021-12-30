package core

import (
	"bytes"
	"strings"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *IsEmpty) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		b := len(t.String()) == 0
		ret = g.BoolOf(b)
	}
	return
}

func (op *Includes) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmdErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Part); e != nil {
		err = cmdErrorCtx(op, "Part", e)
	} else {
		contains := strings.Contains(text.String(), part.String())
		ret = g.BoolOf(contains)
	}
	return
}

func (op *Join) GetText(run rt.Runtime) (ret g.Value, err error) {
	if sep, e := safe.GetOptionalText(run, op.Sep, ""); e != nil {
		err = cmdErrorCtx(op, "Sep", e)
	} else {
		var buf bytes.Buffer
		sep := sep.String()
		for _, part := range op.Parts {
			if txt, e := safe.GetText(run, part); e != nil {
				err = cmdErrorCtx(op, "Part", e)
				break
			} else {
				if buf.Len() > 0 {
					buf.WriteString(sep)
				}
				str := txt.String()
				buf.WriteString(str)
			}
		}
		if err == nil {
			str := buf.String()
			ret = g.StringOf(str)
		}
	}
	return
}
