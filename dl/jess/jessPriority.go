package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Phase = mdl.Phase
type Process func() error

type Context struct {
	Query
	Registrar
	ProcessingList
}

func NewContext(q Query, rar Registrar) *Context {
	return &Context{Query: q, Registrar: rar}
}

// used internally to generate matched phrases in a good order.
type ProcessingList struct {
	pending   [mdl.NumPhases][]Process
	lastPhase mdl.Phase
}

// add to the post processing list.
func (m *ProcessingList) PostProcess(z Phase, p Process) (err error) {
	if next := m.lastPhase + 1; z < next {
		err = fmt.Errorf("scheduled %s(%d) while next phase is %s(%d)", z, z, next, next)
	} else {
		// fix: here out of curiosity; remove.
		if z == next {
			// log.Printf("scheduled for current phase %s", z)
		}
		m.pending[z] = append(m.pending[z], p)
	}
	return
}

// can return with:
// everything in this phase competed ( err is nil )
// nothing was able to run ( err is mdl.Missing )
// some other error
func (m *ProcessingList) UpdatePhase(z Phase) (err error) {
	m.lastPhase = z //
	for prev := m.pending[z]; len(prev) > 0; m.pending[z] = prev {
		// run the set of post process callbacks
		if next, e := updatePhase(z, prev); e != nil && !errors.Is(e, mdl.Missing) {
			err = e
			break
		} else if e != nil && len(next) == len(prev) {
			// return missing only if nothing in the list moved forward
			err = e
			break
		} else {
			// otherwise update the list
			prev = next
		}
	}
	return
}

// run through the passed processing list one time
// return mdl.Missing if any of them failed to process
func updatePhase(z Phase, pending []Process) (ret []Process, err error) {
	var retry int
	var firstMissing error
	for _, callback := range pending {
		if e := callback(); e != nil {
			if !errors.Is(e, mdl.Missing) {
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
