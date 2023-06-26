package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
)

// Params 'Affinity' designates the storage type of a given parameter
// while 'class' is used to indicate an interpretation of that parameter, for example a reference to a kind.
// Pattern locals can have an initial value, other uses of parameter cannot.
type Params struct {
	Affinity  affine.Affinity   `if:"label=have"`
	Name      string            `if:"label=called,type=text"`
	Class     string            `if:"label=of,optional,type=text"`
	Initially assign.Assignment `if:"label=initially,optional"`
}

// ensure fields which reference aspects use the necessary formatting
func AspectParam(aspectName string) Params {
	return Params{Name: aspectName, Affinity: affine.Text, Class: aspectName}
}

func assertFields(kind string, ps []Params, w func(kind, name, class string, aff affine.Affinity, init assign.Assignment) error) (err error) {
	for _, p := range ps {
		if e := w(kind, p.Name, p.Class, p.Affinity, p.Initially); e != nil {
			err = e
			break
		}
	}
	return
}
