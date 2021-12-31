package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Affinity helps implement Assignment, but always returns ""
func (op *GetVar) Affinity() affine.Affinity { return "" }

func (op *GetVar) String() string {
	return op.Name.String()
}

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *GetVar) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, "")
}

func (op *GetVar) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Bool)
}

func (op *GetVar) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Number)
}

func (op *GetVar) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Record)
}

func (op *GetVar) GetText(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Text)
}

func (op *GetVar) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.NumList)
}

func (op *GetVar) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.TextList)
}

func (op *GetVar) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.RecordList)
}

func (op *GetVar) getVar(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Name.String(), aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
