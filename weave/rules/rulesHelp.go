package rules

import (
	"git.sr.ht/~ionous/tapestry/rt/event"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func CanFilterActor(k *g.Kind) (ret bool) {
	// if the focus of an event involves an actor;
	// then we automatically filter for the player
	if k.NumField() > 0 && k.Implements(kindsOf.Action.String()) {
		if f := k.Field(0); f.Type == event.Actors && f.Name == event.Actor {
			ret = true
		}
	}
	return
}
