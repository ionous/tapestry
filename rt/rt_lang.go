// Code generated by "makeops"; edit at your own risk.
package rt

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
)

const Assignment_Type = "assignment"

var Assignment_Optional_Marshal = Assignment_Marshal

func Assignment_Marshal(n jsn.Marshaler, ptr *Assignment) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func Assignment_Repeats_Marshal(n jsn.Marshaler, vals *[]Assignment) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				Assignment_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const BoolEval_Type = "bool_eval"

var BoolEval_Optional_Marshal = BoolEval_Marshal

func BoolEval_Marshal(n jsn.Marshaler, ptr *BoolEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func BoolEval_Repeats_Marshal(n jsn.Marshaler, vals *[]BoolEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				BoolEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const Execute_Type = "execute"

var Execute_Optional_Marshal = Execute_Marshal

func Execute_Marshal(n jsn.Marshaler, ptr *Execute) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func Execute_Repeats_Marshal(n jsn.Marshaler, vals *[]Execute) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				Execute_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const NumListEval_Type = "num_list_eval"

var NumListEval_Optional_Marshal = NumListEval_Marshal

func NumListEval_Marshal(n jsn.Marshaler, ptr *NumListEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func NumListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]NumListEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				NumListEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const NumberEval_Type = "number_eval"

var NumberEval_Optional_Marshal = NumberEval_Marshal

func NumberEval_Marshal(n jsn.Marshaler, ptr *NumberEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func NumberEval_Repeats_Marshal(n jsn.Marshaler, vals *[]NumberEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				NumberEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const RecordEval_Type = "record_eval"

var RecordEval_Optional_Marshal = RecordEval_Marshal

func RecordEval_Marshal(n jsn.Marshaler, ptr *RecordEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func RecordEval_Repeats_Marshal(n jsn.Marshaler, vals *[]RecordEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				RecordEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const RecordListEval_Type = "record_list_eval"

var RecordListEval_Optional_Marshal = RecordListEval_Marshal

func RecordListEval_Marshal(n jsn.Marshaler, ptr *RecordListEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func RecordListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]RecordListEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				RecordListEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const TextEval_Type = "text_eval"

var TextEval_Optional_Marshal = TextEval_Marshal

func TextEval_Marshal(n jsn.Marshaler, ptr *TextEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func TextEval_Repeats_Marshal(n jsn.Marshaler, vals *[]TextEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				TextEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

const TextListEval_Type = "text_list_eval"

var TextListEval_Optional_Marshal = TextListEval_Marshal

func TextListEval_Marshal(n jsn.Marshaler, ptr *TextListEval) {
	if slat := *ptr; slat != nil {
		slat.(jsn.Marshalee).Marshal(n)
	}
	return
}

func TextListEval_Repeats_Marshal(n jsn.Marshaler, vals *[]TextListEval) {
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		if n.RepeatValues(cnt) {
			for _, el := range *vals {
				TextListEval_Marshal(n, &el)
			}
			n.EndValues()
		}
	}
	return
}

var Slots = []interface{}{
	(*Assignment)(nil),
	(*BoolEval)(nil),
	(*Execute)(nil),
	(*NumListEval)(nil),
	(*NumberEval)(nil),
	(*RecordEval)(nil),
	(*RecordListEval)(nil),
	(*TextEval)(nil),
	(*TextListEval)(nil),
}
