package core

import (
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *CallMake) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	if d, e := op.makeRecord(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.RecordOf(d)
	}
	return
}

func (op *CallMake) makeRecord(run rt.Runtime) (ret *g.Record, err error) {
	if k, e := run.GetKindByName(op.Kind); e != nil {
		err = e
	} else {
		out := k.NewRecord()
		if args := op.Arguments.Args; len(args) == 0 {
			ret = out // return the empty record
		} else {
			for _, arg := range args {
				name := lang.Underscore(arg.Name)
				if fin := k.FieldIndex(name); fin < 0 {
					e := g.UnknownField(op.Kind, arg.Name)
					err = errutil.Append(err, e)
				} else if val, e := safe.GetAssignedValue(run, arg.From); e != nil {
					err = errutil.Append(err, e)
				} else if e := out.SetNamedField(name, val); e != nil {
					err = errutil.Append(err, e)
					// fix? we have to set by name to handle traits
					// would it make more sense to switch out here for that?
					// or possibly handle traits at the indexed level?
				}
			}
			if err == nil {
				ret = out
			}
		}
	}
	return
}
