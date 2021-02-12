package list

import (
	"sort"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// SortNumbers implements Sorter
type SortNumbers struct {
	Var     core.Variable `if:"selector=numbers"`
	ByField *SortByField  `if:"selector"`
	Order   `if:"selector"`
}

// SortText implements Sorter
type SortText struct {
	Var     core.Variable `if:"selector=text"`
	ByField *SortByField  `if:"selector"`
	Order   `if:"selector"`
	Case    `if:"selector"`
}

// SortRecords implements Sorter
type SortRecords struct {
	Var   core.Variable `if:"selector=records"`
	Using pattern.PatternName
}

type SortByField struct {
	Name string `if:"selector"`
}

func (op *SortByField) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_by_field",
		Group:  "list",
		Desc:   `Sort numbers: .`,
		Fluent: &composer.Fluid{Name: "byField", Role: composer.Selector},
	}
}

func (op *SortNumbers) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_numbers",
		Group:  "list",
		Desc:   `Sort numbers: .`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}
func (op *SortText) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_text",
		Group:  "list",
		Desc:   `Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}
func (op *SortRecords) Compose() composer.Spec {
	return composer.Spec{
		Name:   "list_sort_using",
		Group:  "list",
		Desc:   `Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.`,
		Fluent: &composer.Fluid{Name: "sort", Role: composer.Command},
	}
}

func (op *SortNumbers) Execute(run rt.Runtime) (err error) {
	if e := op.sortByNum(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SortText) Execute(run rt.Runtime) (err error) {
	if e := op.sortByText(run); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *SortRecords) Execute(rt.Runtime) (err error) {
	return errutil.New("not implemented")
}

func (op *SortNumbers) sortByNum(run rt.Runtime) (err error) {
	if v, e := run.GetField(object.Variables, op.Var.String()); e != nil {
		err = e
	} else {
		if pby := op.ByField; pby == nil {
			err = errutil.New("not implemented")
		} else {
			switch aff := v.Affinity(); aff {
			case affine.RecordList:
				err = sortRecords(run, v.Records(), pby.Name, affine.Number, op.numSorter)

			case affine.TextList:
				err = sortObjects(run, v.Strings(), pby.Name, affine.Number, op.numSorter)

			default:
				err = errutil.New("not implemented")
			}
		}
	}
	return
}

func (op *SortText) sortByText(run rt.Runtime) (err error) {
	if v, e := run.GetField(object.Variables, op.Var.String()); e != nil {
		err = e
	} else {
		if pby := op.ByField; pby == nil {
			err = errutil.New("not implemented")
		} else {
			switch aff := v.Affinity(); aff {
			case affine.RecordList:
				err = sortRecords(run, v.Records(), pby.Name, affine.Text, op.textSorter)

			case affine.TextList:
				err = sortObjects(run, v.Strings(), pby.Name, affine.Text, op.textSorter)

			default:
				err = errutil.New("not implemented")
			}
		}
	}
	return
}

func (op *SortNumbers) numSorter(a, b g.Value) (ret bool) {
	aa, bb := a.Float(), b.Float()
	if !op.Order.Descending() {
		ret = aa < bb
	} else {
		ret = bb < aa
	}
	return
}

func (op *SortText) textSorter(a, b g.Value) (ret bool) {
	aa, bb := a.String(), b.String()
	if !op.Case.Sensitive() {
		aa, bb = strings.ToLower(aa), strings.ToLower(bb)
	}
	if !op.Order.Descending() {
		ret = aa < bb
	} else {
		ret = bb < aa
	}
	return
}

func sortRecords(run rt.Runtime, src []*g.Record, field string, aff affine.Affinity, cmp func(a, b g.Value) bool) (err error) {
	sort.Slice(src, func(i, j int) (ret bool) {
		ret = i < j // provisionally
		a, b := src[i], src[j]
		if aa, e := unpackRecord(a, field, aff); e != nil {
			err = e
		} else if bb, e := unpackRecord(b, field, aff); e != nil {
			err = e
		} else {
			ret = cmp(aa, bb)
		}
		return
	})
	return
}

func sortObjects(run rt.Runtime, src []string, field string, aff affine.Affinity, cmp func(a, b g.Value) bool) (err error) {
	sort.Slice(src, func(i, j int) (ret bool) {
		ret = i < j // provisionally
		if a, e := safe.ObjectFromString(run, src[i]); e != nil {
			err = errutil.Append(err, e)
		} else if b, e := safe.ObjectFromString(run, src[j]); e != nil {
			err = errutil.Append(err, e)
		} else if aa, e := safe.Unpack(a, field, aff); e != nil {
			err = e
		} else if bb, e := safe.Unpack(b, field, aff); e != nil {
			err = e
		} else {
			ret = cmp(aa, bb)
		}
		return
	})
	return
}

func unpackRecord(src *g.Record, field string, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := src.GetNamedField(field); e != nil {
		err = e
	} else if e := safe.Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}
