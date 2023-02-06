package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *GetFromName) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromName(run, op.Name, affine.Bool, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromName) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromName(run, op.Name, affine.Number, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromName) GetText(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromName(run, op.Name, affine.Text, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromName) GetList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromName(run, op.Name, affine.List, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *GetFromName) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if v, e := getFromName(run, op.Name, affine.Record, op.Dot); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

// FIX: convert and warn instead of error on field affinity checks
func getFromName(run rt.Runtime, name rt.TextEval, aff affine.Affinity, path []Dot) (ret g.Value, err error) {
	if name, e := safe.GetText(run, name); e != nil {
		err = e
	} else if root, e := getValueByName(run, name.String(), len(path) > 0); e != nil {
		err = e
	} else if val, e := Peek(run, root, path); e != nil {
		err = e
	} else if e := safe.Check(val, aff); e != nil {
		err = e
	} else {
		ret = val
	}
	return
}

func getValueByName(run rt.Runtime, name string, subFields bool) (ret g.Value, err error) {
	// uppercase names are assumed to be requests for object names.
	if lang.IsCapitalized(name) {
		ret, err = run.GetField(meta.ObjectValue, name)
	} else {
		// otherwise, try as a variable first:
		switch v, e := run.GetField(meta.Variables, name); e.(type) {
		case nil:
			// did the variable contain an object name? then unpack into an object.
			// fix: yikes. can this be removed?
			if aff := v.Affinity(); aff == affine.Text && subFields {
				ret, err = safe.ObjectFromString(run, v.String())
			}
		case g.Unknown:
			// no such variable? try as an object:
			ret, err = run.GetField(meta.ObjectValue, name)
		default:
			err = e
		}
	}
	return
}
