package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// returns the author specified name for the indicated object.
// returns an error if there is no such object;
// returns the empty string for an empty request.
func (op *NameOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if id, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := id.String(); len(obj) == 0 {
		ret = id
	} else if v, e := run.GetField(meta.ObjectName, obj); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *ObjectTraits) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if ts, e := op.getTraits(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = rt.StringsOf(ts)
	}
	return
}

func (op *ObjectTraits) getTraits(run rt.Runtime) (ret []string, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = e
	} else if kind, e := run.GetField(meta.ObjectKind, obj.String()); e != nil {
		err = e
	} else if k, e := run.GetKindByName(kind.String()); e != nil {
		err = e
	} else {
		for i, cnt := 0, k.FieldCount(); i < cnt; i++ {
			if f := k.Field(i); f.Name == f.Type { // aspect like
				if a, e := run.GetKindByName(f.Type); e == nil && a.Implements(kindsOf.Aspect.String()) {
					if str, e := run.GetField(obj.String(), f.Name); e != nil {
						err = e
						break
					} else {
						ret = append(ret, str.String())
					}
				}
			}
		}
	}
	return
}

// returns a list of all objects of the specified kind.
func (op *KindsOf) GetTextList(run rt.Runtime) (rt.Value, error) {
	kind := inflect.Normalize(op.Kind) // fix: at assembly time.
	return run.GetField(meta.ObjectsOfKind, kind)
}

// returns the kind of the indicated object.
// returns an error if there is no such object; returns the empty string for an empty request.
func (op *KindOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Object, op.Nothing); e != nil {
		err = e
	} else if k == nil {
		ret = rt.Empty
	} else {
		ret = rt.StringOf(k.Name()) // tbd: should kind string have a type of meta.ObjectKind?
	}
	return
}

// returns true if the indicated object is of the specified kind.
// returns an error if there is no such object;
// returns the false for an empty request UNLESS nothing objects were specified as being allowed to match.
func (op *IsKindOf) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Object, op.Nothing); e != nil {
		err = cmdError(op, e)
	} else {
		ok := k != nil && k.Implements(inflect.Normalize(op.Kind))
		ret = rt.BoolOf(ok)
	}
	return
}

// returns true if the indicated object is of the specified kind
// but not a kind that derives from the specified kind.
// returns an error if there is no such object; returns the false for an empty request.
func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Object, false); e != nil {
		err = cmdError(op, e)
	} else {
		ok := k != nil && k.Name() == inflect.Normalize(op.Kind)
		ret = rt.BoolOf(ok)
	}
	return
}

// get the kind of the passed object reference
// handles null references which in some cases still have a type
// can return nil for a empty reference
// ( an invalid reference returns error )
func objectKind(run rt.Runtime, eval rt.TextEval, allowNothing bool) (ret *rt.Kind, err error) {
	if eval == nil {
		err = safe.MissingEval("object text")
	} else if text, e := eval.GetText(run); e != nil {
		err = e
	} else {
		// if the object name is blank, we might still be able to glean some info on kind...
		// ( useful for things like "nobody" vs "nothing"
		if name := text.String(); len(name) == 0 {
			if allowNothing {
				kind := text.Type()
				if len(kind) != 0 {
					ret, err = run.GetKindByName(kind)
				}
			}
		} else {
			// fix? we cant fully rely on the type of the text because of lists
			// lists are considered all of one type right now: they are not bags.
			if obj, e := run.GetField(meta.ObjectId, name); e != nil {
				err = e
			} else if kind, e := run.GetField(meta.ObjectKind, obj.String()); e != nil {
				err = e
			} else {
				ret, err = run.GetKindByName(kind.String())
			}
		}
	}
	return
}
