package value

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/detailed"
)

func Text_Marshal(n jsn.Marshaler, val *Text) {
	if _, ok := n.(detailed.Chart); ok {
		Text_Marshal_Customized(n, val)
	} else {
		n.StrValue(val)
	}
}

func (v *Text) GetStr() (ret string) {
	// custom serialization to avoid conflicts with @variables
	if str := v.Str; len(str) > 0 && str[0] == '@' {
		ret = "@" + str
	} else {
		ret = str
	}
	return
}

func (v *Text) SetStr(string) {
	panic("strip text")
}
