package core

import (
	"testing"

	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// test some simple functionality of the object commands using a mock runtime
func TestObjects(t *testing.T) {
	this, that, nothing := T("this"), T("that"), T("nothing")
	base, derived := T("base"), T("derived")

	run := modelTest{objClass: map[string]string{
		// objects:
		"this": base.Text.String(),
		"that": derived.Text.String(),
		// hierarchy:
		"base":    base.Text.String(),
		"derived": derived.Text.String() + "," + base.Text.String(),
	}}

	t.Run("exists", func(t *testing.T) {
		if e := testTrue(t, &run, &ObjectExists{Object: this}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &Not{Test: &ObjectExists{Object: nothing}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if cls, e := safe.GetText(&run, &KindOf{Object: this}); e != nil {
			t.Fatal(e)
		} else if cls.String() != base.Text.String() {
			t.Fatal("unexpected", cls)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &IsKindOf{Object: this, Kind: base.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsKindOf{Object: that, Kind: base.Text}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(t, &run, &IsKindOf{Object: that, Kind: derived.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &Not{&IsKindOf{this, derived.Text}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{A: &KindOf{this}, Is: &Equal{}, B: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{that}, Is: &Unequal{}, B: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{that}, Is: &Equal{}, B: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{this}, Is: &Unequal{}, B: derived}); e != nil {
			t.Fatal(e)
		}
	})
}

type modelTest struct {
	baseRuntime
	objClass map[string]string
}

func (m *modelTest) GetField(target, field string) (ret g.Value, err error) {
	if cls, ok := m.objClass[field]; !ok {
		err = g.UnknownField(target, field)
	} else {
		switch target {
		case object.Id:
			ret = g.StringOf(field)

		case object.Kind:
			ret = g.StringOf(cls)

		case object.Kinds:
			if path, ok := m.objClass[cls]; !ok {
				err = errutil.New("modelTest: unknown class", cls)
			} else {
				ret = g.StringOf(path)
			}

		default:
			err = g.UnknownField(target, field)
		}
	}
	return
}
