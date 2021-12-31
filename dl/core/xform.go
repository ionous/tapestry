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
	} else if t := t.String(); len(t) == 0 {
		ret = g.Empty
	} else {
		singular := run.SingularOf(t)
		ret = g.StringOf(singular)
	}
	return
}

func (op *Pluralize) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if t := t.String(); len(t) == 0 {
		ret = g.Empty
	} else {
		plural := run.PluralOf(t)
		ret = g.StringOf(plural)
	}
	return
}

func (op *Capitalize) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if s := t.String(); len(s) == 0 {
		ret = g.Empty
	} else {
		var out strings.Builder
		for _, ch := range s {
			if out.Len() == 0 {
				ch = unicode.ToUpper(ch)
			}
			out.WriteRune(ch)
		}
		ret = g.StringOf(out.String())
	}
	return
}

func (op *MakeLowercase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		lwr := strings.ToLower(t.String())
		ret = g.StringOf(lwr)
	}
	return
}

func (op *MakeUppercase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		upper := strings.ToUpper(t.String())
		ret = g.StringOf(upper)
	}
	return
}

func (op *MakeTitleCase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if t := t.String(); len(t) == 0 {
		ret = g.Empty
	} else {
		title := lang.Titlecase(t)
		ret = g.StringOf(title)
	}
	return
}

func (op *MakeSentenceCase) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else if t := t.String(); len(t) == 0 {
		ret = g.Empty
	} else {
		sentence := lang.SentenceCase(t)
		ret = g.StringOf(sentence)
	}
	return
}

func (op *MakeReversed) GetText(run rt.Runtime) (ret g.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		a := []rune(t.String())
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		ret = g.StringOf(string(a))
	}
	return
}
