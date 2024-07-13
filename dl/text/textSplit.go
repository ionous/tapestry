package text

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SplitWords) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else {
		vs := strings.Split(t.String(), "\n")
		ret = rt.StringsOf(vs)
	}
	return
}

func (op *SplitLines) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else {
		vs := strings.Fields(t.String())
		ret = rt.StringsOf(vs)
	}
	return
}

func (op *SplitText) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if sep, e := safe.GetText(run, op.Separator); e != nil {
		err = cmd.Error(op, e)
	} else {
		vs := strings.Split(t.String(), sep.String())
		ret = rt.StringsOf(vs)
	}
	return
}
