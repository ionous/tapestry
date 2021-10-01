package core

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
)

func GetVar_Marshal(n jsn.Marshaler, val *GetVar) {
	if ex, ok := n.(cout.VariableMarshaler); !ok || !ex.VariableValue(&val.Name.Str) {
		GetVar_Marshal_Customized(n, val)
	}
}
