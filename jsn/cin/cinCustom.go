package cin

import (
	"encoding/json"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"git.sr.ht/~ionous/iffy/rt"
)

var custom = chart.Customization{
	// package rt:
	rt.Assignment_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.Assignment)
		if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.Assignment_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.BoolEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.BoolEval)
		if v, ok := readBool(dec); ok {
			*ptr = v
		} else if v, ok := readVar(dec); ok {
			dec.Commit("bool var")
			*ptr = v
		} else {
			err = rt.BoolEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.NumberEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.NumberEval)
		if v, ok := readNum(dec); ok {
			*ptr = v
		} else if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.NumberEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.NumListEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.NumListEval)
		if v, ok := readNumList(dec); ok {
			*ptr = v
		} else if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.NumListEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.RecordEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.RecordEval)
		if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.RecordEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.RecordListEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.RecordListEval)
		if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.RecordListEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.TextEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.TextEval)
		// NOTE: a string where a text eval is expected could be a variable providing text or some literal text.
		// variables start with @, and text that starts with an @ winds up prefixed with two @@s
		if v, ok := readVarOrText(dec); ok {
			*ptr = v
		} else {
			err = rt.TextEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	rt.TextListEval_Type: func(n jsn.Marshaler, i interface{}) (err error) {
		dec, ptr := n.(*xDecoder), i.(*rt.TextListEval)
		if v, ok := readTextList(dec); ok {
			*ptr = v
		} else if v, ok := readVar(dec); ok {
			*ptr = v
		} else {
			err = rt.TextListEval_DefaultMarshal(dec, ptr)
		}
		return
	},
	// package core:
	// if someone authored compact data with an explicit text value or get variable
	// ex. { "TextValue:": "some text" } or { "GetVar:": "var name" }
	// we can allow the default handlers to handle reading them.
	// ( we dont expect those to be escaped )
}

func readBool(dec *xDecoder) (ret *core.BoolValue, okay bool) {
	var val bool
	if e := json.Unmarshal(dec.CurrentMessage, &val); e == nil {
		ret, okay = &core.BoolValue{val}, true
		dec.Commit("bool literal")
	}
	return
}
func readNum(dec *xDecoder) (ret *core.NumValue, okay bool) {
	var val float64
	if e := json.Unmarshal(dec.CurrentMessage, &val); e == nil {
		ret, okay = &core.NumValue{val}, true
		dec.Commit("num literal")
	}
	return
}
func readNumList(dec *xDecoder) (ret *core.Numbers, okay bool) {
	var val []float64
	if e := json.Unmarshal(dec.CurrentMessage, &val); e == nil {
		ret, okay = &core.Numbers{val}, true
		dec.Commit("num list literal")
	}
	return
}
func readTextList(dec *xDecoder) (ret *core.Texts, okay bool) {
	var val []string
	if e := json.Unmarshal(dec.CurrentMessage, &val); e == nil {
		ret, okay = &core.Texts{val}, true
		dec.Commit("text list literal")
	}
	return
}

// when there's just a single string that fits a text eval..
// it could be literal text or a variable providing text.
func readVarOrText(dec *xDecoder) (ret rt.TextEval, okay bool) {
	var str string
	if e := json.Unmarshal(dec.CurrentMessage, &str); e == nil {
		if cnt := len(str); cnt == 0 || str[0] != '@' {
			ret, okay = &core.TextValue{str}, true
			dec.Commit("simple text literal")
		} else {
			if cnt > 2 && str[1] == '@' {
				// any text primitive with an @ gets another @ prefixed to it
				// so.. strip that off here.
				ret, okay = &core.TextValue{str[1:]}, true
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

func readVar(dec *xDecoder) (ret *core.GetVar, okay bool) {
	var str string
	if e := json.Unmarshal(dec.CurrentMessage, &str); e == nil {
		if len(str) > 0 && str[0] == '@' {
			ret, okay = newVar(str), true
			dec.Commit("@var")
		}
	}
	return
}

func newVar(str string) *core.GetVar {
	return &core.GetVar{Name: value.VariableName{Str: str[1:]}}
}
