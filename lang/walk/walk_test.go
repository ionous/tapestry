package walk_test

import (
	r "reflect"
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/lang/walk"
)

type Pod struct {
	A string
	B string
}

type Outer struct {
	X Pod
	Y string
	Z string
	W Pod
}

type Slice struct {
	Els []Pod
}

type Slot struct {
	P any
}

func TestPod(t *testing.T) {
	var pod Pod
	Check(t, walk.Walk(r.ValueOf(pod)),
		"{", "A", "Value", "B", "Value", "}",
	)
}

func TestOuter(t *testing.T) {
	var outer Outer
	Check(t, walk.Walk(r.ValueOf(outer)),
		"{",
		"X", "Flow",
		"{", "A", "Value", "B", "Value", "}",
		"Y", "Value",
		"Z", "Value",
		"W", "Flow",
		"{", "A", "Value", "B", "Value", "}",
		"}",
	)
}

func TestSlice(t *testing.T) {
	slice := Slice{make([]Pod, 2)}
	Check(t, walk.Walk(r.ValueOf(slice)),
		"{", "Els", "Flow", "[",
		"0", "{", "A", "Value", "B", "Value", "}",
		"1", "{", "A", "Value", "B", "Value", "}",
		"]", "}",
	)
	empty := Slice{}
	Check(t, walk.Walk(r.ValueOf(empty)),
		"{", "Els", "Flow", "[",
		"]", "}",
	)
}

func TestSlot(t *testing.T) {
	empty := Slot{}
	Check(t, walk.Walk(r.ValueOf(empty)),
		"{", "P", "Slot",
		"<", ">",
		"}")
	filled := Slot{new(Pod)}
	Check(t, walk.Walk(r.ValueOf(filled)),
		"{", "P", "Slot",
		"<", "Pod", // reports on what type fills the slot
		"{", "A", "Value", "B", "Value", "}", // then the flow
		">", "}")
}

func Check(t *testing.T, it walk.Walker, parts ...string) {
	lc := listChecker{parts: parts}
	lc.flow(t, it)

	if lc.idx < len(parts) {
		t.Logf("ended early, missing: %v", parts[lc.idx:])
		t.Fail()
	}
}

func (lc *listChecker) flow(t *testing.T, it walk.Walker) {
	lc.check(t, "{")
	for it.Next() {
		f := it.Field()
		lc.check(t, f.Name(), f.SpecType().String())
		switch f.SpecType() {
		case walk.Flow:
			if f.Repeats() {
				lc.repeatFlow(t, it.Walk())
			} else {
				lc.flow(t, it.Walk())
			}

		case walk.Slot:
			if f.Repeats() {
				panic("should really have a test for this")
			} else {
				lc.slot(t, it.Walk())
			}
		}
	}
	lc.check(t, "}")
	return
}

func (lc *listChecker) slot(t *testing.T, it walk.Walker) {
	lc.check(t, "<")
	for i := 0; it.Next(); i++ {
		lc.check(t, it.Value().Type().Name())
		lc.flow(t, it.Walk())
	}
	lc.check(t, ">")
}

func (lc *listChecker) repeatFlow(t *testing.T, it walk.Walker) {
	lc.check(t, "[")
	for i := 0; it.Next(); i++ {
		lc.check(t, strconv.Itoa(i))
		lc.flow(t, it.Walk())
	}
	lc.check(t, "]")
	return
}

type listChecker struct {
	idx   int
	parts []string
}

// check the passed strings, one at a time.
func (lc *listChecker) check(t *testing.T, strs ...string) {
	for _, got := range strs {
		t.Log(got)
		if lc.idx >= len(lc.parts) {
			t.Log("out of range")
			t.Fail()
		} else if next := lc.parts[lc.idx]; next != got {
			t.Logf("mismatch at %d; expected %q got %q", lc.idx, next, got)
			t.Fail()
		} else {
			lc.idx++
		}
	}
	return
}
