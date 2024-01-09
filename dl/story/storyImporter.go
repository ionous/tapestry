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
	currentCatalog = cat
	cat.SetSource(path)
	return WeaveStatements(cat, tgt.StoryStatements)
}

func WeaveStatements(cat *weave.Catalog, all []StoryStatement) (err error) {
	slice := r.ValueOf(all)
	for i, el := range all {
		if e := rewriteSlot(cat, walk.Walk(slice.Index(i))); e != nil {
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
func rewriteSlot(cat *weave.Catalog, slot walk.Walker) (err error) {
	if flow, ok := unpackSlot(slot); ok {
		// kind of weird: the value of the slot is a pointer;
		// the value of the flow is a struct; we need the pointer.
		i := slot.Value().Interface()
		if tgt, ok := i.(PreImport); ok {
			if rep, e := tgt.PreImport(cat); e != nil {
				err = errutil.New(e, "failed to create pre replacement")
			} else if rep != nil {
				// fix: a more direct setValue?
				slot.Value().Set(r.ValueOf(rep))
			}
		}
		//
		if e := rewriteFlow(cat, flow); e != nil {
			err = e
		} else if tgt, ok := i.(PostImport); ok {
			if rep, e := tgt.PostImport(cat); e != nil {
				err = errutil.New(e, "failed to create post replacement")
			} else if rep != nil {
				slot.Value().Set(r.ValueOf(rep))
			}
		}
	}
	return
}

// a flow cant be swapped out the way a slot can,
// however a slot in a glow might be able to change.
func rewriteFlow(cat *weave.Catalog, flow walk.Walker) (err error) {
	for flow.Next() {
		switch f := flow.Field(); f.SpecType() {
		case walk.Slot:
			if !f.Repeats() {
				err = rewriteSlot(cat, flow.Walk())
			} else {
				for it := flow.Walk(); it.Next(); {
					err = rewriteSlot(cat, it.Walk())
				}
			}
		case walk.Flow:
			if !f.Repeats() {
				err = rewriteFlow(cat, flow.Walk())
			} else {
				for it := flow.Walk(); it.Next(); {
					err = rewriteFlow(cat, it.Walk())
				}
			}
		}
	}
	return
}

// given a slot, return its flow ( false if the slot was empty )
func unpackSlot(slot walk.Walker) (ret walk.Walker, okay bool) {
	if slot.Next() { // next moves the focus to the command if any
		ret, okay = slot.Walk(), true // walk returns the command as a container
	}
	return
}
