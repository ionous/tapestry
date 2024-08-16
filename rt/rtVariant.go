package rt

import (
	"errors"
	"fmt"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
)

// variant implements the Value interface.
// every primitive value is its own unique instance.
// records are as pointers, and lists as pointers to slices;
// their data is therefore shared across variants.
// pointers to slices allows for in-place grow, append, etc.
type variant struct {
	a affine.Affinity
	t string
	i any
}

var _ Value = (*variant)(nil)

func (n variant) Affinity() affine.Affinity {
	return n.a
}

func (n variant) Type() string {
	return n.t
}

func (n variant) Bool() bool {
	return n.i.(bool)
}

func (n variant) Float() (ret float64) {
	// see also: MakeLiteral
	switch v := n.i.(type) {
	case float64:
		ret = v
	case float32:
		ret = float64(v)
	case int:
		ret = float64(v)
	case int64:
		ret = float64(v)
	default:
		log.Panicf("%s is not a number", n.a)
	}
	return
}

func (n variant) Int() (ret int) {
	switch v := n.i.(type) {
	case int:
		ret = v
	case int64:
		ret = int(v)
	case float32:
		ret = int(v)
	case float64:
		ret = int(v)
	default:
		log.Panicf("%s is not a number", n.a)
	}
	return
}

func (n variant) String() (ret string) {
	if a, ok := n.i.(string); ok {
		ret = a
	} else {
		ret = fmt.Sprintf("<%T %v>", n.i, n.i)
	}
	return
}

func (n variant) Record() *Record {
	rec := n.i.(*Record)
	return rec
}

func (n variant) Floats() (ret []float64) {
	vp := n.i.(*[]float64)
	return *vp
}

func (n variant) Strings() (ret []string) {
	vp := n.i.(*[]string)
	return *vp
}

func (n variant) Records() (ret []*Record) {
	vp := n.i.(*[]*Record)
	return *vp
}

func (n variant) Len() (ret int) {
	switch vp := n.i.(type) {
	case string:
		ret = len(vp)
	case *[]float64:
		ret = len(*vp)
	case *[]string:
		ret = len(*vp)
	case *[]*Record:
		ret = len(*vp)
	default:
		log.Panicf("%s is not measurable", n.a)
	}
	return
}

func (n variant) Index(i int) (ret Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		ret = FloatFrom((*vp)[i], n.t)
	case *[]string:
		ret = StringFrom((*vp)[i], n.t)
	case *[]*Record:
		ret = RecordOf((*vp)[i])
	default:
		log.Panicf("%s is not indexable", n.a)
	}
	return
}

func (n variant) FieldByName(name string) (ret Value, err error) {
	switch rec := n.i.(type) {
	case *Record:
		ret, err = rec.GetNamedField(name)
	default:
		log.Panicf("trying to get field %q, but %s types don't have fields", name, n.a)
	}
	return
}

// copies the incoming value
func (n variant) SetFieldByName(name string, v Value) (err error) {
	switch rec := n.i.(type) {
	case *Record:
		newVal := CopyValue(v)
		return rec.SetNamedField(name, newVal)
	default:
		log.Panicf("trying to set field %q, but %s types don't have fields", name, n.a)
	}
	return
}

func (n variant) SetIndex(i int, v Value) (err error) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp)[i] = v.Float()
	case *[]string:
		if len(n.t) > 0 && n.t != v.Type() {
			err = fmt.Errorf("SetIndex(%s) doesnt match value(%s)", n.t, v.Type())
		} else {
			(*vp)[i] = v.String()
		}
	case *[]*Record:
		if n.t != v.Type() {
			err = errors.New("record types dont match")
		} else {
			rec := v.Record()
			n := copyRecordValues(rec)
			(*vp)[i] = n
		}
	default:
		log.Panicf("%s is not index writable", n.a)
	}
	return
}

// Slices copies a chunk out of a list
func (n variant) Slice(i, j int) (ret Value, err error) {
	if i < 0 {
		err = fmt.Errorf("slice at %d can't be negative", i)
	} else if cnt := n.Len(); j > cnt {
		err = fmt.Errorf("slice at %d out of range %d", j, cnt)
	} else if i > j {
		err = fmt.Errorf("invalid slice range %d > %d", i, j)
	} else {
		switch n.a {
		case affine.NumList:
			vp := n.i.(*[]float64)
			ret = FloatsFrom(copyFloats((*vp)[i:j]), n.Type())

		case affine.TextList:
			vp := n.i.(*[]string)
			ret = StringsFrom(copyStrings((*vp)[i:j]), n.Type())

		case affine.RecordList:
			vp := n.i.(*[]*Record)
			ret = RecordsFrom(copyRecords((*vp)[i:j]), n.Type())

		default:
			log.Panicf("%s is not sliceable", n.a)
		}
	}
	return
}

// Splice replaces a range of values
func (n variant) Splice(i, j int, add Value, cutList *Value) (err error) {
	if i < 0 {
		err = fmt.Errorf("splice at %d can't be negative", i)
	} else if cnt := n.Len(); j > cnt {
		err = fmt.Errorf("splice at %d out of range %d", j, cnt)
	} else if i > j {
		err = fmt.Errorf("invalid splice range %d > %d", i, j)
	} else {
		switch n.a {
		case affine.NumList:
			vp := n.i.(*[]float64)
			els := (*vp)
			if cutList != nil {
				cut := copyFloats(els[i:j])
				*cutList = FloatsOf(cut)
			}
			ins := normalizeFloats(add)
			(*vp) = append(els[:i], append(ins, els[j:]...)...)

		case affine.TextList:
			// if this list has a type, then the other list must have a type and they must match;
			// or the other list must have no type and be empty. ( kind of restrictive... not sure what's best )
			if add != nil && len(n.t) > 0 && (n.t != add.Type() && !(add.Type() == "" && add.Len() == 0)) {
				err = fmt.Errorf("Splice additions of %q don't match text type %q", n.t, add.Type())
			} else {
				vp := n.i.(*[]string)
				els := (*vp)
				if cutList != nil {
					cut := copyStrings(els[i:j])
					*cutList = StringsOf(cut)
				}
				ins := normalizeStrings(add)
				(*vp) = append(els[:i], append(ins, els[j:]...)...)
			}

		case affine.RecordList:
			vp := n.i.(*[]*Record)
			if add != nil && n.t != add.Type() {
				err = fmt.Errorf("Splice additions of %q don't match record type %q", n.t, add.Type())
			} else {
				els := (*vp)
				if cutList != nil {
					// move the record pointers
					// no need to copy the record values
					// only one list will have the pointers at a time
					cut := make([]*Record, j-i)
					copy(cut, els[i:j])
					*cutList = RecordsFrom(cut, n.t)
				}
				ins := copyRecords(normalizeRecords(add))
				// read from els before adding to els to avoid stomping overlapping memory.
				(*vp) = append(els[:i], append(ins, els[j:]...)...)
			}
		default:
			log.Panicf("%s is not spliceable", n.a)
		}
	}
	return
}

func (n variant) Appends(add Value) (err error) {
	switch n.a {
	case affine.NumList:
		vp := n.i.(*[]float64)
		ins := normalizeFloats(add)
		(*vp) = append((*vp), ins...)

	case affine.TextList:
		if len(n.t) > 0 && n.t != add.Type() {
			err = fmt.Errorf("Appends(%s) doesnt match value(%s)", n.t, add.Type())
		} else {
			vp := n.i.(*[]string)
			ins := normalizeStrings(add)
			(*vp) = append((*vp), ins...)
		}

	case affine.RecordList:
		vp := n.i.(*[]*Record)
		if n.t != add.Type() {
			err = errors.New("record types dont match")
		} else {
			ins := copyRecords(normalizeRecords(add))
			(*vp) = append((*vp), ins...)
		}

	default:
		log.Panicf("%s is not appendable", n.a)
	}
	return
}
