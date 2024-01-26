package story

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
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
	evts := inspect.Callbacks{
		// given a slot, replace its command using PreImport or PostImport
		// and, walk the contents of its (replaced) for additional pre or post imports.
		// a command usually would only implement either Pre or Post ( or neither. )
		OnSlot: func(slot inspect.Iter) (err error) {
			t := slot.TypeInfo().(*typeinfo.Slot)
			updateActivityDepth(cat, t, 1)
			var tgt PreImport
			if ok := slot.GetSlot(&tgt); ok {
				if rep, e := tgt.PreImport(cat); e != nil {
					err = errutil.New(e, "failed to create pre import")
				} else if rep != nil {
					if !slot.SetSlot(rep) {
						err = errutil.New("couldnt assign pre import")
					}
				}
			}
			return
		},
		OnEnd: func(slot inspect.Iter) (err error) {
			if t, ok := slot.TypeInfo().(*typeinfo.Slot); ok && !slot.Repeating() {
				var tgt PostImport
				if ok := slot.GetSlot(&tgt); ok {
					if rep, e := tgt.PostImport(cat); e != nil {
						err = errutil.New(e, "failed to create post import")
					} else if rep != nil {
						if !slot.SetSlot(rep) {
							err = errutil.New("couldnt assign post import")
						}
					}
				}
				updateActivityDepth(cat, t, -1)
			}
			return
		},
	}
	//
	currentCatalog = cat
	for _, el := range all {
		slot := StoryStatement_Slot{Value: el}
		if e := inspect.Visit(&slot, &evts); e != nil {
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

// fix: for comment logging; remove?
func updateActivityDepth(cat *weave.Catalog, t *typeinfo.Slot, inc int) {
	if t == &rtti.Zt_Execute {
		cat.Env.Inc(activityDepth, inc)
	}
}
