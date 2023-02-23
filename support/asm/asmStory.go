package asm

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

func ImportStory(k *imp.Importer, path string, tgt jsn.Marshalee) (err error) {
	k.SetSource(path)
	return importStory(k, tgt)
}

// post-processing hooks
func importStory(k *imp.Importer, tgt jsn.Marshalee) error {
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
				} else if op, ok := flow.GetFlow().(*assign.CallPattern); !ok {
					err = errutil.Fmt("trying to import something other than a response")
				} else {
					k.WriteEphemera(story.ImportCall(op))
				}
				return
			},
		},
		core.Response_Type: chart.KeyMap{
			chart.BlockStart: func(b jsn.Block, v interface{}) (err error) {
				if flow, ok := b.(jsn.FlowBlock); !ok {
					err = errutil.Fmt("trying to import something other than a flow")
				} else if resp, ok := flow.GetFlow().(*core.Response); !ok {
					err = errutil.Fmt("trying to import something other than a response")
				} else {
					k.WriteEphemera(&eph.EphKinds{
						Kinds: kindsOf.Response.String(),
						Contain: []eph.EphParams{{
							Affinity:  eph.Affinity{eph.Affinity_Text},
							Name:      resp.Name,
							Initially: &assign.FromText{Value: resp.Text},
						}},
					})
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
				// sometimes we also get slice blocks...
				if slot, ok := b.(jsn.SlotBlock); ok {
					// sometimes we get empty slots...
					if val, ok := slot.GetSlot(); ok {
						if stmt, ok := val.(imp.PostImport); ok {
							err = stmt.PostImport(k)
						}
					}
				}
				return
			},
		},
	}))
}
