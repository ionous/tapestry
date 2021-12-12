package cin

import (
	"encoding/json"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

func (dec *xDecoder) customFlow(flow jsn.FlowBlock, msg json.RawMessage) (err error) {
	switch typeName := flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case story.NamedNoun_Type:
		var str string
		if e := json.Unmarshal(msg, &str); e != nil {
			err = chart.Unhandled(typeName)
		} else {
			var out story.NamedNoun
			if space := strings.Index(str, " "); space < 0 {
				out.Name.Str = str
			} else {
				jsn.MakeEnum(&out.Determiner, &out.Determiner.Str).SetValue(str[:space])
				out.Name.Str = str[space+1:]
			}
			if !flow.SetFlow(&out) {
				err = errutil.New("could set result to flow", typeName, flow)
			}
		}
	}
	return
}
func (dec *xDecoder) customSlot(slot jsn.SlotBlock, msg json.RawMessage) (err error) {
	switch typeName := slot.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomSlot")

	case rt.Assignment_Type:
		if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.BoolEval_Type:
		if v, ok := readBool(dec, msg); ok {
			slot.SetSlot(v)
		} else if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.NumberEval_Type:
		if v, ok := readNum(dec, msg); ok {
			slot.SetSlot(v)
		} else if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.NumListEval_Type:
		if v, ok := readNumList(dec, msg); ok {
			slot.SetSlot(v)
		} else if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.RecordEval_Type:
		if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.RecordListEval_Type:
		if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.TextEval_Type:
		// NOTE: a string where a text eval is expected could be a variable providing text or some literal text.
		// variables start with @, and text that starts with an @ winds up prefixed with two @@s
		if v, ok := readVarOrText(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

	case rt.TextListEval_Type:
		if v, ok := readTextList(dec, msg); ok {
			slot.SetSlot(v)
		} else if v, ok := readVar(dec, msg); ok {
			slot.SetSlot(v)
		} else {
			err = chart.Unhandled(typeName)
		}

		// package core:
		// if someone authored compact data with an explicit text value or get variable
		// ex. { "TextValue:": "some text" } or { "GetVar:": "var name" }
		// we can allow the default handlers to handle reading them.
		// ( we dont expect those to be escaped )
	}
	return
}

func readBool(dec *xDecoder, msg json.RawMessage) (ret *literal.BoolValue, okay bool) {
	var val bool
	if e := json.Unmarshal(msg, &val); e == nil {
		ret, okay = &literal.BoolValue{val}, true
		dec.Commit("bool literal")
	}
	return
}
func readNum(dec *xDecoder, msg json.RawMessage) (ret *literal.NumValue, okay bool) {
	var val float64
	if e := json.Unmarshal(msg, &val); e == nil {
		ret, okay = &literal.NumValue{val}, true
		dec.Commit("num literal")
	}
	return
}
func readNumList(dec *xDecoder, msg json.RawMessage) (ret *literal.NumValues, okay bool) {
	var val []float64
	if e := json.Unmarshal(msg, &val); e == nil {
		ret, okay = &literal.NumValues{val}, true
		dec.Commit("num list literal")
	}
	return
}
func readTextList(dec *xDecoder, msg json.RawMessage) (ret *literal.TextValues, okay bool) {
	var val []string
	if e := json.Unmarshal(msg, &val); e == nil {
		ret, okay = &literal.TextValues{val}, true
		dec.Commit("text list literal")
	}
	return
}

// when there's just a single string that fits a text eval..
// it could be literal text or a variable providing text.
func readVarOrText(dec *xDecoder, msg json.RawMessage) (ret rt.TextEval, okay bool) {
	var str string
	if e := json.Unmarshal(msg, &str); e == nil {
		if cnt := len(str); cnt == 0 || str[0] != '@' {
			ret, okay = &literal.TextValue{str}, true
			dec.Commit("simple text literal")
		} else {
			if cnt > 2 && str[1] == '@' {
				// any text primitive with an @ gets another @ prefixed to it
				// so.. strip that off here.
				ret, okay = &literal.TextValue{str[1:]}, true
				dec.Commit("escaped text literal")
			} else {
				// only has one @.. that means its a var
				ret, okay = newVar(str), true
				dec.Commit("text var")
			}
		}
	}
	return
}

func readVar(dec *xDecoder, msg json.RawMessage) (ret *core.GetVar, okay bool) {
	var str string
	if e := json.Unmarshal(msg, &str); e == nil {
		if len(str) > 0 && str[0] == '@' {
			ret, okay = newVar(str), true
			dec.Commit("@var")
		}
	}
	return
}

func newVar(str string) *core.GetVar {
	return &core.GetVar{Name: core.VariableName{Str: str[1:]}}
}
