package call

import (
	"testing"
)

func TestTriggers(t *testing.T) {
	t.Run("trigger once", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     F(3),
			Trigger: &TriggerOnce{},
		}); b != 0b00000100 { // bits are read left to right
			t.Fatalf("mismatch %b", b)
		}
	})
	t.Run("trigger cycle", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     F(3),
			Trigger: &TriggerCycle{},
		}); b != 0b00100100 {
			t.Fatalf("mismatch %b", b)
		}
	})
	t.Run("trigger switch", func(t *testing.T) {
		if b := triggerTest(&CallTrigger{
			Name:    t.Name(),
			Num:     F(3),
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
