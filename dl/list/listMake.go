package list

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *MakeNumList) GetNumList(run rt.Runtime) (ret rt.Value, err error) {
	vs := make([]float64, len(op.Values))
	for i, a := range op.Values {
		if v, e := safe.GetNum(run, a); e != nil {
			err = e
		} else {
			vs[i] = v.Float()
		}
	}
	if err == nil {
		ret = rt.FloatsOf(vs)
	}
	return
}

func (op *MakeTextList) GetTextList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.makeList(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *MakeTextList) makeList(run rt.Runtime) (ret rt.Value, err error) {
	vs := make([]string, len(op.Values))
	for i, a := range op.Values {
		if v, e := safe.GetText(run, a); e != nil {
			err = e
		} else {
			vs[i] = v.String()
		}
	}
	if err == nil {
		ret = rt.StringsOf(vs)
	}
	return
}

func (op *MakeRecordList) GetRecordList(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.makeList(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *MakeRecordList) makeList(run rt.Runtime) (ret rt.Value, err error) {
	if cnt := len(op.List); cnt == 0 {
		err = errors.New("record list is empty")
	} else {
		var subtype string
		vs := make([]*rt.Record, cnt)
		for i, a := range op.List {
			if v, e := safe.GetRecord(run, a); e != nil {
				err = e
			} else if vt := v.Type(); len(subtype) == 0 {
				subtype = vt
			} else if vt != subtype {
				err = fmt.Errorf("record %d of type %q not %q", i, vt, subtype)
				break
			} else {
				vs[i] = v.Record()
			}
		}
		if err == nil {
			ret = rt.RecordsFrom(vs, subtype)
		}
	}
	return
}
