package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"github.com/ionous/errutil"
)

// every primitive value is its own unique instance.
// records are as pointers, and lists as pointers to slices;
// their data is therefore shared across refValues.
// pointers to slices allows for in-place grow, append, etc.
type refValue struct {
	a affine.Affinity
	t string
	i interface{}
}

var _ Value = (*refValue)(nil)

func (n refValue) Affinity() affine.Affinity {
	return n.a
}

func (n refValue) Type() string {
	return n.t
}

func (n refValue) Bool() bool {
	return n.i.(bool)
}

func (n refValue) Float() (ret float64) {
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
		panic(n.a.String() + " is not a number")
	}
	return
}

func (n refValue) Int() (ret int) {
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
		panic(n.a.String() + " is not a number")
	}
	return
}

func (n refValue) String() string {
	return n.i.(string)
}

func (n refValue) Record() *Record {
	return n.i.(*Record)
}

func (n refValue) Floats() (ret []float64) {
	vp := n.i.(*[]float64)
	return *vp
}

func (n refValue) Strings() (ret []string) {
	vp := n.i.(*[]string)
	return *vp
}

func (n refValue) Records() (ret []*Record) {
	vp := n.i.(*[]*Record)
	return *vp
}

func (n refValue) Len() (ret int) {
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
		panic(n.a.String() + " is not measurable")
	}
	return
}

func (n refValue) Index(i int) (ret Value) {
	switch vp := n.i.(type) {
	case *[]float64:
		ret = FloatFrom((*vp)[i], n.t)
	case *[]string:
		ret = StringFrom((*vp)[i], n.t)
	case *[]*Record:
		ret = RecordOf((*vp)[i])
	default:
		panic(n.a.String() + " is not indexable")
	}
	return
}

func (n refValue) FieldByName(f string) (ret Value, err error) {
	name := lang.SpecialUnderscore(f)
	if v, e := n.Record().GetNamedField(name); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (n refValue) SetFieldByName(f string, v Value) (err error) {
	rec := n.Record()
	name := lang.SpecialUnderscore(f)
	newVal := CopyValue(v)
	return rec.SetNamedField(name, newVal)
}

func (n refValue) SetIndex(i int, v Value) (err error) {
	switch vp := n.i.(type) {
	case *[]float64:
		(*vp)[i] = v.Float()
	case *[]string:
		if len(n.t) > 0 && n.t != v.Type() {
			err = errutil.Fmt("SetIndex(%s) doesnt match value(%s)", n.t, v.Type())
		} else {
			(*vp)[i] = v.String()
		}
	case *[]*Record:
		if n.t != v.Type() {
			err = errutil.New("record types dont match")
		} else {
			n := copyRecord(v.Record())
			(*vp)[i] = n
		}
	default:
		panic(n.a.String() + " is not index writable")
	}
	return
}

// Slices copies a chunk out of a list
func (n refValue) Slice(i, j int) (ret Value, err error) {
	if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
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
			panic(n.a.String() + " is not sliceable")
		}
	}
	return
}

// Splice replaces a range of values
func (n refValue) Splice(i, j int, add Value) (ret Value, err error) {
	if i < 0 {
		err = Underflow{i, 0}
	} else if cnt := n.Len(); j > cnt {
		err = Overflow{j, cnt}
	} else if i > j {
		err = errutil.New("bad range", i, j)
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
				err = errutil.Fmt("Splice(%s) doesnt match value(%s)", n.t, add.Type())
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
				err = errutil.New("record types dont match")
			} else {
				els := (*vp)
				// move the record pointers
				// no need to copy the record values
				// only one list will have the pointers at a time
				cut := make([]*Record, j-i)
				copy(cut, els[i:j])
				// make a list out of one or more records from add
				ins := copyRecords(normalizeRecords(add))
				// read from els before adding to els to avoid stomping overlapping memory.
				(*vp) = append(els[:i], append(ins, els[j:]...)...)
				// return our cut pointers
				ret = RecordsFrom(cut, n.t)
			}
		default:
			panic(n.a.String() + " is not spliceable")
		}
	}
	return
}

func (n refValue) Appends(add Value) (err error) {
	switch n.a {
	case affine.NumList:
		vp := n.i.(*[]float64)
		ins := normalizeFloats(add)
		(*vp) = append((*vp), ins...)

	case affine.TextList:
		if len(n.t) > 0 && n.t != add.Type() {
			err = errutil.Fmt("Appends(%s) doesnt match value(%s)", n.t, add.Type())
		} else {
			vp := n.i.(*[]string)
			ins := normalizeStrings(add)
			(*vp) = append((*vp), ins...)
		}

	case affine.RecordList:
		vp := n.i.(*[]*Record)
		if n.t != add.Type() {
			err = errutil.New("record types dont match")
		} else {
			ins := copyRecords(normalizeRecords(add))
			(*vp) = append((*vp), ins...)
		}

	default:
		panic(n.a.String() + " is not appendable")
	}
	return
}
