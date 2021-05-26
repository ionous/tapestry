package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// Affinity helps implement Assignment, but always returns ""
func (op *Var) Affinity() affine.Affinity { return "" }

func (op *Var) String() string {
	return op.Name
}

// GetAssignedValue implements Assignment so we can SetXXX values from variables without a FromXXX statement in between.
func (op *Var) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, "")
}

func (op *Var) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Bool)
}

func (op *Var) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Number)
}

func (op *Var) GetRecord(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Record)
}

func (op *Var) GetText(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.Text)
}

func (op *Var) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.NumList)
}

func (op *Var) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.TextList)
}

func (op *Var) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	return op.getVar(run, affine.RecordList)
}

func (op *Var) getVar(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := safe.CheckVariable(run, op.Name, aff); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
