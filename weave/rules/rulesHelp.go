package rules

import (
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func CanFilterActor(k *g.Kind) (ret bool) {
	// by default: all event handlers are filtered to the player and the innermost target.
	eventLike := k.Implements(kindsOf.Event.String()) || k.Implements(kindsOf.Action.String())
	// if the focus of the event involves an actor;
	// then we automatically filter for the player
	if k.NumField() > 0 && eventLike {
		if f := k.Field(0); f.Type == event.Actors && f.Name == event.Actor {
			ret = true
		}
	}
	return
}
