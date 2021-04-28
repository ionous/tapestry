package core

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// PrintNum writes a number using numerals, eg. "1".
// Numeral:
type PrintNum struct {
	Num rt.NumberEval
}

// PrintNumWord writes a number using english: eg. "one".
// Numeral words:5
type PrintNumWord struct {
	Num rt.NumberEval `if:"pb=words"`
}

// Compose defines a spec for the composer editor.
func (*PrintNum) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "numeral",
		Spec:  "as text {num:number_eval}",
		Desc:  "A number as text: Writes a number using numerals, eg. '1'.",
		Group: "printing",
	}
}

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

// Compose defines a spec for the composer editor.
func (*PrintNumWord) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "numeral",
		Desc:  "A number in words: Writes a number in plain english: eg. 'one'",
		Group: "printing",
	}
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
