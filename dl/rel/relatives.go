package rel

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RelativeOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.Empty
	} else {
		noun, rel := a, op.Via
		if vs, e := run.RelativesOf(noun, rel.String()); e != nil {
			err = cmdError(op, e)
		} else if cnt := vs.Len(); cnt > 1 {
			e := errutil.New("expected at most one relative for", noun, "in", rel)
			err = cmdError(op, e)
		} else {
			if cnt != 0 { // having no relatives is considered okay
				ret = vs.Index(0)
			} else {
				ret = g.StringFrom("", vs.Type())
			}
		}
	}
	return
}

func (op *RelativesOf) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.StringsOf(nil)
	} else if vs, e := run.RelativesOf(a, op.Via.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *ReciprocalOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.Empty
	} else {
		noun, rel := a, op.Via
		if vs, e := run.ReciprocalsOf(noun, rel.String()); e != nil {
			err = cmdError(op, e)
		} else if cnt := vs.Len(); cnt > 1 {
			e := errutil.New("expected at most one reciprocal for", noun, "in", rel.Str)
			err = cmdError(op, e)
		} else {
			if cnt != 0 { // having no relatives is considered okay
				ret = vs.Index(0)
			} else {
				ret = g.StringFrom("", vs.Type())
			}
		}
	}
	return
}

func (op *ReciprocalsOf) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.StringsOf(nil)
	} else if vs, e := run.ReciprocalsOf(a, op.Via.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vs
	}
	return
}
