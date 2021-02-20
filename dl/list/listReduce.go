package list

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// A normal reduce would return a value, instead we accumulate into a variable
type Reduce struct {
	IntoValue    string
	FromList     rt.Assignment
	UsingPattern pattern.PatternName
}

func (*Reduce) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_reduce",
		Group: "list",
		Desc: `Reduce List: Transform the values from one list by combining them into a single value.
		The named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).`,
	}
}

func (op *Reduce) Execute(run rt.Runtime) (err error) {
	if e := op.reduce(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Reduce) reduce(run rt.Runtime) (err error) {
	if fromList, e := safe.GetAssignedValue(run, op.FromList); e != nil {
		err = e
	} else if outVal, e := safe.CheckVariable(run, op.IntoValue, ""); e != nil {
		err = e
	} else {
		var pat pattern.Pattern
		if e := run.GetEvalByName(op.UsingPattern.String(), &pat); e != nil {
			err = e
		} else {
			for it := g.ListIt(fromList); it.HasNext(); {
				if inVal, e := it.GetNext(); e != nil {
					err = e
					break
				} else {
					args := []*core.Argument{
						{"$1", &fromVal{inVal}},
						{"$2", &fromVal{outVal}},
					}
					if newVal, e := pat.Run(run, args, outVal.Affinity()); e != nil {
						err = e
						break
					} else {
						// send it back in for the next time.
						outVal = newVal
					}
				}
			}
			if err == nil {
				// write back the completed value
				err = run.SetField(object.Variables, op.IntoValue, outVal)
			}
		}
	}
	return
}
