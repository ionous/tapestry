package call

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestTriggers(t *testing.T) {
	t.Run("trigger once", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     literal.F(3),
			Trigger: &TriggerOnce{},
		}); b != 0b00000100 { // bits are read left to right
			t.Fatalf("mismatch %b", b)
		}
	})
	t.Run("trigger cycle", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     literal.F(3),
			Trigger: &TriggerCycle{},
		}); b != 0b00100100 {
			t.Fatalf("mismatch %b", b)
		}
	})
	t.Run("trigger switch", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     literal.F(3),
			Trigger: &TriggerSwitch{},
		}); b != 0b11111100 {
			t.Fatalf("mismatch %b", b)
		}
	})
}

func triggerTest(op *CallTrigger) (ret int) {
	run := seqTest{counters: make(map[string]int)}
	for i := 0; i < 8; i++ {
		if v, e := op.GetBool(&run); e != nil {
			panic(e)
		} else if v.Bool() {
			ret |= 1 << i
		}
	}
	return
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type seqTest struct {
	baseRuntime
	counters map[string]int
}

func (m *seqTest) Random(inclusiveMin, exclusiveMax int) int {
	return (exclusiveMax-inclusiveMin)/2 + inclusiveMin
}

func (m *seqTest) GetField(target, field string) (ret rt.Value, err error) {
	if target != meta.Counter {
		err = rt.UnknownField(target, field)
	} else {
		v := m.counters[field]
		ret = rt.IntOf(v)
	}
	return
}

func (m *seqTest) SetField(target, field string, value rt.Value) (err error) {
	if target != meta.Counter {
		err = rt.UnknownField(target, field)
	} else {
		m.counters[field] = value.Int()
	}
	return
}
