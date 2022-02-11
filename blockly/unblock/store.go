package unblock

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
	var bff File
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

func NewTopBlock(m *chart.Machine, reg TypeCreator, bff *Info) *chart.StateMix {
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

func newInnerBlock(m *chart.Machine, reg TypeCreator, flow jsn.FlowBlock, bff *Info) *chart.StateMix {
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
				m.PushState(newStackReader(m, reg, stack.Info))
				okay = true
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
			if field, ok := bff.Fields.Find(term); !ok {
				err = jsn.Missing
			} else {
				err = storeValue(pv, field.Msg)
			}
			return
		},
		// end of the dl structure
		OnEnd: func() {
			m.FinishState(true)
		},
	}
}

// a stack is a repeating slot
// we expect to get OnMap/OnEnd for every element
func newStackReader(m *chart.Machine, reg TypeCreator, next *Info) *chart.StateMix {
	return &chart.StateMix{
		// create the value for the slot
		OnSlot: func(_ string, slot jsn.SlotBlock) (okay bool) {
			// the typename we want is (munged) in the block file
			if typeName, ok := unstackName(next.Type); !ok {
				log.Println("couldnt unstack", next.Type)
			} else if i, ok := reg.NewType(typeName); !ok {
				log.Println("couldn't create", typeName)
			} else if !slot.SetSlot(i) {
				log.Printf("couldn't set flow %T", i)
			} else {
				okay = true
			}
			return
		},
		// happens after OnSlot for every block of data in the stack
		OnMap: func(_ string, flow jsn.FlowBlock) (alwaysTrue bool) {
			m.PushState(newInnerBlock(m, reg, flow, next))
			return true
		},
		// called after each the map's inner block completes
		// if we are out of data, we end the stack reader
		OnCommit: func(interface{}) {
			// advance the function level's next pointer.
			if outer := next.Next; outer != nil {
				next = outer.Info
			} else {
				m.FinishState(true)
			}
		},
	}
}

func unstackName(n string) (ret string, okay bool) {
	const suffix = "_stack"
	if cnt := len(n); cnt > len(suffix) && n[0] == '_' && n[cnt-len(suffix):] == suffix {
		ret = n[1 : cnt-len(suffix)]
		okay = true
	}
	return
}

func newListReader(m *chart.Machine, b *Info, idx, end int) *chart.StateMix {
	return &chart.StateMix{
		OnValue: func(n string, pv interface{}) (err error) {
			field := b.Fields[idx]
			if e := storeValue(pv, field.Msg); e != nil {
				err = e
			} else if idx = idx + 1; idx >= end {
				m.FinishState(true)
			}
			return
		},
	}
}

func storeValue(pv interface{}, msg json.RawMessage) (err error) {
	if el, ok := pv.(interface{ SetValue(interface{}) bool }); ok {
		var i interface{}
		if e := json.Unmarshal(msg, &i); e != nil {
			err = e
		} else if !el.SetValue(i) {
			err = errutil.New("couldnt set value", i)
		}
	} else {
		if e := json.Unmarshal(msg, pv); e != nil {
			err = e // couldnt unmarshal directly into the target value
		}
	}
	return
}
