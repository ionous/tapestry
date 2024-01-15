package story

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// StoryStatement - a marker interface for commands which produce facts about the game world.
type StoryStatement interface {
	Weave(*weave.Catalog) error
}

// hacky: go interfaces arent vtables;
// so when a runtime helper implements the rt interface:
// it has no access to the full implementation of the interface.
// meaning inside rule application the importer isnt accessible via casting.
// we'd need a context maybe ( ex. pass an interface{} through Options );
// a global is fine for now.
var currentCatalog *weave.Catalog

func ImportStory(cat *weave.Catalog, path string, tgt *StoryFile) (err error) {
	cat.SetSource(path)
	return WeaveStatements(cat, tgt.StoryStatements)
}

func WeaveStatements(cat *weave.Catalog, all []StoryStatement) (err error) {
	currentCatalog = cat
	slice := r.ValueOf(all)
	evts := walk.Events{
		BeforeSlot: beforeSlot,
		AfterSlot:  afterSlot,
	}
	for i, el := range all {
		w := walk.Walk(slice.Index(i))
		if e := walk.VisitSlot(w, evts); e != nil {
			err = e
			break
		} else if e := el.Weave(cat); e != nil {
			err = e
			break
		}
	}
	return
}

// transform a story statement's execution ( ex. during a macro )
// into a weave so that it can generate facts for the database
// expects that the runtime is the importer's own runtime.
// ( as opposed to the story's playtime. )
func Weave(run rt.Runtime, op StoryStatement) (err error) {
	if cat := currentCatalog; cat.Runtime() != run {
		err = errutil.Fmt("mismatched runtimes?")
	} else {
		err = op.Weave(cat)
	}
	return
}

// given a slot, replace its command using PreImport or PostImport
// and, walk the contents of its (replaced) for additional pre or post imports.
// a command usually would only implement either Pre or Post ( or neither. )
func beforeSlot(slot walk.Walker) (_ walk.Events, err error) {
	cat := currentCatalog
	if i := unpackSlot(slot); i != nil {
		updateActivityDepth(cat, i, 1)
		//
		if tgt, ok := i.(PreImport); ok {
			if rep, e := tgt.PreImport(cat); e != nil {
				err = errutil.New(e, "failed to create pre replacement")
			} else if rep != nil {
				// fix: a more direct setValue?
				slot.Value().Set(r.ValueOf(rep))
			}
		}
	}
	return
}

func afterSlot(slot walk.Walker) (_ walk.Events, err error) {
	cat := currentCatalog
	if i := unpackSlot(slot); i != nil {
		if tgt, ok := i.(PostImport); ok {
			if rep, e := tgt.PostImport(cat); e != nil {
				err = errutil.New(e, "failed to create post replacement")
			} else if rep != nil {
				slot.Value().Set(r.ValueOf(rep))
			}
		}
		updateActivityDepth(cat, i, -1)
	}
	return
}

// fix: for comment logging; remove?
func updateActivityDepth(cat *weave.Catalog, i any, inc int) {
	// fix:
	// this used to switch on the *slot*
	// now its switching on the interface
	// *in* that slot. getting that zero value for a switch is a bit of a pain.

	// if _, ok := i.(rt.Execute); ok {
	// 	cat.Env.Inc(activityDepth, inc)
	// }
}

// kind of weird: the value of the slot is a pointer;
// the value of the flow is a struct; we need the pointer.

func unpackSlot(slot walk.Walker) (ret interface{}) {
	if v := slot.Value(); v.IsValid() {
		ret = v.Interface()
	}
	return
}
