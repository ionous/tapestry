package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *Map) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) remap(run rt.Runtime) (err error) {
	if fromList, e := safe.GetAssignedValue(run, op.FromList); e != nil {
		err = errutil.New("from_list:", op.FromList, e)
	} else if toList, e := safe.List(run, op.ToList.Value()); e != nil {
		err = errutil.New("to_list:", op.ToList.Value(), e)
	} else {
		pat := op.UsingPattern
		aff := affine.Element(toList.Affinity())
		//
		for it := g.ListIt(fromList); it.HasNext(); {
			if inVal, e := it.GetNext(); e != nil {
				err = e
				break
			} else {
				if newVal, e := run.Call(pat, aff, []rt.Arg{
					{"$1", &fromVal{inVal}},
				}); e != nil {
					// note: we treat no result as an error because
					// we are trying to map *all* of the elements from one list into another
					err = e
					break
				} else {
					toList.Append(newVal)
				}
			}
		}
	}
	return
}

// allows values to be passed as arguments ( arguments are usually evals )
type fromVal struct{ val g.Value }

func (op *fromVal) GetAssignedValue(rt.Runtime) (g.Value, error) {
	return op.val, nil
}
func (op *fromVal) Affinity() affine.Affinity {
	return op.val.Affinity()
}
