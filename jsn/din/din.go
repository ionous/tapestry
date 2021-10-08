package din

import (
	"encoding/json"

	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

type Decoder struct {
	*chart.Machine
}

func NewDecoder(msg []byte) Decoder {
	return Decoder{chart.NewDecoder(func(m *chart.Machine) *chart.StateMix {
		return newBlock(m, msg)
	})}
}

func readEnum(m *chart.Machine, msg json.RawMessage, val chart.EnumMarshaler) {
	var str string
	if e := json.Unmarshal(msg, &str); e != nil {
		m.Warning(e)
	} else {
		val.SetEnum(str)
	}
	return
}

func newValue(m *chart.Machine, msg json.RawMessage, next *chart.StateMix) *chart.StateMix {
	next.OnValue = func(typeName string, pv interface{}) {
		if enum, ok := pv.(chart.EnumMarshaler); ok {
			readEnum(m, msg, enum)
		} else if e := json.Unmarshal(msg, pv); e != nil {
			m.Warning(e)
		}
	}
	return next
}

func newBlock(m *chart.Machine, msg json.RawMessage) *chart.StateMix {
	next := chart.NewReportingState(m)
	next.OnMap = func(l, k string) (okay bool) {
		var v dinMap
		if e := json.Unmarshal(msg, &v); e != nil {
			m.Warning(e)
		} else if v.Type != k {
			m.Warning(errutil.New("expected", k, "found", v.Type))
		} else {
			m.PushState(newFlow(m, v.Fields))
		}
		return
	}
	next.OnLiteral = func(field string) bool {
		return next.OnMap("", field)
	}
	// // ex."noun_phrase" "$KIND_OF_NOUN"
	// next.OnPick = func(t string, p jsn.Picker) (okay bool) {
	// 	if choice, ok := p.GetChoice(); !ok {
	// 		m.Error(errutil.New("couldnt determine choice of", p))
	// 	} else if len(choice) > 0 {
	// 		kind := p.GetType()
	// 		m.PushState(newSwap(m, choice, detMap{
	// 			Id:   m.FlushCursor(),
	// 			Type: kind,
	// 		}))
	// 		okay = true
	// 	}
	// 	return
	// }
	next.OnRepeat = func(t string, vs jsn.Slicer) (okay bool) {
		if hint := vs.GetSize(); hint > 0 {
			// m.PushState(newSlice(m, make([]interface{}, 0, hint)))
			okay = true
		}
		return
	}
	return next
}

func newFlow(m *chart.Machine, fields dinFields) *chart.StateMix {
	next := chart.NewReportingState(m)
	next.OnKey = func(_, key string) (okay bool) {
		if msg, ok := fields[key]; ok {
			m.ChangeState(newKey(m, *next, msg))
			okay = true
		}
		return okay
	}
	next.OnEnd = func() {
		m.FinishState(nil)
	}
	return next
}

func newKey(m *chart.Machine, prev chart.StateMix, msg json.RawMessage) *chart.StateMix {
	next := newValue(m, msg, &prev)
	next.OnCommit = func(interface{}) {
		m.ChangeState(&prev)
	}
	return next
}
