package core

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
)

func GetVar_Marshal(n jsn.Marshaler, val *GetVar) bool {
	ex, ok := n.(cout.VariableMarshaler)
	return (ok && ex.VariableValue(GetVar_Type, &val.Name.Str)) ||
		GetVar_Marshal_Customized(n, val)
}
