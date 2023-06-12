package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
)

type UniformField struct {
	Name, Type string
	Affinity   affine.Affinity
	At         string
}

// normalize the values of the field
func MakeUniformField(aff affine.Affinity, fieldName, fieldClass, at string) (ret UniformField, err error) {
	if name, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		// shortcut: if we specify a field name for a record and no class, we'll expect the class to be the name.
		if len(class) == 0 && isRecordAffinity(aff) {
			class = name
		}
		ret = UniformField{Name: name, Affinity: aff, Type: class, At: at}
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isRecordAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case affine.Record, affine.RecordList:
		okay = true
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isClassAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case "", affine.Text, affine.TextList, affine.Record, affine.RecordList:
		okay = true
	}
	return
}
