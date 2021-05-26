package core

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func (op *PrintNum) GetText(run rt.Runtime) (ret g.Value, err error) {
	if n, e := safe.GetNumber(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s := strconv.FormatFloat(n.Float(), 'g', -1, 64); len(s) > 0 {
		ret = g.StringOf(s)
	} else {
		ret = g.StringOf("<num>")
	}
	return
}

func (op *PrintNumWord) GetText(run rt.Runtime) (ret g.Value, err error) {
	if n, e := safe.GetNumber(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s, ok := lang.NumToWords(n.Int()); ok {
		ret = g.StringOf(s)
	} else {
		ret = g.StringOf("<num>")
	}
	return
}
