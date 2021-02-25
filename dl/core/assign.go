package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// Let assigns a value to a local variable.
type Assign struct {
	Var  Variable      `if:"selector"`
	From rt.Assignment `if:"selector=be"`
}

// FromBool - implements Assignment
type FromBool struct {
	Val rt.BoolEval `if:"selector"`
}

// FromNum - implements Assignment
type FromNum struct {
	Val rt.NumberEval `if:"selector"`
}

// FromText - implements Assignment
type FromText struct {
	Val rt.TextEval `if:"selector"`
}

// FromRecord - implements Assignment
type FromRecord struct {
	Val rt.RecordEval `if:"selector"`
}

// FromNumbers - implements Assignment
type FromNumbers struct {
	Vals rt.NumListEval `if:"selector"`
}

// FromTexts - implements Assignment
type FromTexts struct {
	Vals rt.TextListEval `if:"selector"`
}

// FromRecords - implements Assignment
type FromRecords struct {
	Vals rt.RecordListEval `if:"selector"`
}

func (*Assign) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "Let: Assigns a variable to a value.",
		Fluent: &composer.Fluid{Name: "let", Role: composer.Command},
	}
}

func (op *Assign) Execute(run rt.Runtime) (err error) {
	if v, e := safe.GetAssignedValue(run, op.From); e != nil {
		err = cmdError(op, e)
	} else if e := run.SetField(object.Variables, op.Var.String(), v); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (*FromBool) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Bool: Assigns the calculated boolean value.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromBool) Affinity() affine.Affinity {
	return affine.Bool
}
func (op *FromBool) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetBool(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromNum) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Number: Assigns the calculated number.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromNum) Affinity() affine.Affinity {
	return affine.Number
}
func (op *FromNum) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetNumber(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromText) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Text: Assigns the calculated piece of text.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromText) Affinity() affine.Affinity {
	return affine.Text
}
func (op *FromText) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if val, e := safe.GetText(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = val
	}
	return
}

func (*FromRecord) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Record: Assigns the calculated record.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromRecord) Affinity() affine.Affinity {
	return affine.Record
}
func (op *FromRecord) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.GetRecord(run, op.Val); e != nil {
		err = cmdError(op, e)
	} else {
		ret = obj
	}
	return
}

func (*FromNumbers) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Numbers: Assigns the calculated numbers.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromNumbers) Affinity() affine.Affinity {
	return affine.NumList
}
func (op *FromNumbers) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetNumList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (*FromTexts) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Texts: Assigns the calculated texts.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromTexts) Affinity() affine.Affinity {
	return affine.TextList
}
func (op *FromTexts) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if vals, e := safe.GetTextList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = vals
	}
	return
}

func (*FromRecords) Compose() composer.Spec {
	return composer.Spec{
		Group:  "variables",
		Desc:   "From Records: Assigns the calculated records.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}
func (op *FromRecords) Affinity() affine.Affinity {
	return affine.RecordList
}
func (op *FromRecords) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if objs, e := safe.GetRecordList(run, op.Vals); e != nil {
		err = cmdError(op, e)
	} else {
		ret = objs
	}
	return
}
