package pattern

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

type FieldInit struct {
	Field string
	Init  rt.Assignment
}

// pop from the front of the slice
func shift[S ~[]E, E any](s S) (E, S) {
	return s[0], s[1:]
}

// InitRecord fill the passed record with the arguments named by keys and values
// keys, if any, must match the order of fields in the kind ( but don't need to be contiguous. )
// if there are more values than keys, the first values are treated as indexed values;
// having fewer values than keys is invalid.
func InitRecord(run rt.Runtime, k *rt.Kind, keys []string, vals []rt.Value) (ret *rt.Record, err error) {
	rec := rt.NewRecord(k) // todo: add an "NewUninitializedRecord" to save a few cycles?
	scope := windowingScope{rec, make(map[string]bool)}
	run.PushScope(scope)
	vars := initHelper{keys, vals}
	for i, cnt := 0, k.NumField(); i < cnt; i++ {
		ft := k.Field(i)
		if v, e := vars.nextValue(run, ft); e != nil {
			err = e
			break
		} else if v != nil {
			if convertedVal, e := safe.RectifyText(run, v, ft.Affinity, ft.Type); e != nil {
				err = fmt.Errorf("%w while initializing arg %d(%s)", e, i, ft.Name)
				break
			} else if e := rec.SetIndexedField(i, convertedVal); e != nil {
				err = fmt.Errorf("%w while setting arg %d(%s)", e, i, ft.Name)
				break
			}
		}
		// add this field to the scope
		scope.UpdateWindow(i)
	}
	if err == nil {
		if e := vars.sanity(); e != nil {
			err = e
		} else {
			ret = rec
		}
	}
	run.PopScope()
	return
}

type initHelper struct {
	keys []string
	vals []rt.Value
}

// sanity checks
func (in *initHelper) sanity() (err error) {
	if cnt := len(in.keys); cnt > 0 {
		err = fmt.Errorf("initialization has %d unused keys", cnt)
	} else if cnt := len(in.vals); cnt > 0 {
		err = fmt.Errorf("initialization has %d unused values", cnt)
	}
	return
}

// can return nil
func (in *initHelper) nextValue(run rt.Runtime, ft rt.Field) (ret rt.Value, err error) {
	if len(in.keys) < len(in.vals) { // use indexed values
		ret, in.vals = shift(in.vals)
	} else if len(in.keys) > 0 && in.keys[0] == ft.Name {
		ret, in.vals = shift(in.vals) // then named fields
		_, in.keys = shift(in.keys)
	} else if a := ft.Init; a != nil { // falling back to initializers
		ret, err = a.GetAssignedValue(run)
	} else {
		// cheat to make default records when calling patterns
		if ft.Affinity == affine.Record {
			if subk, e := run.GetKindByName(ft.Type); e != nil {
				err = e
			} else {
				ret = rt.RecordOf(rt.NewRecord(subk))
			}
		}
	}
	return
}
