package core

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *ObjectExists) GetBool(run rt.Runtime) (ret g.Value, err error) {
	switch obj, e := safe.ObjectText(run, op.Object); e.(type) {
	case nil:
		if len(obj.String()) == 0 {
			ret = g.False
		} else {
			ret = g.True
		}
	case g.Unknown:
		ret = g.False // fix: is this branch even possible?
	default:
		err = cmdError(op, e)
	}
	return
}

func (op *IdOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (op *NameOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.Empty // fix: or, should it be "nothing"
	} else if v, e := run.GetField(meta.ObjectName, obj); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// ex. repeating across all things
func (op *KindsOf) GetTextList(run rt.Runtime) (g.Value, error) {
	kind := lang.Underscore(op.Kind) // fix: at assembly time.
	return run.GetField(meta.ObjectsOfKind, kind)
}

// returns the kind of this type
func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if k, e := objectKind(run, op.Object); e != nil {
		err = e
	} else if k == nil {
		ret = g.Empty
	} else {
		ret = g.StringOf(k.Name())
	}
	return
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if k, e := objectKind(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ok := k != nil && k.Implements(lang.Underscore(op.Kind))
		ret = g.BoolOf(ok)
	}
	return
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if k, e := objectKind(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ok := k != nil && k.Name() == lang.Underscore(op.Kind)
		ret = g.BoolOf(ok)
	}
	return
}

// get the kind of the passed object reference
// handles null references which in some cases still have a type
// can return nil for a empty reference
// ( an invalid reference returns error )
func objectKind(run rt.Runtime, eval rt.TextEval) (ret *g.Kind, err error) {
	if eval == nil {
		err = safe.MissingEval("object text")
	} else if text, e := eval.GetText(run); e != nil {
		err = e
	} else {
		// if the object name is blank, we might still be able to glean some info on kind...
		// ( useful for things like "nobody" vs "nothing"
		if name := text.String(); len(name) == 0 {
			kind := text.Type()
			if len(kind) != 0 {
				ret, err = run.GetKindByName(kind)
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
