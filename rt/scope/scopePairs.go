package scope

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// a scope which provides a single named field.
type Pairs struct {
	ks      []string
	vs      []rt.Value
	changed map[string]bool
}

func NewPairs(ks []string, vs []rt.Value) *Pairs {
	return &Pairs{ks: ks, vs: vs, changed: make(map[string]bool)}
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
		err = rt.UnknownVariable(field)
	}
	return
}

func (w *Pairs) FieldChanged(field string) bool {
	return w.changed[field]
}

func (w *Pairs) FieldByName(field string) (ret rt.Value, err error) {
	if i, e := w.find(field); e != nil {
		err = e
	} else {
		ret = w.vs[i]
	}
	return
}

func (w *Pairs) SetFieldByName(field string, val rt.Value) (err error) {
	if i, e := w.find(field); e != nil {
		err = e
	} else {
		w.vs[i] = val
		w.changed[field] = true
	}
	return
}
