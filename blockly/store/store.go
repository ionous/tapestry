package store

import (
	"encoding/json"
	"log"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

type TypeCreator interface {
	//HasType(string) bool
	// ex. test_scene
	// fix? _test_scene_stack
	NewType(string) (interface{}, bool)
}

// assumes dst is something similar to story file
func Decode(dst jsn.Marshalee, reg TypeCreator, msg json.RawMessage) (err error) {
	// we first unmarshal into a structure we can poke around in
	// a "streaming" decoder, would be a bit trickier to write.
	var bff BlockFile
	if e := json.Unmarshal(msg, &bff); e != nil {
		err = e
	} else if top, ok := bff.FindFirst("story_file"); !ok {
		err = errutil.New("couldnt find story file in block file")
	} else {
		dec := chart.MakeDecoder()
		err = dec.Marshal(dst, NewTopBlock(&dec, reg, top))
	}
	return
}

func NewTopBlock(m *chart.Machine, reg TypeCreator, bff *BlockInfo) *chart.StateMix {
	return &chart.StateMix{
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			// bff.type should equal flow.GetType()
			m.PushState(newInnerBlock(m, reg, flow, bff))
			return true
		},
		OnCommit: func(interface{}) {
			// blk.R(close)
			// m.FinishState(nil)
		},
	}
}

func newInnerBlock(m *chart.Machine, reg TypeCreator, flow jsn.FlowBlock, bff *BlockInfo) *chart.StateMix {
	var term string
	return &chart.StateMix{
		// a member of the dl, which might exist in bff.inputs, .fields, or .next;
		// it depends on what the next call is.
		OnKey: func(_ string, field string) (noerr error) {
			term = field[1:] // ex. StoryFile_Field_StoryLines = "$STORY_LINES"
			return
		},
		// a member that is a flow.
		// value lives in inputs
		OnMap: func(typeName string, flow jsn.FlowBlock) bool {
			// new inner block...
			return false
		},
		// a value that fills a slot; this will be an input
		OnSlot: func(string, jsn.SlotBlock) (alwaysTrue bool) {
			return false
		},
		// a member that's a swap will always be an input
		// for simple values ( strs in swaps ) there will be a faux for that input
		OnSwap: func(_ string, swap jsn.SwapBlock) (alwaysTrue bool) {
			return false
		},
		// a member that repeats:
		// could be a slice of a specific flow of a series of slots ( of numbered inputs ),
		// a stack ( starting in input, continuing in next ), a list ( of numbered fields )
		// we could return true/false depending on whether the block file has data
		// but first we have to find out where that data lives
		OnRepeat: func(typeName string, outBlocks jsn.SliceBlock) (okay bool) {
			if stack, ok := bff.Inputs[term]; ok {
				// might be nicer if count could grow, rather than counting in advance
				// could also sink them into a flat list as we count.
				cnt := 1 + stack.CountNext() // all of the next blocks connected to the input, plus the input block itself.
				outBlocks.SetSize(cnt)
				m.PushState(newStackReader(m, reg, stack.BlockInfo))

			} else {
				if i, cnt := bff.CountField(term); cnt > 0 {
					outBlocks.SetSize(cnt)
					m.PushState(newListReader(m, bff, i, i+cnt))
					okay = true
				}
			}
			return
		},
		// simple values live in bff.fields
		OnValue: func(_ string, pv interface{}) (err error) {
			if _, ok := bff.Fields.Find(term); !ok {
				err = jsn.Missing
			} else {
				err = errutil.New("not implemented")
			}
			return
		},
		// end of the dl structure
		OnEnd: func() {
		},
	}
}

// a stack is a repeating slot
// we expect to get OnMap/OnEnd for every element
func newStackReader(m *chart.Machine, reg TypeCreator, next *BlockInfo) *chart.StateMix {
	return &chart.StateMix{
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			if next != nil {
				if i, ok := reg.NewType(typeName); !ok {
					log.Println("couldn't create", typeName)
				} else if !flow.SetFlow(i) {
					log.Printf("couldn't set flow %T", i)
				} else {
					m.PushState(newInnerBlock(m, reg, flow, next))
					okay = true
				}
			}
			return
		},
		OnCommit: func(interface{}) {
			// advance to the next incoming data.
			if next = next.Next.BlockInfo; next == nil {
				m.FinishState(true)
			}
		},
	}
}

func newListReader(m *chart.Machine, b *BlockInfo, idx, end int) *chart.StateMix {
	return &chart.StateMix{
		OnValue: func(n string, pv interface{}) (err error) {
			field := b.Fields[idx]
			if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
				var i interface{}
				if e := json.Unmarshal(field.Msg, &i); e != nil {
					err = e
				} else if !el.SetValue(i) {
					err = errutil.New("couldnt set value", i)
				}
			} else {
				if e := json.Unmarshal(field.Msg, pv); e != nil {
					err = e // couldnt unmarshal directly into the target value
				}
			}
			if idx = idx + 1; idx >= end {
				m.FinishState(true)
			}
			return
		},
	}
}
