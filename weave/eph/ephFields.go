package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"github.com/ionous/errutil"
)

// Params 'Affinity' designates the storage type of a given parameter
// while 'class' is used to indicate an interpretation of that parameter, for example a reference to a kind.
// Pattern locals can have an initial value, other uses of parameter cannot.
type Params struct {
	Affinity  Affinity          `if:"label=have"`
	Name      string            `if:"label=called,type=text"`
	Class     string            `if:"label=of,optional,type=text"`
	Initially assign.Assignment `if:"label=initially,optional"`
}

func weaveFields(kind string, ps []Params, w func(kind, name, class string, aff affine.Affinity, init assign.Assignment) error) (err error) {
	for _, p := range ps {
		if e := weaveField(kind, p, w); e != nil {
			err = e
			break
		}
	}
	return
}

func weaveField(kind string, p Params, w func(kind, name, class string, aff affine.Affinity, init assign.Assignment) error) (err error) {
	if aff, e := fromAffinity(p.Affinity); e != nil {
		err = e
	} else if e := w(kind, p.Name, p.Class, aff, p.Initially); e != nil {
		err = e
	}
	return
}

func fromAffinity(fieldAffinity Affinity) (ret affine.Affinity, err error) {
	if aff, ok := composer.FindChoice(&fieldAffinity, fieldAffinity.Str); !ok && len(fieldAffinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else {
		ret = affine.Affinity(aff)
	}
	return
}
