package core

import (
	"git.sr.ht/~ionous/iffy/export/jsn"
	"git.sr.ht/~ionous/iffy/export/jsn/detailed"
	"github.com/ionous/errutil"
)

func GetVar_Marshal(n jsn.Marshaler, val *GetVar) {
	if _, ok := n.(*detailed.Chart); ok {
		GetVar_Marshal_Customized(n, val)
	} else if str := val.Name.Str; len(str) > 0 && str[0] == '@' {
		// this would conflict with @@ text serialization.
		n.Warning(errutil.New("serialization doesn't support variables names starting with @"))
	} else {
		n.SpecifyValue(GetVar_Type, "@"+str)
	}
}
