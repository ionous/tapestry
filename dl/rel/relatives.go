package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *RelativeOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmd.Error(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = rt.Empty // fix: if there's 'a' has a type, we should probably return the reciprocated type
	} else {
		noun, rel := a, op.Via
		if vs, e := run.RelativesOf(noun, rel); e != nil {
			err = cmd.Error(op, e)
		} else if cnt := vs.Len(); cnt > 1 {
			e := errutil.New("expected at most one relative for", noun, "in", rel)
			err = cmd.Error(op, e)
		} else {
			if cnt != 0 { // having no relatives is considered okay
				ret = vs.Index(0)
			} else {
				ret = rt.StringFrom("", vs.Type())
			}
		}
	}
	return
}

func (op *RelativesOf) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmd.Error(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = rt.StringsOf(nil)
	} else if vs, e := run.RelativesOf(a, op.Via); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *ReciprocalOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmd.Error(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = rt.Empty // fix: if there's 'a' has a type, we should probably return the reciprocated type
	} else {
		noun, rel := a, op.Via
		if vs, e := run.ReciprocalsOf(noun, rel); e != nil {
			err = cmd.Error(op, e)
		} else if cnt := vs.Len(); cnt > 1 {
			e := errutil.New("expected at most one reciprocal for", noun, "in", rel)
			err = cmd.Error(op, e)
		} else {
			if cnt != 0 { // having no relatives is considered okay
				ret = vs.Index(0)
			} else {
				ret = rt.StringFrom("", vs.Type())
			}
		}
	}
	return
}

func (op *ReciprocalsOf) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmd.Error(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = rt.StringsOf(nil)
	} else if vs, e := run.ReciprocalsOf(a, op.Via); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = vs
	}
	return
}
