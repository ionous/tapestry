package object

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// returns the author specified name for the indicated object.
// returns an error if there is no such object;
// returns the empty string for an empty request.
func (op *ObjectName) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if obj, e := getObjectTarget(run, op.Target); e != nil {
		err = cmd.Error(op, e)
	} else {
		if name := obj.String(); len(name) == 0 {
			ret = rt.Nothing
		} else if v, e := run.GetField(meta.ObjectName, name); e != nil {
			err = cmd.Error(op, e)
		} else {
			ret = v
		}
	}
	return
}

func (op *ObjectStates) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if ts, e := op.getTraits(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.StringsOf(ts)
	}
	return
}

func (op *ObjectStates) getTraits(run rt.Runtime) (ret []string, err error) {
	if obj, e := getObjectTarget(run, op.Target); e != nil {
		err = cmd.Error(op, e)
	} else if name := obj.String(); len(name) > 0 {
		if kind, e := run.GetField(meta.ObjectKind, name); e != nil {
			err = e
		} else if k, e := run.GetKindByName(kind.String()); e != nil {
			err = e
		} else {
			for i, cnt := 0, k.FieldCount(); i < cnt; i++ {
				if f := k.Field(i); f.Name == f.Type { // aspect like
					if a, e := run.GetKindByName(f.Type); e == nil && a.Implements(kindsOf.Aspect.String()) {
						if str, e := run.GetField(name, f.Name); e != nil {
							err = e
							break
						} else {
							ret = append(ret, str.String())
						}
					}
				}
			}
		}
	}
	return
}

// returns a list of all objects of the specified kind.
func (op *KindsOf) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := safe.GetText(run, op.KindName); e != nil {
		err = cmd.Error(op, e)
	} else {
		kind := inflect.Normalize(k.String())
		ret, err = run.GetField(meta.ObjectsOfKind, kind)
	}
	return
}

// returns the kind of the indicated object.
// returns an error if there is no such object; returns the empty string for an empty request.
func (op *KindOf) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Target, op.Nothing); e != nil {
		err = e
	} else if k == nil {
		ret = rt.Nothing
	} else {
		ret = rt.StringOf(k.Name()) // tbd: should kind string have a type of meta.ObjectKind?
	}
	return
}

// returns true if the indicated object is of the specified kind.
// returns an error if there is no such object;
// returns the false for an empty request UNLESS nothing objects were specified as being allowed to match.
func (op *IsKindOf) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Target, op.Nothing); e != nil {
		err = cmd.Error(op, e)
	} else if checkKind, e := safe.GetText(run, op.KindName); e != nil {
		err = cmd.Error(op, e)
	} else {
		ok := k != nil && k.Implements(inflect.Normalize(checkKind.String()))
		ret = rt.BoolOf(ok)
	}
	return
}

// returns true if the indicated object is of the specified kind
// but not a kind that derives from the specified kind.
// returns an error if there is no such object; returns the false for an empty request.
func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if k, e := objectKind(run, op.Target, false); e != nil {
		err = cmd.Error(op, e)
	} else if checkKind, e := safe.GetText(run, op.KindName); e != nil {
		err = cmd.Error(op, e)
	} else {
		ok := k != nil && k.Name() == inflect.Normalize(checkKind.String())
		ret = rt.BoolOf(ok)
	}
	return
}

func getObjectTarget(run rt.Runtime, tgt rt.Address) (ret rt.Value, err error) {
	if ref, e := safe.GetReference(run, tgt); e != nil {
		err = e
	} else if val, e := ref.GetValue(); e != nil {
		err = e
	} else if val.Affinity() != affine.Text {
		err = fmt.Errorf("expected a string, have %s", val.Affinity())
	} else {
		ret = val
	}
	return
}

// get the kind of the passed object reference
// handles null references which in some cases still have a type
// can return nil for a empty reference
// ( an invalid reference returns error )
func objectKind(run rt.Runtime, tgt rt.Address, allowNothing bool) (ret *rt.Kind, err error) {
	if text, e := getObjectTarget(run, tgt); e != nil {
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
