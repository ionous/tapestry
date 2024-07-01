package mdl

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

type FieldInfo struct {
	Name, Class string
	Affinity    affine.Affinity
	Init        rt.Assignment
	Error       error
}

func (f *FieldInfo) isAspectLike() (ret bool) {
	return f.Affinity == affine.Text && f.Name == f.Class
}

// shortcut: if we specify a record affinity, but no class
// the returned class is the name of the field.
func (f *FieldInfo) getDefaultClass() (ret string) {
	if cls := f.Class; len(cls) > 0 {
		ret = cls
	} else if isRecordAffinity(f.Affinity) {
		ret = f.Name
	}
	return
}

// does this field have a (non-zero) name and affinity?
func (f *FieldInfo) validate() (err error) {
	if f.Error != nil {
		err = f.Error
	} else if len(f.Name) == 0 {
		err = errors.New("missing name")
	} else if len(f.Affinity) == 0 {
		err = errors.New("missing affinity")
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
