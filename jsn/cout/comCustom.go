package cout

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

var custom = chart.Customization{
	core.BoolValue_Type: func(n jsn.Marshaler, i interface{}) bool {
		var out bool = i.(*core.BoolValue).Bool
		return n.MarshalValue(core.BoolValue_Type, out)
	},
	core.NumValue_Type: func(n jsn.Marshaler, i interface{}) bool {
		var out float64 = i.(*core.NumValue).Num
		return n.MarshalValue(core.NumValue_Type, out)
	},
	core.Numbers_Type: func(n jsn.Marshaler, i interface{}) bool {
		var out []float64 = i.(*core.Numbers).Values
		return n.MarshalValue(core.Numbers_Type, out)
	},

	// write text as a raw string
	core.TextValue_Type: func(n jsn.Marshaler, i interface{}) bool {
		str := i.(*core.TextValue).Text
		// if the text starts with an @, add another @
		if len(str) > 0 && str[0] == '@' {
			str = "@" + str
		}
		return n.MarshalValue(core.TextValue_Type, str)
	},
	//
	core.Texts_Type: func(n jsn.Marshaler, i interface{}) bool {
		var out []string = i.(*core.Texts).Values
		return n.MarshalValue(core.Texts_Type, out)
	},

	// write variables as a string prepended by @
	core.GetVar_Type: func(n jsn.Marshaler, i interface{}) (okay bool) {
		ptr := i.(*core.GetVar)
		str := ptr.Name.Str
		// a leading ampersand would conflict with @@ escaped text serialization.
		if leadingAmp := len(str) > 0 && str[0] == '@'; !leadingAmp {
			okay = n.MarshalValue(core.GetVar_Type, "@"+str)
		} else {
			okay = core.GetVar_DefaultMarshal(n, ptr)
		}
		return
	},
}
