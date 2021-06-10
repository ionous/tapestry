package core

import (
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
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
	} else if v, e := run.GetField(object.Name, obj); e != nil {
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
	} else if v, e := run.GetField(object.Kind, obj); e != nil {
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
	} else {
		kind := lang.Breakcase(op.Kind)
		if objectPath, e := run.GetField(object.Kinds, obj); e != nil {
			err = cmdError(op, e)
		} else {
			// Contains reports whether second is within first.
			cp, ck := objectPath.String()+",", kind+","
			ok := strings.Contains(cp, ck)
			ret = g.BoolOf(ok)
		}
	}
	return
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else {
		kind := lang.Breakcase(op.Kind)
		if objectPath, e := run.GetField(object.Kinds, obj); e != nil {
			err = cmdError(op, e)
		} else {
			// Contains reports whether second is within first.
			cp, ck := objectPath.String()+",", kind+","
			ok := strings.HasPrefix(cp, ck)
			ret = g.BoolOf(ok)
		}
	}
	return
}

func (op *KindsOf) GetTextList(run rt.Runtime) (g.Value, error) {
	kind := lang.Breakcase(op.Kind) // fix: break case at assembly time.
	return run.GetField(object.Nouns, kind)
}
