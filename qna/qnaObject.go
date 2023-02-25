package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// expects field to be a normalized name already.
func (run *Runner) setObjectField(obj qdb.NounInfo, field string, newValue g.Value) (err error) {
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
func (run *Runner) getObjectField(obj qdb.NounInfo, field string) (ret g.Value, err error) {
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

func (run *Runner) setFieldCache(obj qdb.NounInfo, field g.Field, val g.Value) (err error) {
	// fix: convert when appropriate.
	if aff := val.Affinity(); aff != field.Affinity {
		err = errutil.Fmt(`mismatched affinity "%s.%s(%s)" writing %s`, obj, field.Name, field.Affinity, aff)
	} else {
		key := makeKey(obj.Domain, obj.Id, field.Name)
		run.nounValues[key] = cachedValue{v: g.CopyValue(val)}
	}
	return
}

func (run *Runner) getFieldCache(obj qdb.NounInfo, field g.Field) (ret g.Value, err error) {
	if c, e := run.nounValues.cache(func() (ret interface{}, err error) {
		// note: in the original version of this, we queried *all* fields
		// ( unioning in those with traits, and those without defaults )
		if b, e := run.qdb.NounValue(obj.Id, field.Name); e != nil {
			err = e
		} else {
			// fields be empty, have literal values, or dynamic values.
			if len(b) == 0 {
				if v, e := g.NewDefaultValue(run, field.Affinity, field.Type); e != nil {
					err = e
				} else {
					ret = v
				}
			} else if isEvalLike := b[0] == '{'; !isEvalLike {
				if v, e := parseLiteral(field.Affinity, field.Type, b); e != nil {
					err = e
				} else {
					ret = v
				}
			} else {
				if a, e := parseEval(field.Affinity, b, run.signatures); e != nil {
					err = e
				} else {
					ret = a // a is an assignment which generate values.
				}
			}
		}
		return
	}, obj.Domain, obj.Id, field.Name); e != nil {
		err = e
	} else {
		switch c := c.(type) {
		case g.Value:
			ret = c
		case assign.Assignment:
			// evaluate the assignment to get the current value
			if v, e := assign.GetSafeAssignment(run, c); e != nil {
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
