package scope

import (
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// a scope which provides a single named field.
type Pairs struct {
	ks []string
	vs []g.Value
}

//
func NewPairs(ks []string, vs []g.Value) *Pairs {
	return &Pairs{ks: ks, vs: vs}
}

func (w *Pairs) find(field string) (ret int, err error) {
	var found bool
	for i, k := range w.ks {
		if k == field {
			ret, found = i, true
			break
		}
	}
	if !found {
		err = g.UnknownVariable(field)
	}
	return
}

func (w *Pairs) FieldByName(field string) (ret g.Value, err error) {
	if i, e := w.find(field); e != nil {
		err = e
	} else {
		ret = w.vs[i]
	}
	return
}

func (w *Pairs) SetFieldByName(field string, val g.Value) (err error) {
	if i, e := w.find(field); e != nil {
		err = e
	} else {
		w.vs[i] = val
	}
	return
}

// placeholder method: for now, determines whether the field exists
func (w *Pairs) SetFieldDirty(field string) (err error) {
	_, err = w.find(field)
	return
}
