package value

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"git.sr.ht/~ionous/iffy/export/jsn/detailed"
)

func Text_Marshal(n jsn.Marshaler, val *Text) {
	if _, ok := n.(*detailed.DetailedMarshaler); ok {
		Text_Marshal_Customized(n, val)
	} else {
		// custom serialization to avoid conflicts with @variables
		str := val.Str
		if len(str) > 0 && str[0] == '@' {
			str = "@" + str
		}
		n.WriteValue(Text_Type, str)
	}
}
