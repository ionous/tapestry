package story

import (
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// StoryStatement - a marker interface for commands which produce facts about the game world.
type StoryStatement interface {
	Weave(*weave.Catalog) error
}

// backcompat
func (op *StoryFile) Weave(cat *weave.Catalog) error {
	return Weave(cat, op.Statements)
}

// visits each statement calling PreImport and PostImport
// to transform the statements; then calls Weave on each.
func Weave(cat *weave.Catalog, all []StoryStatement) (err error) {
	evts := inspect.Callbacks{
		// given a slot, replace its command using PreImport or PostImport
		// and, walk the contents of its (replaced) for additional pre or post imports.
		// a command usually would only implement either Pre or Post ( or neither. )
		OnSlot: func(slot inspect.It) (err error) {
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
		OnEnd: func(slot inspect.It) (err error) {
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

const activityDepth = "activityDepth"

// was for comment logging; remove?
func updateActivityDepth(cat *weave.Catalog, t *typeinfo.Slot, inc int) {
	if t == &rtti.Zt_Execute {
		cat.Env.Inc(activityDepth, inc)
	}
}
