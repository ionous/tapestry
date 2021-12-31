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

func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.Empty
	} else if v, e := run.GetField(meta.ObjectKind, obj); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else if ok, e := safe.IsKindOf(run, obj, lang.Underscore(op.Kind)); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(ok)
	}
	return
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else {
		kind := lang.Underscore(op.Kind)
		if k, e := run.GetField(meta.ObjectKind, obj); e != nil {
			err = cmdError(op, e)
		} else {
			ok := kind == k.String()
			ret = g.BoolOf(ok)
		}
	}
	return
}

// ex. repeating across all things
func (op *KindsOf) GetTextList(run rt.Runtime) (g.Value, error) {
	kind := lang.Underscore(op.Kind) // fix: at assembly time.
	return run.GetField(meta.ObjectsOfKind, kind)
}
