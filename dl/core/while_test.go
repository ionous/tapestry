package core

import (
	"errors"
	"testing"

	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

func TestLoopBreak(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			B(true), MakeActivity(
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &AtLeast{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				// &Next{},
				&Assign{N("j"), &FromNum{&SumOf{V("j"), I(1)}}},
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
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
				&ChooseAction{
					If: &CompareNum{V("i"), &AtLeast{}, I(4)},
					Do: MakeActivity(
						&Break{},
					),
				},
				&Next{},
				&Assign{N("j"), &FromNum{&SumOf{V("j"), I(1)}}},
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
				&Assign{N("i"), &FromNum{&SumOf{V("i"), I(1)}}},
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
	case field == "i" && target == meta.Variables:
		ret = g.IntOf(k.i)
	case field == "j" && target == meta.Variables:
		ret = g.IntOf(k.j)
	default:
		panic("unexpected get")
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *loopRuntime) SetField(target, field string, v g.Value) (err error) {
	switch {
	case field == "i" && target == meta.Variables:
		k.i = v.Int()
	case field == "j" && target == meta.Variables:
		k.j = v.Int()
	default:
		panic("unexpected set")
	}
	return
}
