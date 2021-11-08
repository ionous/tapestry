package core

import (
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Trigger interface{ Trigger() Trigger }

func (op *TriggerOnce) Trigger() Trigger   { return op }
func (op *TriggerCycle) Trigger() Trigger  { return op }
func (op *TriggerSwitch) Trigger() Trigger { return op }

func (op *CallTrigger) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if ok, e := op.update(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(ok)
	}
	return
}

func (op *CallTrigger) update(run rt.Runtime) (okay bool, err error) {
	name := op.Name
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
