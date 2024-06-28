package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// expects field to be a normalized name already.
func (run *Runner) setObjectField(obj query.NounInfo, field string, newValue rt.Value) (err error) {
	if kind, e := run.getKind(obj.Kind); e != nil {
		err = e
	} else if fieldIndex := kind.FieldIndex(field); fieldIndex < 0 {
		err = rt.UnknownField(obj.String(), field)
	} else {
		fieldData := kind.Field(fieldIndex)
		if fieldData.Name == field {
			// the field we found is the name we asked for: its a normal field.
			err = run.writeNounValue(obj, fieldData, newValue)
		} else {
			// when the name differs, we must have found the aspect for a trait.
			if aff := newValue.Affinity(); aff != affine.Bool {
				err = errutil.New("can only set a trait with booleans, have", aff)
			} else if trait, e := oppositeDay(run, fieldData.Name, field, newValue.Bool()); e != nil {
				err = e
			} else {
				// set the aspect to the value of the requested trait
				traitValue := rt.StringFrom(trait, fieldData.Type)
				if was, e := run.readNounValue(obj, fieldData); e != nil {
					err = e
				} else if was := was.String(); was != trait {
					if e := run.writeNounValue(obj, fieldData, traitValue); e != nil {
						err = e
					} else if notify := run.notify.ChangedState; notify != nil {
						notify(obj.Noun, fieldData.Name, was, field)
					}
				}
			}
		}
	}
	return
}

// expects field to be a normalized name already.
func (run *Runner) getObjectField(obj query.NounInfo, field string) (ret rt.Value, err error) {
	if kind, e := run.getKind(obj.Kind); e != nil {
		err = e
	} else if fieldIndex := kind.FieldIndex(field); fieldIndex < 0 {
		err = rt.UnknownField(obj.String(), field)
	} else {
		fieldData := kind.Field(fieldIndex)
		if v, e := run.readNounValue(obj, fieldData); e != nil {
			err = e
		} else if fieldData.Name == field {
			// the field we found is the name we asked for:
			// so this is a regular field.
			ret = v
		} else {
			// when the name differs, we must have found the aspect for a trait.
			// note: kind also does this -- but since the data here isnt stored in a record
			// ( it's stored in the noun value cache ) we have to duplicate the aspect/field check.
			// return true if the aspect field holds the particular requested field
			ret = rt.BoolOf(field == v.String())
		}
	}
	return
}

func (run *Runner) writeNounValue(obj query.NounInfo, field rt.Field, val rt.Value) (err error) {
	// fix: convert when appropriate?
	if aff := val.Affinity(); aff != field.Affinity {
		err = errutil.Fmt(`mismatched affinity "%s.%s(%s)" writing %s`, obj, field.Name, field.Affinity, aff)
	} else {
		key := query.MakeKey(obj.Domain, obj.Noun, field.Name)
		userVal := UserValue{rt.CopyValue(val)}
		run.dynamicVals.Store(key, userVal)
	}
	return
}

// return the (cached) value of a noun's field
// if the noun's field contains an assignment it's evaluated each time.
func (run *Runner) readNounValue(obj query.NounInfo, ft rt.Field) (ret rt.Value, err error) {
	key := query.MakeKey(obj.Domain, obj.Noun, ft.Name)
	// first ensure the value is in the cache
	if _, e := run.dynamicVals.Ensure(key, func() (ret any, err error) {
		ret, err = run.query.NounValue(obj.Noun, ft.Name)
		return
	}); e != nil {
		err = e
	} else {
		// then ask for the value again to unpack it.
		ret, err = run.unpackDynamicValue(key, ft.Affinity, ft.Type)
	}
	return
}
