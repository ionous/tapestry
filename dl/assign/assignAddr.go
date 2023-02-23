package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// uniform access to objects and variables
// implemented by the members of Address
type Address interface {
	// returns the object name, meta.Variables, or an error
	GetObjectName(rt.Runtime) (string, error)
	// returns the targeted field in the named object, the targeted variable, or an error
	GetFieldName(rt.Runtime) (string, error)
	// path within the value identified by this reference
	GetPath() []Dot
}

// determine the location of a value
func GetReference(run rt.Runtime, addr Address) (ret RefValue, err error) {
	if objName, e := addr.GetObjectName(run); e != nil {
		err = e
	} else if fieldName, e := addr.GetFieldName(run); e != nil {
		err = e
	} else if path, e := ResolvePath(run, addr.GetPath()); e != nil {
		err = e
	} else {
		ret = RefValue{objName, fieldName, path}
	}
	return
}

// pull the object.field value from the runtime ( without expanding its path )
func GetRootValue(run rt.Runtime, addr Address) (ret RootValue, err error) {
	if ref, e := GetReference(run, addr); e != nil {
		err = e
	} else {
		ret, err = ref.GetRootValue(run)
	}
	return
}
