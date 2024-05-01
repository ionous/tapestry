package list

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *MakeNumList) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	vs := make([]float64, len(op.Values))
	for i, a := range op.Values {
		if v, e := safe.GetNumber(run, a); e != nil {
			err = e
		} else {
			vs[i] = v.Float()
		}
	}
	if err == nil {
		ret = g.FloatsOf(vs)
	}
	return
}

func (op *MakeNumList) makeList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.makeList(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *MakeTextList) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.makeList(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *MakeTextList) makeList(run rt.Runtime) (ret g.Value, err error) {
	vs := make([]string, len(op.Values))
	for i, a := range op.Values {
		if v, e := safe.GetText(run, a); e != nil {
			err = e
		} else {
			vs[i] = v.String()
		}
	}
	if err == nil {
		ret = g.StringsOf(vs)
	}
	return
}

func (op *MakeRecordList) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.makeList(run); e != nil {
		err = CmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *MakeRecordList) makeList(run rt.Runtime) (ret g.Value, err error) {
	if subtype, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if subtype := subtype.String(); len(subtype) == 0 {
		err = errutil.New("expected a valid record name")
	} else if k, e := run.GetKindByName(subtype); e != nil {
		err = errutil.Fmt("expected a known record name, got %q", subtype)
	} else {
		vs := make([]*g.Record, len(op.Values))
		for i, a := range op.Values {
			if v, e := safe.GetRecord(run, a); e != nil {
				err = e
			} else if v.Type() != k.Name() {
				err = errutil.Fmt("record %d of type %q not %q", i, v.Type(), subtype)
			} else if rec, ok := v.Record(); !ok {
				err = errutil.Fmt("record %d of type %q was nil", i, v.Type())
			} else {
				vs[i] = rec
			}
		}
		if err == nil {
			ret = g.RecordsFrom(vs, subtype)
		}
	}
	return
}
