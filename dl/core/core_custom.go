package core

import (
	"git.sr.ht/~ionous/iffy/jsn"
)

func GetVar_Marshal(n jsn.Marshaler, val *GetVar) bool {
	// ex, ok := n.(cout.VariableMarshaler)
	// return (ok && ex.VariableValue(GetVar_Type, &val.Name.Str)) ||
	return GetVar_Marshal_Customized(n, val)
}
