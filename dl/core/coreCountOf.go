package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type CountOf struct {
	Pos reader.Position `if:"internal"` // generated at import time to provide a unique counter for each sequence
	Num rt.NumberEval   `if:"selector"`
	Trigger
}

type TriggerOnce struct{}
type TriggerCycle struct{}
type TriggerSwitch struct{}

type Trigger interface{ Trigger() Trigger }

func (op *TriggerOnce) Trigger() Trigger   { return op }
func (op *TriggerCycle) Trigger() Trigger  { return op }
func (op *TriggerSwitch) Trigger() Trigger { return op }

func (*TriggerOnce) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "once", Role: composer.Selector},
		Group:  "comparison",
		Stub:   true,
	}
}

func (*TriggerCycle) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "cycle", Role: composer.Selector},
		Group:  "comparison",
	}
}

func (*TriggerSwitch) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "switch", Role: composer.Selector},
		Group:  "comparison",
	}
}

func (*CountOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "logic",
		Fluent: &composer.Fluid{Name: "countOf", Role: composer.Function},
		Desc: `CountOf: A guard which returns true based on a counter. 
Counters start at zero and are incremented every time the guard gets checked.`,
		Stub: true,
	}
}

func (op *CountOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if ok, e := op.update(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(ok)
	}
	return
}

func (op *CountOf) update(run rt.Runtime) (okay bool, err error) {
	name := op.Pos.String()
	if p, e := run.GetField(object.Counter, name); e != nil {
		err = e
	} else if count := p.Int(); count >= 0 {
		//
		if target, e := safe.GetNumber(run, op.Num); e != nil {
			err = e
		} else {
			// determine the next value
			next := count + 1
			// are we at or above the target?
			if okay = next >= target.Int(); okay {
				switch op.Trigger.(type) {
				case *TriggerOnce:
					next = -1
				case *TriggerCycle:
					next = 0
				case *TriggerSwitch:
					next = count
				}
			}
			// set back the counter
			if e := run.SetField(object.Counter, name, g.IntOf(next)); e != nil {
				err = e
			}
		}
	}
	return
}
