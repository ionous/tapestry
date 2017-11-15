package event

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/class"
	"github.com/ionous/iffy/rt"
)

// Map event id to listener list
type EventMap map[ident.Id]EventListeners

// we have two kinds of listeners, class listeners, and object listeners.
type EventListeners struct {
	Classes, Objects PhaseMap
}

// id ( either class or object ) to list
type PhaseMap map[ident.Id]PhaseList

// PhaseList contains lists of capture and bubble handlers.
type PhaseList [ListenerTypes][]Handler

type Handler struct {
	Options
	Exec rt.Execute
}

type Target struct {
	obj      rt.Object
	cls      rt.Class
	handlers PhaseList
}

// CollectAncestors to create targets from the parents of the passed object. The  target order is: instance's parent, parent classes, container instance, repeat.
func (els EventListeners) CollectAncestors(run rt.Runtime, obj rt.Object) (ret []Target, err error) {
	if at, e := run.GetAncestors(run, obj); e != nil {
		err = e
	} else {
		var tgt []Target
		for at.HasNext() {
			if obj, e := at.GetObject(); e != nil {
				err = e
				break
			} else {
				tgt = els.CollectTargets(obj, tgt)
			}
		}
		if err == nil {
			ret = tgt
		}
	}
	return
}

func (els EventListeners) CollectTargets(obj rt.Object, tgt []Target) []Target {
	// check instance listeners
	if ls, ok := els.Objects[obj.Id()]; ok {
		tgt = append(tgt, Target{obj: obj, handlers: ls})
	}
	// check class listeners
	for cls := obj.Type(); ; {
		if ls, ok := els.Classes[class.Id(cls)]; ok {
			tgt = append(tgt, Target{obj: obj, cls: cls, handlers: ls})
		}
		// move to parent class
		if next, ok := class.Parent(cls); !ok {
			break
		} else {
			cls = next
		}
	}
	return tgt
}
