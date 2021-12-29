package generic

import (
	"git.sr.ht/~ionous/iffy/affine"
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
	// note: the field is a trait when the field that was found doesnt match the field requested
	k := d.kind
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
	return
}

// GetIndexedField can't ask for traits, only their aspects.
func (d *Record) GetIndexedField(i int) (ret Value, err error) {
	if fv, ft := d.values[i], d.kind.fields[i]; fv != nil {
		ret = fv
	} else {
		if ft.isAspectLike() {
			// if we're asking for an aspect, the default value will be the string of its first trait
			if k, e := d.kind.kinds.GetKindByName(ft.Type); e != nil {
				err = e
			} else {
				firstTrait := k.Field(0) // first trait is the default
				nv := StringFrom(firstTrait.Name, k.Name())
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
	k := d.kind
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField(k.name, field)
	} else {
		// set a normal field ( when the name doesn't match the request: the caller requested a trait )
		if ft := k.fields[i]; ft.Name == field {
			err = d.SetIndexedField(i, val)
		} else {
			// set the aspect to the value of the requested trait
			if yes := val.Affinity() == affine.Bool && val.Bool(); !yes {
				err = errutil.Fmt("error setting trait: couldn't determine the meaning of %q %s %v", field, val.Affinity(), val)
			} else {
				d.values[i] = StringFrom(field, ft.Type)
			}
		}
	}
	return
}

// SetIndexedField - note this doesn't handle trait translation.
// Unlike the Value interface, this doesnt panic and it doesnt copy values.
func (d *Record) SetIndexedField(i int, val Value) (err error) {
	var okay bool
	ft, aff, cls := d.kind.fields[i], val.Affinity(), val.Type()
	if aff == ft.Affinity {
		switch ft.Affinity {
		// the most flexible: anything can fit into anything.
		case affine.Bool, affine.Number, affine.NumList:
			okay = true

		// the least flexible: exact matches are needed.
		case affine.Record, affine.RecordList:
			okay = cls == ft.Type

		case affine.Text:
			if change := textConvert(d.kind.kinds, ft, val); change < 0 {
				okay = true // accept as is
			} else if change > 0 {
				val = StringFrom(val.String(), ft.Type)
				okay = true
			}

		case affine.TextList:
			if change := textConvert(d.kind.kinds, ft, val); change < 0 {
				okay = true
			} else if change > 0 {
				val = StringsFrom(val.Strings(), ft.Type)
				okay = true
			}
		}
	}
	if !okay {
		err = errutil.Fmt("couldnt set field %s ( %s of type %q ) with val %s of type %q", ft.Name, ft.Affinity, ft.Type, aff, cls)
	} else {
		d.values[i] = val
	}
	return
}

// typed text can only fit in certain typed text fields
// return: keep(<0); not convertible(0); converts(>0)
func textConvert(ks Kinds, ft Field, val Value) (ret int) {
	switch cls := val.Type(); {
	// if the classes match exactly: accept as is.
	case ft.Type == cls:
		ret = -1

	// untyped text can accept any src type.
	case len(ft.Type) == 0:
		ret = -1 // keep the value as is.

	// typed text can only accept certain typed src values.
	// ex. a field wants only "cats"; my value (cls) is "things, animals, cats, tigers".
	case len(cls) > 0:
		if vk, e := ks.GetKindByName(cls); e == nil {
			if vk.Implements(ft.Type) {
				ret = -1 // keep the value as is (vs. "downgrading" the value to fit the field.)
			}
		}

	// untyped blank strings can fit into any type of text
	// ( they become typed blank strings )
	case val.Len() == 0:
		ret = 1

	// untyped traits can fit into aspects ( if the trait is valid )
	// this doesnt apply to text lists ( they can never be true aspects )
	case ft.isAspectLike():
		if ak, e := ks.GetKindByName(ft.Type); e == nil {
			if ak.FieldIndex(val.String()) >= 0 {
				ret = 1
			}
		}
	}
	return
}
