package weave

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// domains are processed one at a time,
// starting with those which have zero dependencies,
// ending with those which have the most dependencies.
type Processing struct {
	phase weaver.Phase
	queue []scheduled
}

type scheduled struct {
	on weaver.Phase
	// because the requested phase can be zero
	// ( for "try every phrase" ) the actual phase is passed in the callback
	req func(weaver.Phase) error
}

// queue the scheduled callback for processing.
// or, if the desired phase is the current phase, try to run immediately.
// ( trying immediately allows jess sentence matching to rely on prior sentences. )
func (proc *Processing) schedule(when weaver.Phase, req func(weaver.Phase) error) (err error) {
	if now := proc.phase; now > when && when != 0 {
		err = fmt.Errorf("processing %s phase %s already passed", now, when)
	} else {
		if when == now && now != 0 {
			// if its the same phase, try to run immediately:
			if e := req(now); e != nil && !errors.Is(e, mdl.ErrMissing) {
				// queue on missing, otherwise error
				err = e
			}
		}
		if err == nil {
			proc.queue = append(proc.queue, scheduled{when, req})
		}
	}
	return
}

func (proc *Processing) runAll() (err error) {
	for z := weaver.Phase(1); z < weaver.NumPhases; z++ {
		// fix? necessary for re entrant scheduling
		// if the weave callback contained a scheduler
		// this ( and the processing object itself ) could probably go on the stack
		proc.phase = z
		if e := proc.runPhase(z); e != nil {
			err = e
			break
		}
	}
	return
}

// run the passed phase until all scheduled callbacks have finished.
// error if they didn't finish ( for example, after trying all of them and failing )
func (proc *Processing) runPhase(now weaver.Phase) (err error) {
	var afterFirst bool
Error:
	for {
		var lastRetry error
		var progress []scheduled
		// empty out the queue so we can see newly scheduled requests
		progress, proc.queue = proc.queue, nil
		prevCount := len(progress)
		// loop through the queue, compacting into "keep"
		keep := progress[:0]
		for _, next := range progress {
			// phase zero is special ( for jess )
			// run once per phase, and keep in the original order.
			if ((next.on != 0) && (next.on != now)) ||
				(next.on == 0 && afterFirst) {
				keep = append(keep, next)
			} else {
				// run a scheduled request for this phase
				// on success, drop it ( don't place into keep. )
				if e := next.req(now); e != nil {
					// exit on error:
					if !errors.Is(e, mdl.ErrMissing) {
						err = e
						break Error
					} else {
						// otherwise, keep and retry.
						keep = append(keep, next)
						// we dont want "try every frame" to mean error
						if next.on != 0 {
							lastRetry = e
						}
					}
				}
			}
		}
		//
		newlyScheduled := len(proc.queue)
		succeded := len(keep) < prevCount || newlyScheduled > 0
		// if nothing was removed, and nothing added
		// then any error in trying to remove is critical.
		if !succeded && lastRetry != nil {
			e := fmt.Errorf("couldn't finish phase %s", now)
			err = errors.Join(lastRetry, e)
			break
		} else {
			// append any newly scheduled elements
			if newlyScheduled > 0 {
				keep = append(keep, proc.queue...)
			}
			// restore the queue
			proc.queue = keep
			afterFirst = true
			// exit the loop if nothing was attempted.
			if lastRetry == nil {
				break
			}
		}
	}
	return
}
