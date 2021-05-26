package core

import (
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *Response) GetText(run rt.Runtime) (ret g.Value, err error) {
	if safe.GetFlag(run, object.PrintResponseNames) {
		ret = g.StringOf(op.Name.Value())
	} else if op.Text == nil {
		err = cmdError(op, errutil.New("response doesnt have external lookup yet"))
	} else {
		if v, e := safe.GetText(run, op.Text); e != nil {
			err = cmdError(op, e)
		} else {
			ret = v
		}
	}
	return
}
