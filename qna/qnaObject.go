package qna

import (
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// expects field to be a normalized name already.
func (run *Runner) setObjectField(obj query.NounInfo, field string, newValue g.Value) (err error) {
	// tbd: cache the kind in the object info?
	// or even... cache the ( last n ) field info into "obj.field"?
	if kind, e := run.getKind(obj.Kind); e != nil {
		err = e
	} else if fieldIndex := kind.FieldIndex(field); fieldIndex < 0 {
		err = g.UnknownField(obj.String(), field)
	} else {
		fieldData := kind.Field(fieldIndex)
		if fieldData.Name == field {
			// the field we found is the name we asked for: its a normal field.
			err = run.setFieldCache(obj, fieldData, newValue)
		} else {
			// when the name differs, we must have found the aspect for a trait.
			// FIX: should we transform the value so that it has type of the aspect?
			// FIX: records dont have opposite day so this seems ... unfair.
			// FIX: im also curious about aspects that only have one trait, or blank ( nothing ).
			if aff := newValue.Affinity(); aff != affine.Bool {
				err = errutil.New("can only set a trait with booleans, have", aff)
			} else if trait, e := oppositeDay(run, fieldData.Name, field, newValue.Bool()); e != nil {
				err = e
			} else {
				// set the aspect to the value of the requested trait
				traitValue := g.StringFrom(trait, fieldData.Type)
				err = run.setFieldCache(obj, fieldData, traitValue)
			}
		}
	}
	return
}

// expects field to be a normalized name already.
func (run *Runner) getObjectField(obj query.NounInfo, field string) (ret g.Value, err error) {
	// tbd: cache the kind in the object info?
	// or even... cache the ( last n ) field info into "obj.field"?
	if kind, e := run.getKind(obj.Kind); e != nil {
		err = e
	} else if fieldIndex := kind.FieldIndex(field); fieldIndex < 0 {
		err = g.UnknownField(obj.String(), field)
	} else {
		fieldData := kind.Field(fieldIndex)
		if v, e := run.getFieldCache(obj, fieldData); e != nil {
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
			ret = g.BoolOf(field == v.String())
		}
	}
	return
}

func (run *Runner) setFieldCache(obj query.NounInfo, field g.Field, val g.Value) (err error) {
	// fix: convert when appropriate.
	if aff := val.Affinity(); aff != field.Affinity {
		err = errutil.Fmt(`mismatched affinity "%s.%s(%s)" writing %s`, obj, field.Name, field.Affinity, aff)
	} else {
		key := makeKey(obj.Domain, obj.Id, field.Name)
		run.nounValues[key] = cachedValue{v: g.CopyValue(val)}
	}
	return
}

func (run *Runner) getFieldCache(obj query.NounInfo, field g.Field) (ret g.Value, err error) {
	if c, e := run.nounValues.cache(func() (ret any, err error) {
		// note: in the original version of this, we queried *all* fields
		// ( unioning in those with traits, and those without defaults )
		if pairs, e := run.query.NounValues(obj.Id, field.Name); e != nil {
			err = e
		} else if len(pairs) == 0 { // fields be empty, have literal values, or dynamic values.
			ret, e = readFields(run, field, pairs)
		} else {
			// tbd: needed for create record, wish there was a nicer way
			// at the very least it'd be nice if the decoder could hold this
			// but then we'd have to set the runtime into it? its loopy and confusing.
			var ks g.Kinds = run
			ret, err = g.NewDefaultValue(ks, field.Affinity, field.Type)
		}
		return
	}, obj.Domain, obj.Id, field.Name); e != nil {
		err = e
	} else {
		switch c := c.(type) {
		case g.Value:
			ret = c
		case rt.Assignment:
			// evaluate the assignment to get the current value
			if v, e := safe.GetAssignment(run, c); e != nil {
				err = e
			} else {
				ret, err = safe.AutoConvert(run, field, v)
			}
		default:
			err = errutil.Fmt("unexpected type in object cache %T", c)
		}
	}
	return
}

// return can be an assignment ( which gets evaluated )
// or a literal value ( a fixed value )
func readFields(run *Runner, field g.Field, pairs []string) (ret any, err error) {
	if !dotted(pairs) {
		// a single top level value? then its an assignment
		value := pairs[1]
		ret, err = run.decode.DecodeAssignment([]byte(value))
	} else {
		// sparse pairs of values? then its a sparse record of literals
		if k, e := run.GetKindByName(field.Type); e != nil {
			err = e
		} else {
			rec := k.NewRecord()
			if e := makeRecords(run, rec, pairs); e != nil {
				err = e
			} else {
				ret = g.RecordOf(rec)
			}
		}
	}
	return
}

func makeRecords(run *Runner, rec *g.Record, pairs []string) (err error) {
	for i, cnt := 0, len(pairs); i < cnt; i += 2 {
		if e := fillRecord(run, rec, pairs[i], pairs[i+1]); e != nil {
			err = e
			break
		}
	}
	return
}

func fillRecord(run *Runner, rec *g.Record, fullpath, value string) (err error) {
	for path := fullpath; len(path) > 0; {
		part, rest := dotscan(path)
		k := rec.Kind() // has aff if needed
		if i := k.FieldIndex(part); i < 0 {
			err = errutil.New("error")
		} else {
			if len(rest) == 0 {
				field := k.Field(i)
				// FIX: how does fieldType actually get recorded!?!
				if l, e := run.decode.DecodeField([]byte(value), field.Affinity, field.Type); e != nil {
					err = e
				} else if v, e := l.GetLiteralValue(run); e != nil {
					err = e
				} else {
					err = rec.SetIndexedField(i, v)
				}
				break // all done regardless
			} else {
				// a part ending with a dot is a record:
				// the Get() will auto-create the value --
				// fix, future: this is questionable requires rec to know Kinds.
				// the caller could surely handle that ( ex. Dotted and this ) when needed.
				if v, e := rec.GetIndexedField(i); e != nil {
					err = e
					break
				} else if v.Affinity() != affine.Record {
					err = errutil.New("error")
					break
				} else {
					rec, path = v.Record(), rest
				}
			}
		}
	}
	return
}

// do the pairs contain a dotted path?
// if so then they are *all* literals
// -- only the root level allows evals
// because that's all that can be queried for
// (dotted paths live in core, not in the runtime interface )
func dotted(pairs []string) (ret bool) {
	// any more than two elements requires it
	if cnt := len(pairs); cnt > 2 {
		ret = true
	} else if cnt == 2 {
		path := pairs[0]
		_, rhs := dotscan(path)
		ret = len(rhs) > 0 // a dot always has two parts left and right
	}
	return
}

// return the string up to the next dot, and everything after.
func dotscan(str string) (lhs, rhs string) {
	lhs = str // provisionally
	for accum := 0; accum < len(str); {
		if r, n := utf8.DecodeRuneInString(str[accum:]); r == utf8.RuneError {
			break // all done or error
		} else if r == '.' {
			lhs, rhs = str[:accum], str[:accum+n]
			break
		} else {
			accum += n
		}
	}
	return
}
