package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// create the "event object" containing the event name, target, interrupt and cancel status
func newEventRecord(run rt.Runtime, name string, tgt g.Value) (ret *g.Record, err error) {
	if eventFields == nil {
		eventFields = make([]g.Field, event.NumFields)
		for i := 0; i < event.NumFields; i++ {
			f := event.Field(i)
			eventFields[i] = g.Field{
				Name:     f.String(),
				Affinity: f.Affine(),
			}
		}
	}
	out := g.NewAnonymousRecord(run, eventFields)
	if e := out.SetIndexedField(event.Name.Index(), g.StringOf(name)); e != nil {
		err = e
	} else if e := out.SetIndexedField(event.Target.Index(), tgt); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

// cache of the event object fields as described in package event
var eventFields []g.Field
