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
		log.Panicf("%s doesn't have fields", n.a)
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
		log.Panicf("%s doesn't have fields", n.a)
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
func (n variant) Splice(i, j int, add Value) (ret Value, err error) {
	if i < 0 {
		err = fmt.Errorf("slice at %d can't be negative", i)
	} else if cnt := n.Len(); j > cnt {
		err = fmt.Errorf("splice at %d out of range %d", j, cnt)
	} else if i > j {
		err = fmt.Errorf("invalid splice range %d > %d", i, j)
	} else {
		switch n.a {
		case affine.NumList:
			vp := n.i.(*[]float64)
			els := (*vp)
			cut := copyFloats(els[i:j])
			ins := normalizeFloats(add)
			(*vp) = append(els[:i], append(ins, els[j:]...)...)
			ret = FloatsOf(cut)

		case affine.TextList:
			if len(n.t) > 0 && n.t != add.Type() {
				err = fmt.Errorf("Splice(%s) doesnt match value(%s)", n.t, add.Type())
			} else {
				vp := n.i.(*[]string)
				els := (*vp)
				cut := copyStrings(els[i:j])
				ins := normalizeStrings(add)
				(*vp) = append(els[:i], append(ins, els[j:]...)...)
				ret = StringsOf(cut)
			}

		case affine.RecordList:
			vp := n.i.(*[]*Record)
			if n.t != add.Type() {
				err = errors.New("record types dont match")
			} else if src, e := normalizeRecords(add); e != nil {
				err = e // // make a list out of one or more records from add
			} else {
				els := (*vp)
				// move the record pointers
				// no need to copy the record values
				// only one list will have the pointers at a time
				cut := make([]*Record, j-i)
				copy(cut, els[i:j])
				ins := copyRecords(src)
				// read from els before adding to els to avoid stomping overlapping memory.
				(*vp) = append(els[:i], append(ins, els[j:]...)...)
				// return our cut pointers
				ret = RecordsFrom(cut, n.t)
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
		} else if els, e := normalizeRecords(add); e != nil {
			err = e
		} else {
			ins := copyRecords(els)
			(*vp) = append((*vp), ins...)
		}

	default:
		log.Panicf("%s is not appendable", n.a)
	}
	return
}
