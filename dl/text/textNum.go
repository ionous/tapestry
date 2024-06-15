package text

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *PrintNumDigits) Execute(run rt.Runtime) (err error) {
	return safe.WriteText(run, op)
}
func (op *PrintNumDigits) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if n, e := safe.GetNum(run, op.Num); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.StringOf(formatNumber(n))
	}
	return
}

func (op *PrintNumWords) Execute(run rt.Runtime) (err error) {
	return safe.WriteText(run, op)
}
func (op *PrintNumWords) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if n, e := safe.GetNum(run, op.Num); e != nil {
		err = cmd.Error(op, e)
	} else {
		s, ok := inflect.NumToWords(n.Int())
		if !ok {
			s = formatNumber(n)
		}
		ret = rt.StringOf(s)
	}
	return
}

func formatNumber(n rt.Value) string {
	return strconv.FormatFloat(n.Float(), 'g', -1, 64)
}
