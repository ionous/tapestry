package cout

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"github.com/ionous/errutil"
)

// xEncoder - marker so callers can see where a machine pointer came from.
type xEncoder struct {
	chart.Machine
	customFlow CustomFlow
}

type CustomFlow func(jsn.Marshaler, jsn.FlowBlock) error

// NewEncoder create an empty serializer to produce compact script data.
func Encode(in jsn.Marshalee, customFlow CustomFlow) (ret interface{}, err error) {
	if customFlow == nil {
		customFlow = func(jsn.Marshaler, jsn.FlowBlock) error {
			return chart.Unhandled("no custom encoder")
		}
	}
	m := xEncoder{Machine: chart.MakeEncoder(), customFlow: customFlow}
	next := m.newValue(m.newBlock())
	next.OnCommit = func(v interface{}) {
		if ret != nil {
			m.Error(errutil.New("can only write data once"))
		} else {
			ret = v
		}
	}
	m.ChangeState(next)
	in.Marshal(&m)
	return ret, m.Errors()
}

func unpack(pv interface{}) (ret interface{}) {
	switch pv := pv.(type) {
	case interface{ GetCompactValue() interface{} }:
		ret = pv.GetCompactValue()
	case interface{ GetValue() interface{} }:
		ret = pv.GetValue()
	default:
		ret = pv // provisionally
		if pstr, isString := pv.(*string); isString {
			strs := strings.Split(*pstr, "\n")
			if len(strs) > 1 {
				ret = strs
			}
		}
	}
	return
}

// compact data represents primitive values as their value.
func (m *xEncoder) newValue(next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(_ string, pv interface{}) error {
		m.Commit(unpack(pv))
		return nil
	}
	return next
}

// blocks handle beginning new flows, swaps, or repeats
// end ( and how they collect data ) gets left to the caller
func (m *xEncoder) newBlock() *chart.StateMix {
	var next chart.StateMix
	return m.addBlock(&next)
}

func (m *xEncoder) addBlock(next *chart.StateMix) *chart.StateMix {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	next.OnMap = func(typeName string, block jsn.FlowBlock) (okay bool) {
		// maybe there was no key, so this is the first comment call we've gotten for the block.
		var c string
		if m.Machine.Comment != nil {
			c = *m.Machine.Comment
			m.Machine.Comment = nil
		}
		if e := m.customFlow(m, block); e != nil {
			var unhandled chart.Unhandled
			if !errors.As(e, &unhandled) {
				m.Error(e)
			} else {
				m.PushState(m.newFlow(newComFlow(block.GetLede(), c)))
				okay = true // return true to indicate caller should descend into the flow
			}
		}
		return
	}
	next.OnSlot = func(_ string, slot jsn.SlotBlock) (okay bool) {
		if _, ok := slot.GetSlot(); ok {
			m.PushState(m.newSlot())
			okay = true
		}
		return
	}
	// ex."noun_phrase" "$KIND_OF_NOUN"
	next.OnSwap = func(typeName string, p jsn.SwapBlock) (okay bool) {
		if choice, _ := p.GetSwap(); len(choice) > 0 {
			m.PushState(m.newSwap(&comSwap{typeName: typeName, choice: choice}))
			okay = true
		}
		return okay
	}
	next.OnRepeat = func(t string, vs jsn.SliceBlock) bool {
		var slice []interface{}
		if cnt := vs.GetSize(); cnt >= 0 {
			slice = make([]interface{}, 0, cnt)
		}
		m.PushState(m.newSlice(slice))
		return true
	}
	// next.OnEnd... gets determined by the specific block
	return next
}

func (m *xEncoder) newFlow(d *comFlow) *chart.StateMix {
	var next chart.StateMix
	next.OnKey = func(key, _ string) error {
		m.ChangeState(m.newKeyValue(next, d, key))
		return nil
	}
	next.OnEnd = func() {
		// doesnt worry if there's a pending key/value
		// writing a value to a key is always considered optional
		m.FinishState(d.finalize())
	}
	return &next
}

// writes the value into the key and change back to the flow state
func (m *xEncoder) newKeyValue(prev chart.StateMix, d *comFlow, key string) *chart.StateMix {
	next := m.newValue(m.addBlock(&prev))
	next.OnCommit = func(v interface{}) {
		if c, ok := v.(*comSwap); ok {
			d.addMsgPair(key, c.choice, c.value)
		} else {
			d.addMsg(key, v)
			m.ChangeState(&prev)
		}
	}
	return next
}

// every slice pushes a brand new machine
func (m *xEncoder) newSlice(vals []interface{}) *chart.StateMix {
	next := m.newValue(m.newBlock())
	// a new value is being added to our slice
	next.OnCommit = func(v interface{}) {
		vals = append(vals, v)
	}
	// the slice is done, write it to our parent whomever that is.
	next.OnEnd = func() {
		m.FinishState(vals)
	}
	return next
}

// every slot pushes a brand new machine
func (m *xEncoder) newSlot() *chart.StateMix {
	next := m.newValue(m.newBlock())
	next.OnCommit = func(v interface{}) {
		// write our choice and change into an error checking state
		m.ChangeState(chart.NewBlockResult(&m.Machine, v))
	}
	// fix? what should an uncommitted choice write?
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func (m *xEncoder) newSwap(swap *comSwap) *chart.StateMix {
	next := m.newValue(m.newBlock())
	next.OnCommit = func(v interface{}) {
		m.ChangeState(chart.NewBlockResult(&m.Machine, swap.SetValue(v)))
	}
	return next
}
