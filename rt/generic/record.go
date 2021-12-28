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
		case affine.Record, affine.RecordList, affine.TextList:
			okay = cls == ft.Type

		// various compatible text references are allowed.
		case affine.Text:
			if str := val.String(); ft.Type == cls {
				okay = true // accept as is.
			} else if len(ft.Type) == 0 {
				val = StringOf(str) // downgrades to blank
				okay = true
			} else {
				// is the source untyped? then maybe we can upgrade it to a trait.
				if len(cls) == 0 {
					if ft.isAspectLike() {
						if aspect := findAspect(str, d.kind.traits); aspect == ft.Type {
							val = StringFrom(str, ft.Type) // upgrade the value to the aspect.
							okay = true
						}
					}
				} else {
					// a field wants only "cats"; my value is "things, animals, cats, tigers".
					// so we ask: does the value implement field:
					if vk, e := d.kind.kinds.GetKindByName(cls); e == nil {
						// ( downgrade the value's type to fit the field? for now, lets not. )
						okay = vk.Implements(ft.Type)
					}
				}
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
