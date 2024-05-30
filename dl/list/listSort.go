package list

import (
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

func (op *ListSortNumbers) Execute(run rt.Runtime) (err error) {
	if e := op.sortByNum(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListSortText) Execute(run rt.Runtime) (err error) {
	if e := op.sortByText(run); e != nil {
		err = CmdError(op, e)
	}
	return
}

func (op *ListSortNumbers) sortByNum(run rt.Runtime) (err error) {
	if at, e := assign.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else {
		byField := inflect.Normalize(op.ByField)
		switch listAff := vs.Affinity(); listAff {
		case affine.RecordList:
			err = sortRecords(run, vs.Records(), byField, affine.Num, op.numSorter)

		case affine.TextList:
			err = sortObjects(run, vs.Strings(), byField, affine.Num, op.numSorter)

		default:
			err = errutil.New("number sort not implemented for", listAff)
		}
	}
	return
}

func (op *ListSortText) sortByText(run rt.Runtime) (err error) {
	if at, e := assign.GetReference(run, op.Target); e != nil {
		err = e
	} else if vs, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(vs, affine.TextList); e != nil {
		err = e
	} else {
		name := inflect.Normalize(op.ByField)
		switch listAff := vs.Affinity(); listAff {
		case affine.TextList:
			err = sortObjects(run, vs.Strings(), name, affine.Text, op.textSorter)

		case affine.RecordList:
			err = sortRecords(run, vs.Records(), name, affine.Text, op.textSorter)

		default:
			err = errutil.New("text sort not implemented for", listAff)
		}
	}
	return
}

func (op *ListSortNumbers) numSorter(run rt.Runtime, a, b rt.Value) (ret bool, err error) {
	aa, bb := a.Float(), b.Float()
	if descending, e := safe.GetOptionalBool(run, op.Descending, false); e != nil {
		err = e
	} else if !descending.Bool() {
		ret = aa < bb
	} else {
		ret = bb < aa
	}
	return
}

func (op *ListSortText) textSorter(run rt.Runtime, a, b rt.Value) (ret bool, err error) {
	aa, bb := a.String(), b.String()
	if sensitive, e := safe.GetOptionalBool(run, op.UsingCase, false); e != nil {
		err = e
	} else if !sensitive.Bool() {
		aa, bb = strings.ToLower(aa), strings.ToLower(bb)
	}
	if descending, e := safe.GetOptionalBool(run, op.Descending, false); e != nil {
		err = e
	} else if !descending.Bool() {
		ret = aa < bb
	} else {
		ret = bb < aa
	}
	return
}

type compareFn func(run rt.Runtime, a, b rt.Value) (bool, error)

func sortRecords(run rt.Runtime, src []*rt.Record, field string, aff affine.Affinity, cmp compareFn) (err error) {
	sort.Slice(src, func(i, j int) (ret bool) {
		ret = i < j // provisionally
		a, b := src[i], src[j]
		if aa, e := unpackRecord(a, field, aff); e != nil {
			err = e
		} else if bb, e := unpackRecord(b, field, aff); e != nil {
			err = e
		} else {
			ret, err = cmp(run, aa, bb)
		}
		return
	})
	return
}

func sortObjects(run rt.Runtime, src []string, field string, aff affine.Affinity, cmp compareFn) (err error) {
	sort.Slice(src, func(i, j int) (ret bool) {
		ret = i < j // provisionally
		if a, e := run.GetField(meta.ObjectId, src[i]); e != nil {
			err = errutil.Append(err, e)
		} else if b, e := run.GetField(meta.ObjectId, src[j]); e != nil {
			err = errutil.Append(err, e)
		} else if aa, e := run.GetField(a.String(), field); e != nil {
			err = errutil.Append(err, e)
		} else if bb, e := run.GetField(b.String(), field); e != nil {
			err = errutil.Append(err, e)
		} else if e := safe.Check(aa, aff); e != nil {
			err = errutil.Append(err, e)
		} else if e := safe.Check(bb, aff); e != nil {
			err = errutil.Append(err, e)
		} else if b, e := cmp(run, aa, bb); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = b
		}
		return
	})
	return
}

func unpackRecord(src *rt.Record, field string, aff affine.Affinity) (ret rt.Value, err error) {
	if v, e := src.GetNamedField(field); e != nil {
		err = e
	} else if e := safe.Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}
