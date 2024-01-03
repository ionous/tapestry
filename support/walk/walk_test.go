package walk_test

import (
	"errors"
	"fmt"
	r "reflect"
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/walk"
)

func TestPod(t *testing.T) {
	type Pod struct {
		A string
		B string
	}
	if e := Check(t, walk.NewWalker(r.ValueOf(new(Pod)).Elem()),
		"Pod", "Initial",
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
	); e != nil {
		t.Fatal(e)
	}

	type Outer struct {
		X Pod
		Y string
		Z string
		W Pod
	}
	if e := Check(t, walk.NewWalker(r.ValueOf(new(Outer)).Elem()),
		"Outer", "Initial",
		"Outer", "X",
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
		"Outer", "Y",
		"Outer", "Z",
		"Outer", "W",
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
		"Outer", "Terminal",
	); e != nil {
		t.Fatal(e)
	}
}

func TestSlice(t *testing.T) {
	type Pod struct {
		A string
		B string
	}

	if e := Check(t, walk.NewWalker(r.ValueOf(make([]Pod, 2))),
		"", "Initial", // the slice as container
		"", "0", // the pod element; ( the container is still the slice )
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
		"", "1",
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
		"", "Terminal",
	); e != nil {
		t.Fatal(e)
	}
}

func TestSlot(t *testing.T) {
	type Pod struct {
		A string
		B string
	}
	type Outer struct {
		Slot any
	}

	empty := Outer{}
	if e := Check(t, walk.NewWalker(r.ValueOf(empty)),
		"Outer", "Initial",
		"Outer", "Slot",
		"", "Terminal", // the slot is empty
		"Outer", "Terminal",
	); e != nil {
		t.Fatal(e)
	}

	filled := Outer{Slot: new(Pod)}
	if e := Check(t, walk.NewWalker(r.ValueOf(filled)),
		"Outer", "Initial",
		"Outer", "Slot",
		"", "Interface", // the pod element in the slot
		"Pod", "A",
		"Pod", "B",
		"Pod", "Terminal",
		"", "Terminal", // the end of the pod in the slot
		"Outer", "Terminal",
	); e != nil {
		t.Fatal(e)
	}
}

func Check(t *testing.T, it *walk.Walker, parts ...string) (err error) {
	var l listChecker
	for ok := true; ok; ok = it.Next() {
		c := it.Container().Type()
		if e := l.check(t, parts, c.Name()); e != nil && err == nil {
			// t.Error(e)
			err = e
		}
		var n string
		if v := it.Value(); !v.IsValid() {
			n = "Terminal"
		} else if it.Initial() {
			n = "Initial"
		} else if k := c.Kind(); k == r.Struct {
			n = it.Field().Name
		} else if k == r.Slice {
			i := it.Index()
			n = strconv.Itoa(i)
		} else if k == r.Interface {
			n = "Interface"
		}
		if e := l.check(t, parts, n); e != nil && err == nil {
			// t.Error(e)
			err = e
		}
	}
	if err == nil && int(l) < len(parts) {
		err = fmt.Errorf("ended early at line %d", l/2)
	}
	return
}

type listChecker int

// check the passed strings, one at a time.
func (lc *listChecker) check(t *testing.T, parts []string, got string) (err error) {
	t.Log(got)
	if i := int(*lc); i >= len(parts) {
		err = errors.New("out of range")
	} else if next := parts[i]; next != got {
		err = fmt.Errorf("mismatch at %d; expected %q got %q", i, next, got)
	} else {
		(*lc)++
	}
	return
}
