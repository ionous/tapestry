package jesstest

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Process func(weaver.Weaves, rt.Runtime) error

// used internally to generate matched phrases in a good order.
type ProcessingList struct {
	pending   [weaver.NumPhases][]Process
	lastPhase weaver.Phase
}

func (m *ProcessingList) Schedule(z weaver.Phase, p func(weaver.Weaves, rt.Runtime) error) (err error) {
	return m.SchedulePos(compact.Source{}, z, p)
}

// add to the post processing list.
func (m *ProcessingList) SchedulePos(at compact.Source, z weaver.Phase, p func(weaver.Weaves, rt.Runtime) error) (err error) {
	if next := m.lastPhase + 1; z < next {
		err = fmt.Errorf("scheduled %s(%d) while next phase is %s(%d)", z, z, next, next)
		panic(err)
	} else {
		m.pending[z] = append(m.pending[z], p)
	}
	return
}

// run through the passed processing list one time
// return the number of callbacks which completed;
// return weaver.Missing if any failed to complete.
func (m *ProcessingList) UpdatePhase(z weaver.Phase, w weaver.Weaves, run rt.Runtime) (ret int, err error) {
	m.lastPhase = z //
	if prev := m.pending[z]; len(prev) > 0 {
		if next, e := updatePhase(w, run, prev); e != nil && !errors.Is(e, weaver.ErrMissing) {
			err = e
		} else {
			ret = len(prev) - len(next) // how many were processed
			m.pending[z] = next
			err = e // return what was missing ( if anything )
		}
	}
	return
}

// run through the passed processing list one time
// return weaver.Missing if any of them failed to process
func updatePhase(w weaver.Weaves, run rt.Runtime, pending []Process) (ret []Process, err error) {
	var retry int
	var firstMissing error
	for _, callback := range pending {
		if e := callback(w, run); e != nil {
			if !errors.Is(e, weaver.ErrMissing) {
				err = e
				break // return any critical errors
			} else {
				pending[retry] = callback
				retry++ // remember any callbacks that need to try again
				if firstMissing == nil {
					firstMissing = e // remember the first reason
				}
			}
		}
	}
	if err == nil {
		ret = pending[:retry]
		err = firstMissing
	}
	return
}
