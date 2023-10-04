package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
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

// return whether the indexed field has ever been written to.
func (d *Record) HasValue(i int) (ret bool) {
	return d.values[i] != nil
}

// GetNamedField picks a value or trait from this record.
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
		// fix? weave doesnt flag these as class trait
		// we could add that, or put the name of the aspect they came from ( so class stays a "kind" )
	}
	return
}

// GetIndexedField panics if out of range.
// note: traits are never indexed fields ( although their aspect is )
// fix? GetIndexedField writes defaults into the record if there was no value.
func (d *Record) GetIndexedField(i int) (ret Value, err error) {
	// is the stored value valid? return it
	if fv, ft := d.values[i], d.kind.fields[i]; fv != nil {
		ret = fv
	} else {
		// first try set up default aspects
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			for _, a := range d.kind.aspects {
				if a.Name == ft.Name {
					// first trait:
					nv := StringFrom(a.Traits[0], a.Name)
					ret, d.values[i] = nv, nv
					break
				}
			}
		}
		// fallback to other fields:
		if ret == nil {
			if nv, e := NewDefaultValue(ft.Affinity, ft.Type); e != nil {
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
				err = errutil.Fmt("error setting trait: couldn't determine the meaning of %q %s", field, val.Affinity())
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
	if val == nil {
		// probably a terrible idea: resets status of the value to default
		// ( ex. for event interrupt )
		d.values[i] = val
	} else {
		var okay bool
		ft, aff, cls := d.kind.fields[i], val.Affinity(), val.Type()
		switch aff {
		// the most flexible: anything can fit into anything.
		case affine.Bool, affine.Number, affine.Text, affine.TextList, affine.NumList:
			okay = aff == ft.Affinity

		// the least flexible: exact matches are needed.
		case affine.Record, affine.RecordList:
			okay = aff == ft.Affinity && cls == ft.Type

		}
		if !okay {
			err = errutil.Fmt("couldnt set field %s ( %s of type %q ) with val %s of type %q", ft.Name, ft.Affinity, ft.Type, aff, cls)
		} else {
			d.values[i] = val
		}
	}
	return
}
