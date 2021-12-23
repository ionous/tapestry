package generic

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"github.com/ionous/errutil"
)

// Record - provides low level access to named field/values pairs.
// The fields of a record are defined by its kind.
type Record struct {
	kind   *Kind
	values []Value
}

func (d *Record) Kind() *Kind {
	return d.kind
}

func (d *Record) Type() string {
	return d.kind.name
}

// GetNamedField picks out a value from this record.
func (d *Record) GetNamedField(field string) (ret Value, err error) {
	switch k := d.kind; field {
	case meta.ObjectName:
		err = errutil.New("records don't have names")

	case meta.ObjectKind, meta.ObjectKinds:
		ret = StringOf(d.kind.name)

	default:
		// note: the field is a trait when the field that was found doesnt match the field requested
		if i := k.FieldIndex(field); i < 0 {
			err = UnknownField(k.name, field)
		} else if v, e := d.GetIndexedField(i); e != nil {
			err = e
		} else if ft := k.fields[i]; ft.Name == field {
			ret = v
		} else {
			trait := v.String()
			ret = BoolFrom(trait == field, "" /*"trait"*/)
			// fix? the assembler doesnt flag these as class trait
			// we could add that, or put the name of the aspect they came from ( so class stays a "kind" )
		}
	}
	return
}

// GetIndexedField can't ask for traits, only their aspects.
func (d *Record) GetIndexedField(i int) (ret Value, err error) {
	if fv, ft := d.values[i], d.kind.fields[i]; fv != nil {
		ret = fv
	} else {
		if cls, ok := ft.isAspectLike(); ok {
			// if we're asking for an aspect, the default value will be the string of its first trait
			if k, e := d.kind.kinds.GetKindByName(cls); e != nil {
				err = e
			} else {
				firstTrait := k.Field(0)                          // first trait is the default
				nv := StringFrom(firstTrait.Name, "" /*"trait"*/) // fix? should the class be set to something intersting?
				ret, d.values[i] = nv, nv
			}
		} else {
			if nv, e := NewDefaultValue(d.kind.kinds, ft.Affinity, ft.Type); e != nil {
				err = e
			} else {
				ret, d.values[i] = nv, nv
			}
		}
	}
	return
}

// SetNamedField - pokes the passed value into the record.
// Unlike the Value interface, this doesnt panic and it doesnt copy values.
func (d *Record) SetNamedField(field string, val Value) (err error) {
	k := d.kind // note: the field is a trait when the field that was found doesnt match the field requested
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField(k.name, field)
	} else if ft := k.fields[i]; ft.Name == field {
		err = d.SetIndexedField(i, val)
	} else if yes := val.Affinity() == affine.Bool && val.Bool(); !yes {
		err = errutil.Fmt("error setting trait: couldn't determine the meaning of %q %s %v", field, val.Affinity(), val)
	} else {
		// set the aspect to the value of the requested trait
		d.values[i] = StringFrom(field, "" /*"aspect"*/)
	}
	return
}

// SetIndexedField - note this doesn't handle trait translation.
// Unlike the Value interface, this doesnt panic and it doesnt copy values.
func (d *Record) SetIndexedField(i int, val Value) (err error) {
	ft := d.kind.fields[i]
	if a, t := val.Affinity(), val.Type(); !matchTypes(d.kind.kinds, ft.Affinity, ft.Type, a, t) {
		err = errutil.Fmt("couldnt set field %s ( %s of type %s ) because val %s of type %s doesnt match", ft.Name, ft.Affinity, ft.Type, a, t)
	} else {
		d.values[i] = val
	}
	return
}

// fix: if type includes path then this Kinds isnt necessary, which sure would be nce
// ( of course, for normal records Implements() doesnt matter anyway. records dont have hierarchy )
func matchTypes(ks Kinds, fa affine.Affinity, ft string, va affine.Affinity, vt string) (okay bool) {
	if fa == va {
		recordLike := fa == affine.Object || fa == affine.Record || fa == affine.RecordList
		if !recordLike {
			okay = true
		} else if vt == ft {
			okay = true // direct match
		} else if vk, e := ks.GetKindByName(vt); e == nil {
			// a field takes: cats
			// my value is things, animals, cats, tigers.
			// so we should ask: does the value's path contains the field's kind.
			okay = vk.Implements(ft)
		}
	}
	return
}
