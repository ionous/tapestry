package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// StoryStatement - a marker interface for commands which produce facts about the game world.
type StoryStatement interface {
	PostImport(k *imp.Importer) error
}

func ImportStory(k *imp.Importer, path string, tgt *StoryFile) error {
	k.SetSource(path)
	return importStory(k, tgt)
}

// post-processing hooks
func importStory(k *imp.Importer, tgt jsn.Marshalee) error {
	ts := chart.MakeEncoder()
	macroDepth := 0
	return ts.Marshal(tgt, chart.Map(&ts, chart.BlockMap{
		rt.Execute_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, _ interface{}) (err error) {
				k.Env().ActivityDepth++
				return
			},
			chart.BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				k.Env().ActivityDepth--
				return
			},
		},
		DefineMacro_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, _ interface{}) (err error) {
				macroDepth++
				return
			},
			chart.BlockEnd: func(b jsn.Block, _ interface{}) (err error) {
				macroDepth--
				return
			},
		},
		assign.CallPattern_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); !ok {
					err = errutil.Fmt("trying to import something other than a flow")
				} else if _, ok := flow.GetFlow().(*assign.CallPattern); !ok {
					err = errutil.Fmt("trying to import something other than a response")
				} else {
					// k.WriteEphemera(ImportCall(op))
				}
				return
			},
		},
		chart.OtherBlocks: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if slot, ok := b.(jsn.SlotBlock); ok {
					if slat, ok := slot.GetSlot(); !ok {
						err = jsn.Missing
					} else if tgt, ok := slat.(imp.PreImport); ok {
						if rep, e := tgt.PreImport(k); e != nil {
							err = errutil.New(e, "failed to create replacement")
						} else if rep != tgt && !slot.SetSlot(rep) {
							err = errutil.New("failed to set replacement")
						}
					}
				}
				return
			},
			chart.BlockEnd: func(b jsn.Block, v interface{}) (err error) {
				// ignore the contents of macros, because those will get executed by the macro
				// tbd: macros in macros, and macros containing pre-import statements.
				if macroDepth == 0 {
					// sometimes we also get slice blocks...
					if slot, ok := b.(jsn.SlotBlock); ok {
						// sometimes we get empty slots...
						if val, ok := slot.GetSlot(); ok {
							if stmt, ok := val.(imp.PostImport); ok {
								err = stmt.PostImport(k)
							}
						}
					}
				}
				return
			},
		},
	}))
}
