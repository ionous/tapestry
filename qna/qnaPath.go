package qna

import (
	"slices"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
)

// assumes the passed tgt is a text value referring to the object ( actor or other noun ) targeted by an event.
func (run *Runner) newPathForTarget(tgt rt.Value) (ret eventPath, err error) {
	if els, e := run.Call(event.CapturePattern, affine.TextList, nil, []rt.Value{tgt}); e != nil {
		err = e
	} else {
		ret = eventPath{els.Strings(), event.Bubbles}
	}
	return
}

type eventPath struct {
	path  []string
	order event.Flow
}

// for a given event
func (p *eventPath) slice(evt event.Phase) (ret []string) {
	if flow := evt.Flow(); flow == event.Targets {
		start, end := 0, len(p.path)
		if p.order == event.Bubbles {
			end = start + 1 // use only the front
		} else {
			start = end - 1 // use only the end
		}
		ret = p.path[start:end]
	} else {
		if flow != p.order {
			slices.Reverse(p.path)
			p.order = flow
		}
		ret = p.path
	}
	return
}
