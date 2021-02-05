package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type Map struct {
	ToList       string
	FromList     core.Assignment
	UsingPattern pattern.PatternName
}

func (*Map) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_map",
		Group: "list",
		Desc: `Map List: Transform the values from one list and place the results in another list.
		The named pattern is called with two records 'in' and 'out' from the source and output lists respectively.`,
	}
}

func (op *Map) Execute(run rt.Runtime) (err error) {
	if e := op.remap(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Map) remap(run rt.Runtime) (err error) {
	if fromList, e := core.GetAssignedValue(run, op.FromList); e != nil {
		err = errutil.New("from_list:", op.FromList, e)
	} else if toList, e := safe.List(run, op.ToList); e != nil {
		err = errutil.New("to_list:", op.ToList, e)
	} else {
		var pat pattern.Pattern
		if e := run.GetEvalByName(op.UsingPattern.String(), &pat); e != nil {
			err = e
		} else if ps, e := pat.NewRecord(run); e != nil {
			err = e
		} else {
			var pk *g.Kind
			for it := g.ListIt(fromList); it.HasNext(); {
				// create a new set of parameters each loop
				if pk == nil {
					pk = ps.Kind()
				} else {
					ps = pk.NewRecord()
				}
				in, out := 0, 1
				if inVal, e := it.GetNext(); e != nil {
					err = e
					break
				} else if e := ps.SetIndexedField(in, inVal); e != nil {
					err = e
					break
				} else if _, e := pat.Run(run, ps, ""); e != nil {
					err = e
					break
				} else if newVal, e := ps.GetIndexedField(out); e != nil {
					err = e
					break
				} else if src, dst := newVal.Affinity(), toList.Affinity(); src != affine.Element(dst) ||
					((src == affine.Record) && newVal.Type() != toList.Type()) {
					err = errutil.New("elements don't match")
					break
				} else {
					toList.Append(newVal)
				}
			}
			if err == nil {
				err = run.SetField(object.Variables, op.ToList, toList)
			}
		}
	}
	return
}
