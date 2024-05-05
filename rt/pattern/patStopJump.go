package pattern

import (
	"fmt"

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
func (n *stopJump) mergeEvent(evt *rt.Record) (err error) {
	if status, e := evt.GetIndexedField(event.Status.Index()); e != nil {
		err = e
	} else {
		// turn the aspect into a trait index
		curr, state := status.String(), event.CancellationStatus(-1)
		for i := 0; i < event.NumStatus; i++ {
			check := event.CancellationStatus(i)
			if curr == check.String() {
				state = check
				break
			}
		}
		//
		switch state {
		case event.ContinueNormally:
			// do nothing
		case event.InterruptLater:
			n.interrupt(false)
		case event.InterruptNow:
			n.interrupt(true)
		case event.CancelLater:
			n.cancel(false)
		case event.CancelNow:
			n.cancel(true)
		default:
			err = fmt.Errorf("unknown event status %s", curr)
		}
	}
	return
}
