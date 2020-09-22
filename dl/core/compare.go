package core

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

type CompareNum struct {
	A  rt.NumberEval
	Is Comparator
	B  rt.NumberEval
}

type CompareText struct {
	A  rt.TextEval
	Is Comparator
	B  rt.TextEval
}

func (*CompareNum) Compose() composer.Spec {
	return composer.Spec{
		Name:  "compare_num",
		Group: "logic",
		Spec:  "{a:number_eval} {is:comparator} {b:number_eval}",
		Desc:  "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
	}
}

func (comp *CompareNum) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := rt.GetNumber(run, comp.A); e != nil {
		err = errutil.New("CompareNum.A", e)
	} else if tgt, e := rt.GetNumber(run, comp.B); e != nil {
		err = errutil.New("CompareNum.B", e)
	} else {
		ret = compare(comp.Is, int(src-tgt))
	}
	return
}

func compare(comp Comparator, d int) (ret bool) {
	switch cmp := comp.Compare(); {
	case d == 0:
		ret = (cmp & Compare_EqualTo) != 0
	case d < 0:
		ret = (cmp & Compare_LessThan) != 0
	case d > 0:
		ret = (cmp & Compare_GreaterThan) != 0
	}
	return
}

func (*CompareText) Compose() composer.Spec {
	return composer.Spec{
		Name:  "compare_text",
		Group: "logic",
		Desc:  "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
		Spec:  "{a:text_eval} {is:comparator} {b:text_eval}",
	}
}

func (comp *CompareText) GetBool(run rt.Runtime) (ret bool, err error) {
	if src, e := rt.GetText(run, comp.A); e != nil {
		err = errutil.New("CompareText.A", e)
	} else if tgt, e := rt.GetText(run, comp.B); e != nil {
		err = errutil.New("CompareText.B", e)
	} else {
		ret = compare(comp.Is, strings.Compare(src, tgt))
	}
	return
}
