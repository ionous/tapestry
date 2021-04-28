package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type Splice struct {
	Var           core.Variable `if:"selector"`
	Start, Remove rt.NumberEval // from start
	Insert        rt.Assignment
}

func (*Splice) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "splice",
		Name:  "list_splice",
		Group: "list",
		Spec:  "Splice {list:text} {at entry%start?number_eval} {removing%remove?number_eval} {inserting%insert?assignment}",
		Desc: `Splice into list: Modify a list by adding and removing elements.
Note: the type of the elements being added must match the type of the list. 
Text cant be added to a list of numbers, numbers cant be added to a list of text.
If the starting index is negative, it will begin that many elements from the end of the array.
If list's length + the start is less than 0, it will begin from index 0.
If the remove count is missing, it removes all elements from the start to the end; 
if it is 0 or negative, no elements are removed.`,
	}
}

func (op *Splice) Execute(run rt.Runtime) (err error) {
	if _, _, e := op.spliceList(run, ""); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Splice) GetNumList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.spliceList(run, affine.NumList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.FloatsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Splice) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if v, _, e := op.spliceList(run, affine.TextList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.StringsOf(nil)
	} else {
		ret = v
	}
	return
}

func (op *Splice) GetRecordList(run rt.Runtime) (ret g.Value, err error) {
	if v, t, e := op.spliceList(run, affine.RecordList); e != nil {
		err = cmdError(op, e)
	} else if v == nil {
		ret = g.RecordsOf(t, nil)
	} else {
		ret = v
	}
	return
}

func (op *Splice) spliceList(run rt.Runtime, aff affine.Affinity) (retVal g.Value, retType string, err error) {
	if els, e := safe.List(run, op.Var.String()); e != nil {
		err = e
	} else if e := safe.Check(els, aff); e != nil {
		err = e
	} else if ins, e := safe.GetAssignedValue(run, op.Insert); e != nil {
		err = e
	} else if !IsAppendable(ins, els) {
		err = insertError{ins, els}
	} else if i, j, e := op.getIndices(run, els.Len()); e != nil {
		err = e
	} else {
		if i >= 0 && j >= i {
			retVal, err = els.Splice(i, j, ins)
		}
		if err == nil {
			retType = els.Type()
		}
	}
	return
}

func (op *Splice) getIndices(run rt.Runtime, cnt int) (reti, retj int, err error) {
	if i, e := safe.GetOptionalNumber(run, op.Start, 0); e != nil {
		err = e
	} else if rng, e := safe.GetOptionalNumber(run, op.Remove, float64(cnt)); e != nil {
		err = e
	} else {
		reti = clipStart(i.Int(), cnt)
		retj = clipRange(reti, rng.Int(), cnt)
	}
	return
}
