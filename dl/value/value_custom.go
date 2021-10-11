package value

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
)

func Text_Marshal(n jsn.Marshaler, val *Text) {
	if ex, ok := n.(cout.TextMarshaler); !ok || !ex.TextValue(Text_Type, &val.Str) {
		Text_Marshal_Customized(n, val)
	}
}
