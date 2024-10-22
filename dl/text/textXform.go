package text

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *Singularize) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		singular := run.SingularOf(str)
		ret = rt.StringFrom(singular, t.Type())
	}
	return
}

func (op *Pluralize) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		plural := run.PluralOf(str)
		ret = rt.StringFrom(plural, t.Type())
	}
	return
}

func (op *Capitalize) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
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
		ret = rt.StringFrom(out.String(), t.Type())
	}
	return
}

func (op *MakeLowercase) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		lwr := strings.ToLower(t.String())
		ret = rt.StringFrom(lwr, t.Type())
	}
	return
}

func (op *MakeUppercase) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		upper := strings.ToUpper(t.String())
		ret = rt.StringFrom(upper, t.Type())
	}
	return
}

func (op *MakeTitleCase) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		title := inflect.Titlecase(str)
		ret = rt.StringFrom(title, t.Type())
	}
	return
}

func (op *MakeSentenceCase) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		sentence := inflect.SentenceCase(str)
		ret = rt.StringFrom(sentence, t.Type())
	}
	return
}

func (op *MakeReversed) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if t, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else if str := t.String(); len(str) == 0 {
		ret = t
	} else {
		a := []rune(t.String())
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		ret = rt.StringFrom(string(a), t.Type())
	}
	return
}
