package pattern

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

type stopJump struct {
	stop     bool // true to stop processing all rules after the jump
	jump     rt.Jump
	runCount int
}

// after a rule has matched; combine its desired stop/jump with the current set
// it can only become more strict, not less.
func (n *stopJump) ranRule(stop bool, jump rt.Jump) {
	n.mergeStop(stop)
	n.mergeJump(jump)
	n.runCount++
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

// reads and resets event cancel, event interrupt from the passed event object
func (n *stopJump) mergeEvent(evtObj *g.Record) (err error) {
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
