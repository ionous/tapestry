package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type Map struct {
	ToList       string
	FromList     rt.Assignment
	UsingPattern pattern.PatternName
}

func (*Map) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_map",
		Group: "list",
		Desc: `Map List: Transform the values from one list and place the results in another list.
		The designated pattern is called with each value from the "from list", one value at a time.`,
	}
}

func (op *Map) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) remap(run rt.Runtime) (err error) {
	if fromList, e := safe.GetAssignedValue(run, op.FromList); e != nil {
		err = errutil.New("from_list:", op.FromList, e)
	} else if toList, e := safe.List(run, op.ToList); e != nil {
		err = errutil.New("to_list:", op.ToList, e)
	} else {
		pat := op.UsingPattern.String()
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
