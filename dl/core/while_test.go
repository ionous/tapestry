package core

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func TestLoopBreak(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			True: B(true), Exe: MakeActivity(
				&assign.SetValue{
					Target: Variable("i"),
					Value:  &assign.FromNumber{Value: &SumOf{A: Variable("i"), B: I(1)}}},
				&ChooseBranch{
					If: &CompareNum{A: Variable("i"), Is: AtLeast, B: I(4)},
					Exe: MakeActivity(
						&Break{},
					),
				},
				&assign.SetValue{
					Target: Variable("j"),
					Value:  &assign.FromNumber{Value: &SumOf{A: Variable("j"), B: I(1)}}},
			)},
	); e != nil {
		t.Fatal("failed to run:", e)
	} else if run.i != 4 && run.j != 3 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopNext(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&While{
			True: B(true), Exe: MakeActivity(
				&assign.SetValue{
					Target: Variable("i"),
					Value:  &assign.FromNumber{Value: &SumOf{A: Variable("i"), B: I(1)}}},
				&ChooseBranch{
					If: &CompareNum{A: Variable("i"), Is: AtLeast, B: I(4)},
					Exe: MakeActivity(
						&Break{},
					),
				},
				&Continue{},
				&assign.SetValue{
					Target: Variable("j"),
					Value:  &assign.FromNumber{Value: &SumOf{A: Variable("j"), B: I(1)}}},
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
			True: B(true), Exe: MakeActivity(
				&assign.SetValue{
					Target: Variable("i"),
					Value:  &assign.FromNumber{Value: &SumOf{A: Variable("i"), B: I(1)}}},
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
