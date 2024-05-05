package rt

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
)

// Record - provides low level access to named field/values pairs.
// The fields of a record are defined by its kind.
type Record struct {
	*Kind
	values []Value
}

func NewRecord(k *Kind) *Record {
	// values are placeholder; zero values generated on demand
	return &Record{Kind: k, values: make([]Value, k.NumField())}
}

// a record without a named kind
func NewAnonymousRecord(fields []Field) *Record {
	anon := NewKind(nil, fields, nil)
	return NewRecord(anon)
}

// helper for trying to create a record using a name from the passed typeset.
func NewRecordByName(ks Kinds, cls string) (ret *Record, err error) {
	if k, e := ks.GetKindByName(cls); e != nil {
		err = e
	} else {
		ret = NewRecord(k)
	}
	return
}

// GetNamedField picks a value or trait from this record.
func (d *Record) GetNamedField(field string) (ret Value, err error) {
	// note: the field is a trait when the field that was found doesnt match the field requested
	if i := d.FieldIndex(field); i < 0 {
		err = UnknownField(d.Name(), field)
	} else if v, e := d.GetIndexedField(i); e != nil {
		err = e
	} else if ft := d.Field(i); ft.Name == field {
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
	if fv, ft := d.values[i], d.Field(i); fv != nil {
		ret = fv
	} else {
		// first try set up default aspects
		// fix? i dont love this, but it does make sense that enumeration have the least value as a default
		// note: qna objects aren't represented by records, so they dont hit this
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			for _, at := range d.aspects {
				if at.Name == ft.Name {
					// first trait:
					nv := StringFrom(at.Traits[0], at.Name)
					ret, d.values[i] = nv, nv
					break
				}
			}
		}
		// fallback to other fields:
		if ret == nil {
			if nv, e := ZeroField(ft.Affinity, ft.Type, i); e != nil {
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
	if i := d.FieldIndex(field); i < 0 {
		err = UnknownField(d.Name(), field)
	} else {
		// set a normal field ( when the name doesn't match the request: the caller requested a trait )
		if ft := d.Field(i); ft.Name == field {
			err = d.SetIndexedField(i, val)
		} else {
			// set the aspect to the value of the requested trait
			if yes := val.Affinity() == affine.Bool && val.Bool(); !yes {
				err = fmt.Errorf("error setting trait: couldn't determine the meaning of %q %s", field, val.Affinity())
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
		ft, aff, cls := d.Field(i), val.Affinity(), val.Type()
		switch aff {
		// the most flexible: anything can fit into anything.
		case affine.Bool, affine.Number, affine.Text, affine.TextList, affine.NumList:
			okay = aff == ft.Affinity

		// the least flexible: exact matches are needed.
		case affine.Record, affine.RecordList:
			okay = aff == ft.Affinity && cls == ft.Type

		}
		if !okay {
			err = fmt.Errorf("couldnt set field %s ( %s of type %q ) with val %s of type %q", ft.Name, ft.Affinity, ft.Type, aff, cls)
		} else {
			d.values[i] = val
		}
	}
	return
}
