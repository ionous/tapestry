package rel

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *RelativeOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if rel, e := safe.GetText(run, op.RelationName); e != nil {
		err = e
	} else if a, e := safe.ObjectText(run, op.NounName); e != nil {
		err = cmd.Error(op, e)
	} else {
		if noun := a.String(); len(noun) == 0 {
			ret = rt.Nothing // fix: if there's 'a' has a type, we should probably return the reciprocated type
		} else {
			if vs, e := run.RelativesOf(noun, rel.String()); e != nil {
				err = cmd.Error(op, e)
			} else if cnt := vs.Len(); cnt > 1 {
				e := fmt.Errorf("expected at most one relative for %q in %q", noun, rel)
				err = cmd.Error(op, e)
			} else {
				if cnt != 0 { // having no relatives is considered okay
					ret = vs.Index(0)
				} else {
					ret = rt.StringFrom("", vs.Type())
				}
			}
		}
	}
	return
}

func (op *RelativesOf) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if rel, e := safe.GetText(run, op.RelationName); e != nil {
		err = e
	} else if a, e := safe.ObjectText(run, op.NounName); e != nil {
		err = cmd.Error(op, e)
	} else {
		if noun := a.String(); len(noun) == 0 {
			ret = rt.StringsOf(nil)
		} else if vs, e := run.RelativesOf(noun, rel.String()); e != nil {
			err = cmd.Error(op, e)
		} else {
			ret = vs
		}
	}
	return
}

func (op *ReciprocalOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if rel, e := safe.GetText(run, op.RelationName); e != nil {
		err = e
	} else if a, e := safe.ObjectText(run, op.NounName); e != nil {
		err = cmd.Error(op, e)
	} else {
		if noun := a.String(); len(noun) == 0 {
			ret = rt.Nothing // fix: if there's 'a' has a type, we should probably return the reciprocated type
		} else if vs, e := run.ReciprocalsOf(noun, rel.String()); e != nil {
			err = cmd.Error(op, e)
		} else if cnt := vs.Len(); cnt > 1 {
			e := fmt.Errorf("expected at most one reciprocal for %q in %q", noun, rel)
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
	if rel, e := safe.GetText(run, op.RelationName); e != nil {
		err = e
	} else if a, e := safe.ObjectText(run, op.NounName); e != nil {
		err = cmd.Error(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = rt.StringsOf(nil)
	} else if vs, e := run.ReciprocalsOf(a, rel.String()); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = vs
	}
	return
}
