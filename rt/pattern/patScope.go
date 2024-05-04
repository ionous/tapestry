package pattern

import "git.sr.ht/~ionous/tapestry/rt"

// as fields of a pattern record are initialized,
// we expose those to the scope for the remaining fields to see.
type windowingScope struct {
	rec    *rt.Record
	window map[string]bool
}

func (w *windowingScope) UpdateWindow(i int) {
	ft := w.rec.Field(i)
	w.window[ft.Name] = true
}

func (w windowingScope) FieldChanged(field string) bool {
	return w.window[field]
}

func (w windowingScope) FieldByName(field string) (ret rt.Value, err error) {
	if !w.window[field] {
		err = rt.UnknownVariable(field)
	} else {
		ret, err = w.rec.GetNamedField(field)
	}
	return
}

func (w windowingScope) SetFieldByName(field string, val rt.Value) (err error) {
	if !w.window[field] {
		err = rt.UnknownVariable(field)
	} else {
		err = w.rec.SetNamedField(field, val)
	}
	return
}
