package unblock

import (
	"encoding/json"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

type TypeCreator interface {
	NewType(string) (interface{}, bool)
}

// where topBlock is the expected topblock type in the file.... ex. story_file
func Decode(dst jsn.Marshalee, topBlock string, reg TypeCreator, msg json.RawMessage) (err error) {
	// we first unmarshal into a structure we can poke around in
	// a "streaming" decoder, would be a bit trickier to write.
	var bff File
	if e := json.Unmarshal(msg, &bff); e != nil {
		err = e
	} else if top, ok := bff.FindFirst(topBlock); !ok {
		err = errutil.New("couldnt find story file in block file")
	} else {
		dec := chart.MakeDecoder()
		err = dec.Marshal(dst, NewBlock(&dec, reg, top))
	}
	return
}

func NewBlock(m *chart.Machine, reg TypeCreator, bff *BlockInfo) *chart.StateMix {
	return &chart.StateMix{
		Name: "Block:" + bff.Id,
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			// bff.type should equal flow.GetType()
			m.PushState(newInnerBlock(m, reg, flow, bff))
			return true
		},
		OnEnd: func() {
			m.FinishState(true)
		},
	}
}

// since this is using the standard decoder, deserialization is based on the go type.
func newInnerBlock(m *chart.Machine, reg TypeCreator, flow jsn.FlowBlock, bff *BlockInfo) *chart.StateMix {
	var termName string // pulled from the golang key
	if m.Comment != nil {
		if c := bff.Icons.Comment; c != nil {
			*m.Comment = c.Text
		}
	}
	return &chart.StateMix{
		Name: "InnerBlock:" + bff.Id,
		// a member of the dl whose value in blockly might exist in bff.inputs, .fields, or .next;
		// it depends on the next OnXxx: callback.
		OnKey: func(_ string, field string) (noerr error) {
			// termName must match shapeWriter's shapeField.name() which is specified as upper(TermSpec.Field())
			// the code generator passes this callback constants for the field parameter ( via gomake's flow.tmpl )
			// ex. StoryFile_Field_StoryLines = term.Value() = "$" + upper(TermSpec.Field())
			termName = field[1:]
			return
		},
		// a member that is a flow; its value lives in an input.
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
				if input, e := bff.ReadInput(idx); e != nil {
					log.Println(e)
				} else {
					m.PushState(newInnerBlock(m, reg, flow, input.BlockInfo))
					okay = true
				}
			}
			return
		},
		// a member that fills a slot; its value lives in an input.
		OnSlot: func(_ string, slot jsn.SlotBlock) (okay bool) {
			if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
				if input, e := bff.ReadInput(idx); e != nil {
					log.Println(e)
				} else if e := fillSlot(reg, slot, input.Type); e != nil {
					log.Println(e)
				} else {
					m.PushState(NewBlock(m, reg, input.BlockInfo))
					okay = true
				}
			}
			return
		},
		// a member that's a swap uses both a field and an input.
		// for simple values ( strs in swaps ) there will be a faux block type for that input.
		OnSwap: func(_ string, swap jsn.SwapBlock) (okay bool) {
			// the field holds a combo box with swap's choice
			if idx := bff.Fields.FindIndex(termName); idx >= 0 {
				var choice string
				if e := storeValue(&choice, bff.Fields[idx].Msg); e != nil {
					log.Println(e)
				} else if !swap.SetSwap(choice) {
					log.Println("unexpected choice for swap", termName, choice)
				} else if idx := bff.Inputs.FindIndex(termName); idx >= 0 {
					// the input hold the value of the swap's block
					if input, e := bff.ReadInput(idx); e != nil {
						log.Println(e)
					} else {
						m.PushState(newSwapContents(m, reg, input.BlockInfo))
						okay = true
					}
				}
			}
			return
		},
		// a member that repeats:
		// could be a slice of a specific flow of a series of slots ( of numbered inputs ),
		// a stack ( starting in input, continuing in next ), a list ( of numbered fields )
		// we could return true/false depending on whether the block file has data
		// but first we have to find out where that data lives
		OnRepeat: func(typeName string, outBlocks jsn.SliceBlock) (okay bool) {
			if i, cnt := bff.CountFields(termName); cnt > 0 {
				outBlocks.SetSize(cnt)
				m.PushState(newListReader(m, bff, termName, i, i+cnt))
				okay = true
			} else if i, cnt := bff.CountInputs(termName); cnt > 0 {
				outBlocks.SetSize(cnt)
				m.PushState(newSeriesSlice(m, reg, bff, i, i+cnt))
				okay = true
			} else if at := bff.Inputs.FindIndex(termName); at >= 0 {
				if input, e := bff.ReadInput(at); e != nil {
					log.Println(e)
				} else {
					// might be nicer if count could grow, rather than counting in advance
					// could also sink them into a flat list as we count.
					cnt := 1 + input.CountNext() // all of the next blocks connected to the input, plus the input block itself.
					outBlocks.SetSize(cnt)
					m.PushState(newStackReader(m, reg, termName, typeName, input.BlockInfo))
					okay = true
				}
			}
			return
		},
		// simple values live in bff.fields
		OnValue: func(_ string, pv interface{}) (err error) {
			if field, ok := bff.Fields.Find(termName); !ok {
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
func newStackReader(m *chart.Machine, reg TypeCreator, termName, typeName string, next *BlockInfo) *chart.StateMix {
	return &chart.StateMix{
		Name: "StackReader:" + termName,
		// create the value for the slot
		OnSlot: func(_ string, slot jsn.SlotBlock) (okay bool) {
			// the typename we want is (munged) in the block file
			if typeName, ok := unstackName(next.Type); !ok {
				log.Println("couldnt unstack", next.Type)
			} else if e := fillSlot(reg, slot, typeName); e != nil {
				log.Println(e)
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
		// end of each slot
		OnEnd: func() {
			// advance the function level's next pointer.
			if outer := next.Next; outer != nil {
				next = outer.BlockInfo
			} else {
				// after we are out of stacked blocks:
				// wait for the end of the stack
				m.ChangeState(chart.NewBlockResult(m, true))
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

// read the insides of a swap:
// it could be a flow filling the input....
// or a fake block wrapping a primitive value ( a "standalone" )
func newSwapContents(m *chart.Machine, reg TypeCreator, bff *BlockInfo) *chart.StateMix {
	next := NewBlock(m, reg, bff)
	next.Name = "Swap:" + bff.Type
	next.OnValue = func(typeName string, pv interface{}) (err error) {
		// see: block.newSwap & shape.writeStandalone:
		// the fields name is the name of the ifspec type 0 )
		field := strings.ToUpper(typeName)
		if idx := bff.Fields.FindIndex(field); idx < 0 {
			err = jsn.Missing // the block might be missing, and that's okay.
		} else {
			field := bff.Fields[idx]
			err = storeValue(pv, field.Msg) // pv is the destination
		}
		// note: even for values, the state is still getting popped by the ending of the swap
		// ( via the handler added in NewBlock )
		return
	}
	return next
}

// we cant tell whether we have a series of repeating slots, or a slice of repeating flows
// not until we hear a request for either the first slot, or the first flow.
func newSeriesSlice(m *chart.Machine, reg TypeCreator, bff *BlockInfo, idx, end int) *chart.StateMix {
	id := bff.Id + " " + bff.Inputs[idx].Key
	return &chart.StateMix{
		Name: "SeriesOrSlice:" + id,
		OnSlot: func(n string, slot jsn.SlotBlock) bool {
			state := newSeries(m, reg, bff, idx, end)
			m.ChangeState(state)
			return state.OnSlot(n, slot)
		},
		OnMap: func(n string, flow jsn.FlowBlock) bool {
			state := newSlice(m, reg, bff, idx, end)
			m.ChangeState(state)
			return state.OnMap(n, flow)
		},
	}
}

// repeating non-stacking slots
// idx is the index of the first input, end is one beyond the last matching input
func newSeries(m *chart.Machine, reg TypeCreator, bff *BlockInfo, idx, end int) *chart.StateMix {
	var next *BlockInfo // the block connected to the current input
	id := bff.Id + " " + bff.Inputs[idx].Key
	return &chart.StateMix{
		Name: "Series:" + id,
		// create the value for the slot
		OnSlot: func(_ string, slot jsn.SlotBlock) (okay bool) {
			if input, e := bff.ReadInput(idx); e != nil {
				log.Println(e)
			} else if e := fillSlot(reg, slot, input.Type); e != nil {
				log.Println(e)
			} else {
				next = input.BlockInfo
				okay = true
			}
			return
		},
		// happens after OnSlot for every block of data in the stack
		OnMap: func(_ string, flow jsn.FlowBlock) (alwaysTrue bool) {
			m.PushState(newInnerBlock(m, reg, flow, next))
			return true
		},
		// end of each slot
		OnEnd: func() {
			if idx++; idx >= end {
				// after we are out of inputs:
				// wait for the end of the series
				m.ChangeState(chart.NewBlockResult(m, true))
			}
		},
	}
}

// repeating flows
// idx is the index of the first input, end is one beyond the last matching input
func newSlice(m *chart.Machine, reg TypeCreator, bff *BlockInfo, idx, end int) *chart.StateMix {
	id := bff.Id + " " + bff.Inputs[idx].Key
	return &chart.StateMix{
		Name: "Slice:" + id,
		// happens for every new flow in the series
		OnMap: func(_ string, flow jsn.FlowBlock) (okay bool) {
			if input, e := bff.ReadInput(idx); e != nil {
				log.Println(e)
			} else {
				m.PushState(newInnerBlock(m, reg, flow, input.BlockInfo))
				okay = true
			}
			return okay
		},
		// end of the inner block, advance to the next input
		OnCommit: func(interface{}) {
			idx++
		},
		// end of the repeating flows
		OnEnd: func() {
			m.FinishState(true)
		},
	}
}

func fillSlot(reg TypeCreator, slot jsn.SlotBlock, typeName string) (err error) {
	if i, ok := reg.NewType(typeName); !ok {
		err = errutil.New("couldn't create", typeName)
	} else if !slot.SetSlot(i) {
		err = errutil.New("couldn't set flow %T", i)
	}
	return
}

// an array of primitives is a list of fields .
func newListReader(m *chart.Machine, bff *BlockInfo, termName string, idx, end int) *chart.StateMix {
	return &chart.StateMix{
		Name: "List:" + bff.Id + ":" + termName,
		OnValue: func(n string, pv interface{}) (err error) {
			if idx < 0 {
				err = errutil.New("list underflow")
			} else if idx >= end {
				err = errutil.New("list overflow")
			} else {
				field := bff.Fields[idx]
				idx++ // next time, next field
				if e := storeValue(pv, field.Msg); e != nil {
					err = e
				}
			}
			return
		},
		OnEnd: func() {
			m.FinishState(true)
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
