package list

import (
	"sort"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
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
	if root, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if els, e := root.GetCheckedValue(run, affine.NumList); e != nil {
		err = e
	} else {
		name := en.Normalize(op.ByField)
		switch listAff := els.Affinity(); listAff {
		case affine.RecordList:
			err = sortRecords(run, els.Records(), name, affine.Number, op.numSorter)

		case affine.TextList:
			err = sortObjects(run, els.Strings(), name, affine.Number, op.numSorter)

		default:
			err = errutil.New("number sort not implemented for", listAff)

			if err == nil {
				root.SetDirty(run)
			}
		}
	}
	return
}

func (op *ListSortText) sortByText(run rt.Runtime) (err error) {
	if root, e := assign.GetRootValue(run, op.Target); e != nil {
		err = e
	} else if els, e := root.GetCheckedValue(run, affine.TextList); e != nil {
		err = e
	} else {
		name := en.Normalize(op.ByField)
		switch listAff := els.Affinity(); listAff {
		case affine.RecordList:
			// fix? would any of this be clearer/smaller if we used els.Index?
			// ( well sort.Slice couldnt work on it directly, but maybe there's a slice index )
			err = sortRecords(run, els.Records(), name, affine.Text, op.textSorter)

		case affine.TextList:
			err = sortObjects(run, els.Strings(), name, affine.Text, op.textSorter)

		default:
			err = errutil.New("text sort not implemented for", listAff)
		}

		if err == nil {
			root.SetDirty(run)
		}
	}
	return
}

func (op *ListSortNumbers) numSorter(run rt.Runtime, a, b g.Value) (ret bool, err error) {
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

func (op *ListSortText) textSorter(run rt.Runtime, a, b g.Value) (ret bool, err error) {
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

type compareFn func(run rt.Runtime, a, b g.Value) (bool, error)

func sortRecords(run rt.Runtime, src []*g.Record, field string, aff affine.Affinity, cmp compareFn) (err error) {
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
