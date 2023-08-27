package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
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
	if e := importStory(cat, tgt); e != nil {
		err = e
	} else {
		err = WeaveStatements(cat, tgt.StoryStatements)
	}
	return err
}

func WeaveStatements(cat *weave.Catalog, all []StoryStatement) (err error) {
	for _, el := range all {
		if e := el.Weave(cat); e != nil {
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

// post-processing hooks:
// after the story has been read we run an encoder on it to visit every node.
func importStory(cat *weave.Catalog, tgt jsn.Marshalee) error {
	ts := chart.MakeEncoder()
	return ts.Marshal(tgt, chart.Map(&ts, chart.BlockMap{
		rt.Execute_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, _ interface{}) (_ error) {
				cat.Env.Inc(activityDepth, 1)
				return
			},
			chart.BlockEnd: func(b jsn.Block, _ interface{}) (_ error) {
				cat.Env.Inc(activityDepth, -1)
				return
			},
		},
		assign.CallPattern_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); !ok {
					err = errutil.Fmt("trying to import something other than a flow")
				} else if _, ok := flow.GetFlow().(*assign.CallPattern); !ok {
					err = errutil.Fmt("trying to import something other unexpected")
				} else {
					// k.WriteEphemera(ImportCall(op))
				}
				return
			},
		},
		grammar.Action_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); !ok {
					err = errutil.Fmt("trying to import something other than a flow")
				} else if op, ok := flow.GetFlow().(*grammar.Action); !ok {
					err = errutil.Fmt("trying to import something other unexpected")
				} else {
					err = importAction(cat, op)
				}
				return
			},
		},
		chart.OtherBlocks: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if slot, ok := b.(jsn.SlotBlock); ok {
					if slat, ok := slot.GetSlot(); !ok {
						err = jsn.Missing
					} else if tgt, ok := slat.(PreImport); ok {
						if rep, e := tgt.PreImport(cat); e != nil {
							err = errutil.New(e, "failed to create replacement")
						} else if rep != nil && !slot.SetSlot(rep) {
							err = errutil.Fmt("failed to set slot %T with replacement %T", slot, rep)
						}
					}
				}
				return
			},
			chart.BlockEnd: func(b jsn.Block, v interface{}) (err error) {
				if slot, ok := b.(jsn.SlotBlock); ok {
					if slat, ok := slot.GetSlot(); !ok {
						err = jsn.Missing
					} else if tgt, ok := slat.(PostImport); ok {
						if rep, e := tgt.PostImport(cat); e != nil {
							err = errutil.New(e, "failed to create replacement")
						} else if rep != nil && !slot.SetSlot(rep) {
							err = errutil.Fmt("failed to set slot %T with replacement %T", slot, rep)
						}
					}
				}
				return
			},
		},
	}))
}
