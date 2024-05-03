package pattern

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
)

type stopJump struct {
	stop     bool // true to stop processing all rules after the jump
	jump     rt.Jump
	runCount int
}

// after a rule has matched; combine its desired stop/jump with the current set
// it can only become more strict, not less.
// if there's a return value, it must be set for the pattern to be considered done.
func (n *stopJump) update(status stopJump, evtObj *rt.Record, result bool) (done bool, err error) {
	if result && status.runCount > 0 {
		n.mergeStop(status.stop)
		n.mergeJump(status.jump)
		n.runCount++
	}
	if evtObj != nil {
		err = n.mergeEvent(evtObj)
	}
	if err == nil {
		done = n.jump == rt.JumpNow
	}
	return
}

// when now is false, delegates the jump behavior to interrupt
func (n *stopJump) cancel(now bool) {
	if now {
		n.mergeJump(rt.JumpNow)
	}
	n.mergeStop(true)
}

func (n *stopJump) interrupt(now bool) {
	if now {
		n.mergeJump(rt.JumpNow)
	} else {
		n.mergeJump(rt.JumpNext)
	}
}

func (n *stopJump) mergeStop(stop bool) {
	if stop && !n.stop {
		n.stop = true
	}
}

func (n *stopJump) mergeJump(jump rt.Jump) {
	if jump < n.jump {
		n.jump = jump
	}
}

// reads and resets event cancel, event interrupt from the passed event object
func (n *stopJump) mergeEvent(evtObj *rt.Record) (err error) {
	if i := event.Cancel.Index(); evtObj != nil && evtObj.HasValue(i) {
		if cancel, e := evtObj.GetIndexedField(i); e != nil {
			err = e
		} else {
			n.cancel(cancel.Bool())
			_ = evtObj.SetIndexedField(i, nil)
		}
	}

	if i := event.Interupt.Index(); evtObj != nil && evtObj.HasValue(i) {
		if interrupt, e := evtObj.GetIndexedField(i); e != nil {
			err = e
		} else {
			n.interrupt(interrupt.Bool())
			_ = evtObj.SetIndexedField(i, nil)
		}
	}
	return
}
