package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
)

// create the "event object" containing the event name, target, interrupt and cancel status
func newEventRecord(run rt.Runtime, name string, tgt rt.Value) (ret *rt.Record, err error) {
	if eventFields == nil {
		eventFields = make([]rt.Field, event.NumFields)
		for i := 0; i < event.NumFields; i++ {
			f := event.Field(i)
			eventFields[i] = rt.Field{
				Name:     f.String(),
				Affinity: f.Affine(),
			}
		}
	}
	out := rt.NewAnonymousRecord(eventFields)
	if e := out.SetIndexedField(event.Name.Index(), rt.StringOf(name)); e != nil {
		err = e
	} else if e := out.SetIndexedField(event.Target.Index(), tgt); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

// cache of the event object fields as described in package event
var eventFields []rt.Field
