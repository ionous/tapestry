package core

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/detailed"
	"github.com/ionous/errutil"
)

func GetVar_Marshal(n jsn.Marshaler, val *GetVar) {
	if _, ok := n.(*detailed.Chart); ok {
		GetVar_Marshal_Customized(n, val)
	} else if str := val.Name.Str; len(str) > 0 && str[0] == '@' {
		// this would conflict with @@ text serialization.
		n.Warning(errutil.New("serialization doesn't support variables names starting with @"))
	} else {
		n.StrValue(val)
	}
}

func (v *GetVar) GetStr() (ret string) {
	return "@" + v.Name.Str
}

func (v *GetVar) SetStr(string) {
	panic("strip text")
}
