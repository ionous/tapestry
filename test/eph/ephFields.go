package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// Params 'Affinity' designates the storage type of a given parameter
// while 'class' is used to indicate an interpretation of that parameter, for example a reference to a kind.
// Pattern locals can have an initial value, other uses of parameter cannot.
type Params struct {
	Affinity  affine.Affinity `if:"label=have"`
	Name      string          `if:"label=called,type=text"`
	Class     string          `if:"label=of,optional,type=text"`
	Initially rt.Assignment   `if:"label=initially,optional"`
}

func (p Params) FieldInfo() mdl.FieldInfo {
	return mdl.FieldInfo{
		Name:     en.Normalize(p.Name),
		Class:    en.Normalize(p.Class),
		Affinity: p.Affinity,
		Init:     p.Initially,
	}
}

// ensure fields which reference aspects use the necessary formatting
func AspectParam(aspectName string) Params {
	return Params{Name: aspectName, Affinity: affine.Text, Class: aspectName}
}
