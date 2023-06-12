package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
	"github.com/ionous/errutil"
)

// StoryStatement - a marker interface for commands which produce facts about the game world.
type StoryStatement interface {
	Schedule(k *weave.Catalog) error
}

func ImportStory(cat *weave.Catalog, path string, tgt *StoryFile) (err error) {
	cat.SetSource(path)
	if e := importStory(cat, tgt); e != nil {
		err = e
	} else {
		err = ScheduleStatements(cat, tgt.StoryStatements)
	}
	return err
}

func ScheduleStatements(cat *weave.Catalog, all []StoryStatement) (err error) {
	for _, el := range all {
		if e := el.Schedule(cat); e != nil {
			err = e
			break
		}
	}
	return
}

// post-processing hooks
func importStory(k *weave.Catalog, tgt jsn.Marshalee) error {
	ts := chart.MakeEncoder()
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
					} else if tgt, ok := slat.(PreImport); ok {
						if rep, e := tgt.PreImport(k); e != nil {
							err = errutil.New(e, "failed to create replacement")
						} else if rep != nil && !slot.SetSlot(rep) {
							err = errutil.New("failed to set replacement")
						}
					}
				}
				return
			},
		},
	}))
}
