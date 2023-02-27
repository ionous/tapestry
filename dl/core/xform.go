package core

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *Singularize) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		singular := run.SingularOf(str)
		ret = g.StringFrom(singular, t.Type())
	}
	return
}

func (op *Pluralize) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		plural := run.PluralOf(str)
		ret = g.StringFrom(plural, t.Type())
	}
	return
}

func (op *Capitalize) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		var out strings.Builder
		for _, ch := range t.String() {
			if out.Len() == 0 {
				ch = unicode.ToUpper(ch)
			}
			out.WriteRune(ch)
		}
		ret = g.StringFrom(out.String(), t.Type())
	}
	return
}

func (op *MakeLowercase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		lwr := strings.ToLower(t.String())
		ret = g.StringFrom(lwr, t.Type())
	}
	return
}

func (op *MakeUppercase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		upper := strings.ToUpper(t.String())
		ret = g.StringFrom(upper, t.Type())
	}
	return
}

func (op *MakeTitleCase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		title := lang.Titlecase(str)
		ret = g.StringFrom(title, t.Type())
	}
	return
}

func (op *MakeSentenceCase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		sentence := lang.SentenceCase(str)
		ret = g.StringFrom(sentence, t.Type())
	}
	return
}

func (op *MakeReversed) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		a := []rune(t.String())
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		ret = g.StringFrom(string(a), t.Type())
	}
	return
}
