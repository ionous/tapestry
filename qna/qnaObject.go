package qna

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"github.com/ionous/errutil"
)

// expects field to be a normalized name already.
func (run *Runner) setObjectField(obj query.NounInfo, field string, newValue rt.Value) (err error) {
	// tbd: cache the kind in the object info?
	// or even... cache the ( last n ) field info into "obj.field"?
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
						notify(obj.Id, fieldData.Name, was, field)
					}
				}
			}
		}
	}
	return
}

// expects field to be a normalized name already.
func (run *Runner) getObjectField(obj query.NounInfo, field string) (ret rt.Value, err error) {
	// tbd: cache the kind in the object info?
	// or even... cache the ( last n ) field info into "obj.field"?
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
	// fix: convert when appropriate.
	if aff := val.Affinity(); aff != field.Affinity {
		err = errutil.Fmt(`mismatched affinity "%s.%s(%s)" writing %s`, obj, field.Name, field.Affinity, aff)
	} else {
		key := makeKey(obj.Domain, obj.Id, field.Name)
		userVal := UserValue{rt.CopyValue(val)}
		run.dynamicVals.store[key] = userVal
	}
	return
}

// return the (cached) value of a noun's field
// if the noun's field contains an assignment it's evaluated each time.
func (run *Runner) readNounValue(obj query.NounInfo, ft rt.Field) (ret rt.Value, err error) {
	key := makeKey(obj.Domain, obj.Id, ft.Name)

	// kind of ugly: first ensure its in the cache
	// by generating a Value or Assignment
	if _, e := run.dynamicVals.ensure(key, func() (ret any, err error) {
		// a record can have multiple path/values
		if vs, e := run.query.NounValues(obj.Id, ft.Name); e != nil {
			err = e
		} else if len(vs) > 0 {
			ret, err = run.readFields(ft, vs)
		} else {
			// if the noun had no values; the kind might have default values.
			ret, err = run.readKindField(obj, ft)
		}
		return
	}); e != nil {
		err = e
	} else {
		// then ask for the value again to unpack it.
		ret, err = run.unpackDynamicValue(key, ft.Affinity, ft.Type)
	}

	return
}

// returns an rt.Assignment, rt.Value, or error
func (run *Runner) readKindField(obj query.NounInfo, field rt.Field) (ret any, err error) {
	if k, e := run.getKind(obj.Kind); e != nil {
		err = e
	} else if fieldIndex := k.FieldIndex(field.Name); fieldIndex < 0 {
		err = errutil.New("couldnt find field %q in kind %q", field.Name, k.Name)
	} else if init := field.Init; init != nil {
		ret = init
	} else {
		// note: this doesnt properly determine the default trait for an aspect
		// weave works around this by providing the correct default value in the db
		ret, err = rt.ZeroField(field.Affinity, field.Type, fieldIndex)
	}
	return
}

// return can be an assignment ( which gets evaluated )
// or a literal value ( a fixed value )
func (run *Runner) readFields(field rt.Field, vals []query.ValueData) (ret any, err error) {
	if !dotted(vals) { // a single top level value? then its an assignment
		ret, err = run.decode.DecodeAssignment(field.Affinity, vals[0].Value)
	} else {
		// sparse pairs of values? then its a sparse record of literals
		if k, e := run.GetKindByName(field.Type); e != nil {
			err = e
		} else {
			rec := rt.NewRecord(k)
			if e := readRecord(run, rec, vals); e != nil {
				err = e
			} else {
				ret = rt.RecordOf(rec)
			}
		}
	}
	return
}

// autocreates default sub records if need be.
func readRecord(run *Runner, rec *rt.Record, vs []query.ValueData) (err error) {
	for _, vd := range vs {
		if e := readRecordPart(run, rec, vd); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func readRecordPart(run *Runner, rec *rt.Record, vd query.ValueData) (err error) {
	pos := dot.MakeValueCursor(run, rt.RecordOf(rec))
	// follow the dots
	path := strings.Split(vd.Path, ".")
	for len(path) > 1 {
		part, rest := path[0], path[1:]
		if next, e := pos.GetAtField(part); e != nil {
			err = e
			break
		} else {
			pos, path = next, rest
		}
	}
	if err == nil {
		// access the container of the last value manually
		rec := pos.CurrentValue().Record()
		i := rec.FieldIndex(path[0])
		ft := rec.Field(i)
		if l, e := run.decode.DecodeField(ft.Affinity, vd.Value, ft.Type); e != nil {
			err = e // fix: record field type?!
		} else if v, e := l.GetLiteralValue(run); e != nil {
			err = e
		} else {
			err = rec.SetIndexedField(i, v)
		}
	}
	return
}

// do the pairs contain a dotted path?
// if so then they are *all* literals
// -- only the root level allows evals
// because that's all that can be queried for
// (dotted paths live in core, not in the runtime interface )
func dotted(vals []query.ValueData) bool {
	return len(vals) > 2 || len(vals[0].Path) > 0
}
