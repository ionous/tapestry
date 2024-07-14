package list

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *ListSort) Execute(run rt.Runtime) (err error) {
	if e := op.sortByText(run); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *ListSort) sortByText(run rt.Runtime) (err error) {
	if at, e := safe.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if opt, e := op.getOptions(run); e != nil {
		err = e
	} else {
		switch listAff := vs.Affinity(); listAff {
		case affine.NumList:
			opt.SortFloats(vs.Floats())

		case affine.TextList:
			if len(opt.Field) == 0 {
				opt.SortStrings(vs.Strings())
			} else if kind := opt.Kind; len(kind) == 0 {
				err = fmt.Errorf("field %s specified but the text list has no particular type", opt.Field)
			} else {
				err = opt.sortObjects(run, kind, vs.Strings())
			}

		case affine.RecordList:
			if len(opt.Field) == 0 {
				err = fmt.Errorf("a field is required to sort a list of records")
			} else if kind := vs.Type(); len(kind) == 0 {
				err = fmt.Errorf("field %s specified but the records have no particular type", opt.Field)
			} else {
				err = opt.SortRecords(run, vs.Records())
			}

		default:
			err = fmt.Errorf("text sort not implemented for %s", listAff)
		}
	}
	return
}

func (op *ListSort) getOptions(run rt.Runtime) (ret Sorter, err error) {
	if kind, e := safe.GetOptionalText(run, op.KindName, ""); e != nil {
		err = e
	} else if field, e := safe.GetOptionalText(run, op.FieldName, ""); e != nil {
		err = e
	} else if caseSensitive, e := safe.GetOptionalBool(run, op.Case, false); e != nil {
		err = e
	} else if descending, e := safe.GetOptionalBool(run, op.Descending, false); e != nil {
		err = e
	} else {
		ret = Sorter{
			Kind:          inflect.Normalize(kind.String()),
			Field:         inflect.Normalize(field.String()),
			CaseSensitive: caseSensitive.Bool(),
			Descending:    descending.Bool(),
		}
	}
	return
}

type Sorter struct {
	Kind          string
	Field         string
	CaseSensitive bool // preserves case while sorting strings
	Descending    bool
}

func (opt Sorter) SortFloats(vs []float64) {
	slices.SortFunc(vs, opt.compareFloats)
}

func (opt Sorter) SortStrings(str []string) {
	slices.SortFunc(str, opt.compareStrings)
}

func (opt Sorter) compareStrings(a, b string) (ret int) {
	if !opt.CaseSensitive {
		a, b = strings.ToLower(a), strings.ToLower(b)
	}
	if v := cmp.Compare(a, b); !opt.Descending {
		ret = v
	} else {
		ret = -v
	}
	return
}

func (opt Sorter) compareFloats(a, b float64) (ret int) {
	if v := cmp.Compare(a, b); !opt.Descending {
		ret = v
	} else {
		ret = -v
	}
	return
}

func (opt Sorter) compareInts(a, b int) (ret int) {
	if v := cmp.Compare(a, b); !opt.Descending {
		ret = v
	} else {
		ret = -v
	}
	return
}

func (opt Sorter) SortRecords(ks rt.Kinds, src []*rt.Record) (err error) {
	if k, e := ks.GetKindByName(opt.Kind); e != nil {
		err = e
	} else if idx := k.FieldIndex(opt.Field); idx < 0 {
		err = fmt.Errorf("unknown sort field %s for kind %s", opt.Field, opt.Kind)
	} else {
		ft := k.Field(idx)
		switch aff := ft.Affinity; aff {
		case affine.Text:
			opt.sortRecordStrings(src, idx)
		case affine.Num:
			opt.sortRecordFloats(src, idx)
		default:
			err = fmt.Errorf("can't sort record %s field %s of type %s", k.Name(), opt.Field, aff)
		}
	}
	return
}

func (opt Sorter) sortRecordStrings(src []*rt.Record, idx int) {
	slices.SortFunc(src, func(a, b *rt.Record) int {
		aa, _ := a.GetIndexedField(idx)
		bb, _ := b.GetIndexedField(idx)
		return opt.compareStrings(aa.String(), bb.String())
	})
}

func (opt Sorter) sortRecordFloats(src []*rt.Record, idx int) {
	slices.SortFunc(src, func(a, b *rt.Record) int {
		aa, _ := a.GetIndexedField(idx)
		bb, _ := b.GetIndexedField(idx)
		return opt.compareFloats(aa.Float(), bb.Float())
	})
}

func (opt Sorter) sortObjects(run rt.Runtime, kind string, names []string) (err error) {
	if k, e := run.GetKindByName(kind); e != nil {
		err = e
	} else if idx := k.FieldIndex(opt.Field); idx < 0 {
		err = fmt.Errorf("unknown sort field %s for kind %s", opt.Field, kind)
	} else {
		ft := k.Field(idx)
		if aspect, e := getAspect(run, ft); e != nil {
			err = e
		} else if vs, e := getObjectValues(run, names, ft, aspect); e != nil {
			err = e
		} else {
			nv := objectValues{opt, names, vs}
			switch aff := ft.Affinity; aff {
			case affine.Num:
				sort.Sort(byFloat{nv})
			case affine.Text:
				if aspect != nil {
					sort.Sort(byState{nv})
				} else {
					sort.Sort(byString{nv})
				}
			default:
				err = fmt.Errorf("can't sort record %s field %s of type %s", k.Name(), opt.Field, aff)
			}
		}
	}
	return
}

func getAspect(run rt.Runtime, ft rt.Field) (ret *rt.Kind, err error) {
	if ft.Affinity == affine.Text && len(ft.Type) != 0 && ft.Type == ft.Name {
		if k, e := run.GetKindByName(ft.Type); e != nil {
			err = e
		} else if k.Implements(kindsOf.Aspect.String()) {
			ret = k
		}
	}
	return
}

// to sort object by a field's value
// take a snapshot of that field for every object.
type objectValues struct {
	sort   Sorter   // the rules for sorting
	names  []string // the actual slice memory being sorted
	values []any    // either rt.Value, or an int representation of an object state
}

func (nv objectValues) Len() int {
	return len(nv.names)
}

func (nv objectValues) Swap(i, j int) {
	nv.names[i], nv.names[j] = nv.names[j], nv.names[i]
	nv.values[i], nv.values[j] = nv.values[j], nv.values[i]
}

type byFloat struct{ objectValues }
type byString struct{ objectValues }
type byState struct{ objectValues }

func (s byString) Less(i, j int) bool {
	a, b := s.values[i].(rt.Value), s.values[j].(rt.Value)
	return s.sort.compareStrings(a.String(), b.String()) < 0
}

func (s byFloat) Less(i, j int) bool {
	a, b := s.values[i].(rt.Value), s.values[j].(rt.Value)
	return s.sort.compareFloats(a.Float(), b.Float()) < 0
}

func (s byState) Less(i, j int) bool {
	a, b := s.values[i], s.values[j]
	return s.sort.compareInts(a.(int), b.(int)) < 0
}

// to sort object by a field's value
// take a snapshot of that field for every object.
func getObjectValues(run rt.Runtime, names []string, ft rt.Field, aspect *rt.Kind) (ret []any, err error) {
	out := make([]any, len(names))
	for i, name := range names {
		if id, e := run.GetField(meta.ObjectId, name); e != nil {
			err = e
			break
		} else if v, e := run.GetField(id.String(), ft.Name); e != nil {
			err = e
			break
		} else if e := safe.Check(v, ft.Affinity); e != nil {
			err = e
			break
		} else {
			if aspect != nil {
				out[i] = aspect.FieldIndex(v.String())
			} else {
				out[i] = v
			}
		}
	}
	if err == nil {
		ret = out
	}
	return
}
