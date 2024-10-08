package weave

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// domains are processed one at a time,
// starting with those which have zero dependencies,
// ending with those which have the most dependencies.
type Processing struct {
	queue []scheduled
}

type scheduled struct {
	on  weaver.Phase   //
	pos compact.Source // source of the request
	// because the requested phase can be zero
	// ( for "try every phrase" ) the actual phase is passed in the callback
	// fix: what do receivers use that for, and why do they need it?
	req func(weaver.Phase, *mdl.Pen) error
}

// manufacture a particular weaver for a particular phase and source position
type PenCreator interface {
	NewPen(weaver.Phase, compact.Source) *mdl.Pen
}

// queue the scheduled callback for processing.
func (proc *Processing) Schedule(when weaver.Phase, pos compact.Source, req func(weaver.Phase, *mdl.Pen) error) (err error) {
	proc.queue = append(proc.queue, scheduled{when, pos, req})
	return
}

func (proc *Processing) RunPhases(cp PenCreator) (err error) {
	for z := weaver.Phase(1); z < weaver.NumPhases; z++ {
		if e := proc.UpdatePhase(cp, z); e != nil {
			err = e
			break
		}
	}
	return
}

// run the passed phase until all scheduled callbacks have finished.
// error if they didn't finish ( for example, after trying all of them and failing )
func (proc *Processing) UpdatePhase(cp PenCreator, now weaver.Phase) (err error) {
	var exitNextLoop bool
Error:
	for {
		var lastRetry error
		keep := proc.queue[:0] // compact onto the same memory block
		var pending, successes, jessCount int
		// loop through the queue
		// additions might happen as we process elements
		for len(proc.queue) > 0 {
			// pop the front element
			next := proc.queue[0]
			proc.queue = proc.queue[1:]
			if now < next.on {
				keep = append(keep, next)
			} else {
				// run the scheduled request:
				// on success, drop it ( don't place into keep. )
				if e := next.req(now, cp.NewPen(now, next.pos)); e == nil {
					successes++
				} else if !errors.Is(e, mdl.ErrMissing) {
					// stop on critical error
					err = e
					break Error
				} else {
					// if it's the special continuous phase
					// don't consider it pending
					if next.on == 0 {
						jessCount++
					} else {
						lastRetry = e
						pending++
					}
					keep = append(keep, next)
				}
			}
		}
		// after emptying the queue
		// restore it with everything we kept.
		proc.queue = keep
		// if nothing is pending for this phase,
		// make sure we've tried all the jess phrases
		// against whatever just got finished.
		if pending == 0 {
			if jessCount == 0 || exitNextLoop {
				break
			} else {
				exitNextLoop = true
				continue
			}
		} else if successes == 0 {
			// if nothing progressed; error.
			e := fmt.Errorf("couldn't finish %s", now)
			err = errors.Join(lastRetry, e)
			break
		}
		// otherwise, loop.
		exitNextLoop = false
	}
	return
}
