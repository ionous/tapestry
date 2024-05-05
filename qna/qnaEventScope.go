package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/event"
)

// create the "event object" containing the event name, target, interrupt and cancel status
func newEventRecord(run rt.Runtime, name string, tgt rt.Value) (ret *rt.Record, err error) {
	out := rt.NewRecord(&eventKind)
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
var eventKind rt.Kind

func init() {
	fields := make([]rt.Field, event.NumFields)
	for i := 0; i < event.NumFields; i++ {
		f := event.Field(i)
		fields[i] = rt.Field{
			Name:     f.String(),
			Affinity: f.Affine(),
			Type:     f.Type(),
		}
	}
	traits := make([]string, event.NumStatus)
	for i := 0; i < event.NumStatus; i++ {
		state := event.CancellationStatus(i)
		traits[i] = state.String()
	}
	eventKind = rt.Kind{
		Fields: fields,
		Aspects: []rt.Aspect{{
			Name:   event.Status.String(),
			Traits: traits,
		}},
	}
}
