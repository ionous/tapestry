package text

import (
	"bytes"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *IsEmpty) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else {
		b := len(t.String()) == 0
		ret = rt.BoolOf(b)
	}
	return
}

func (op *FindText) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.ErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Subtext); e != nil {
		err = cmd.ErrorCtx(op, "Part", e)
	} else {
		contains := strings.Contains(text.String(), part.String())
		ret = rt.BoolOf(contains)
	}
	return
}

func (op *FindText) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.ErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Subtext); e != nil {
		err = cmd.ErrorCtx(op, "Part", e)
	} else {
		idx := strings.Index(text.String(), part.String())
		ret = rt.IntOf(idx)
	}
	return
}

func (op *TextStartsWith) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.ErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Subtext); e != nil {
		err = cmd.ErrorCtx(op, "Part", e)
	} else {
		contains := strings.HasPrefix(text.String(), part.String())
		ret = rt.BoolOf(contains)
	}
	return
}
func (op *TextEndsWith) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.ErrorCtx(op, "Text", e)
	} else if part, e := safe.GetText(run, op.Subtext); e != nil {
		err = cmd.ErrorCtx(op, "Part", e)
	} else {
		contains := strings.HasSuffix(text.String(), part.String())
		ret = rt.BoolOf(contains)
	}
	return
}

func (op *Join) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if sep, e := safe.GetOptionalText(run, op.Sep, ""); e != nil {
		err = cmd.ErrorCtx(op, "Sep", e)
	} else {
		var buf bytes.Buffer
		sep := sep.String()
		for _, part := range op.Parts {
			if txt, e := safe.GetText(run, part); e != nil {
				err = cmd.ErrorCtx(op, "Part", e)
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
			ret = rt.StringOf(str)
		}
	}
	return
}
