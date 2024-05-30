package core

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *PrintNum) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if n, e := safe.GetNum(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s := strconv.FormatFloat(n.Float(), 'g', -1, 64); len(s) > 0 {
		ret = rt.StringOf(s)
	} else {
		ret = rt.StringOf("<num>")
	}
	return
}

func (op *PrintNumWord) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if n, e := safe.GetNum(run, op.Num); e != nil {
		err = cmdError(op, e)
	} else if s, ok := inflect.NumToWords(n.Int()); ok {
		ret = rt.StringOf(s)
	} else {
		ret = rt.StringOf("<num>")
	}
	return
}
