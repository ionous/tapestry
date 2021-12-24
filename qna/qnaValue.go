package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

func decodeAssignment(a affine.Affinity, prog []byte, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	if e := cin.Decode(rt.Assignment_Slot{&ret}, prog, signatures); e != nil {
		err = e
	}
	return
}

func zeroValue(ks g.Kinds, a affine.Affinity, t string) (ret g.Value, err error) {
	switch a {
	case affine.Bool:
		ret = g.BoolFrom(false, t)
	case affine.Number:
		ret = g.FloatFrom(0, t)
	case affine.Text:
		ret = g.StringFrom("", t)
	case affine.NumList:
		ret = g.FloatsFrom(nil, t)
	case affine.TextList:
		ret = g.StringsFrom(nil, t)
	case affine.Record:
		if k, e := ks.GetKindByName(t); e != nil {
			err = e
		} else {
			ret = g.RecordOf(k.NewRecord())
		}
	case affine.RecordList:
		if k, e := ks.GetKindByName(t); e != nil {
			err = e // verify that the kind exists
		} else {
			ret = g.RecordsFrom(nil, k.Name())
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}

// the parsing of data from the database depends on how the runtime interprets the data
// for instance, is `1` a number or a string with a number?
// note too, value in the db can usually contain either literal values, or evals (ex. text vs templates producing text )
func readValue(a affine.Affinity, rawValue []byte, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	switch a {
	case affine.Bool:
		var v rt.BoolEval
		if e := cin.Decode(rt.BoolEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromBool{v}
		}
	case affine.Number:
		var v rt.NumberEval
		if e := cin.Decode(rt.NumberEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromNum{v}
		}
	case affine.Text:
		var v rt.TextEval
		if e := cin.Decode(rt.TextEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromText{v}
		}
	case affine.NumList:
		var v rt.NumListEval
		if e := cin.Decode(rt.NumListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromNumbers{v}
		}
	case affine.TextList:
		var v rt.TextListEval
		if e := cin.Decode(rt.TextListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromTexts{v}
		}
	case affine.Record:
		var v rt.RecordEval
		if e := cin.Decode(rt.RecordEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromRecord{v}
		}
	case affine.RecordList:
		var v rt.RecordListEval
		if e := cin.Decode(rt.RecordListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromRecords{v}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}
