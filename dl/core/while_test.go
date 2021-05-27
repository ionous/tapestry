package core

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func TestLoopBreak(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			B(true), MakeActivity(
				&Assign{"i", &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &GreaterOrEqual{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				// &Next{},
				&Assign{"j", &FromNum{&SumOf{V("j"), I(1)}}},
			)},
	); e != nil {
		t.Fatal(e)
	} else if run.i != 4 && run.j != 3 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopNext(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			B(true), MakeActivity(
				&Assign{"i", &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &GreaterOrEqual{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				&Next{},
				&Assign{"j", &FromNum{&SumOf{V("j"), I(1)}}},
			)},
	); e != nil {
		t.Fatal(e)
	} else if run.i != 4 && run.j != 0 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopInfinite(t *testing.T) {
	MaxLoopIterations = 100
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			B(true), MakeActivity(
				&Assign{"i", &FromNum{&SumOf{V("i"), I(1)}}},
			)},
	); !errors.Is(e, MaxLoopIterations) {
		t.Fatal(e)
	} else if run.i != 100 {
		t.Fatal("bad counters", run.i, run.j)
	} else {
		t.Log("ok, error is expected:", e)
	}
}

type loopRuntime struct {
	baseRuntime
	i, j int
}

func (k *loopRuntime) GetField(target, field string) (ret g.Value, err error) {
	switch {
	case field == "i" && target == object.Variables:
		ret = g.IntOf(k.i)
	case field == "j" && target == object.Variables:
		ret = g.IntOf(k.j)
	default:
		panic("unexpected get")
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *loopRuntime) SetField(target, field string, v g.Value) (err error) {
	switch {
	case field == "i" && target == object.Variables:
		k.i = v.Int()
	case field == "j" && target == object.Variables:
		k.j = v.Int()
	default:
		panic("unexpected set")
	}
	return
}

func V(i string) *Var           { return &Var{i} }
func T(s string) *TextValue     { return &TextValue{value.Text(s)} }
func I(n int) rt.NumberEval     { return &NumValue{float64(n)} }
func N(n float64) rt.NumberEval { return &NumValue{n} }
func B(b bool) rt.BoolEval      { return &BoolValue{b} }
