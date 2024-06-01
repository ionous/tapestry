package logic_test

import (
	"errors"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestLoopBreak(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&logic.While{
			True: literal.B(true), Exe: []rt.Execute{
				&object.SetValue{
					Target: object.Variable("i"),
					Value:  &assign.FromNum{Value: &math.AddValue{A: object.Variable("i"), B: literal.I(1)}}},
				&logic.ChooseBranch{
					Condition: &math.CompareNum{A: object.Variable("i"), Is: math.C_Comparison_AtLeast, B: literal.I(4)},
					Exe: []rt.Execute{
						&logic.Break{},
					},
				},
				&object.SetValue{
					Target: object.Variable("j"),
					Value:  &assign.FromNum{Value: &math.AddValue{A: object.Variable("j"), B: literal.I(1)}}},
			}},
	); e != nil {
		t.Fatal("failed to run:", e)
	} else if run.i != 4 && run.j != 3 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopNext(t *testing.T) {
	var run loopRuntime
	if e := safe.Run(&run,
		&logic.While{
			True: literal.B(true), Exe: []rt.Execute{
				&object.SetValue{
					Target: object.Variable("i"),
					Value:  &assign.FromNum{Value: &math.AddValue{A: object.Variable("i"), B: literal.I(1)}}},
				&logic.ChooseBranch{
					Condition: &math.CompareNum{A: object.Variable("i"), Is: math.C_Comparison_AtLeast, B: literal.I(4)},
					Exe: []rt.Execute{
						&logic.Break{},
					},
				},
				&logic.Continue{},
				&object.SetValue{
					Target: object.Variable("j"),
					Value:  &assign.FromNum{Value: &math.AddValue{A: object.Variable("j"), B: literal.I(1)}}},
			}},
	); e != nil {
		t.Fatal(e)
	} else if run.i != 4 && run.j != 0 {
		t.Fatal("bad counters", run.i, run.j)
	}
}

func TestLoopInfinite(t *testing.T) {
	logic.MaxLoopIterations = 100
	var run loopRuntime
	if e := safe.Run(&run,
		&logic.While{
			True: literal.B(true), Exe: []rt.Execute{
				&object.SetValue{
					Target: object.Variable("i"),
					Value:  &assign.FromNum{Value: &math.AddValue{A: object.Variable("i"), B: literal.I(1)}}},
			}},
	); !errors.Is(e, logic.MaxLoopIterations) {
		t.Fatal(e)
	} else if run.i != 100 {
		t.Fatal("bad counters", run.i, run.j)
	} else {
		t.Log("ok, error is expected:", e)
	}
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type loopRuntime struct {
	baseRuntime
	i, j int
}

func (k *loopRuntime) GetField(target, field string) (ret rt.Value, err error) {
	switch {
	case field == "i" && target == meta.Variables:
		ret = rt.IntOf(k.i)
	case field == "j" && target == meta.Variables:
		ret = rt.IntOf(k.j)
	default:
		panic("unexpected get")
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *loopRuntime) SetField(target, field string, v rt.Value) (err error) {
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
