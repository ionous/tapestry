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
	this, that, nothing := named("this"), named("that"), named("nothing")
	base, derived := &Text{"base"}, &Text{"derived"}

	run := modelTest{objClass: map[string]string{
		// objects:
		"this": base.Text,
		"that": derived.Text,
		// hierarchy:
		"base":    base.Text,
		"derived": derived.Text + "," + base.Text,
	}}

	t.Run("exists", func(t *testing.T) {
		if e := testTrue(t, &run, &ObjectExists{this}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&ObjectExists{nothing}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if cls, e := safe.GetText(&run, &KindOf{this}); e != nil {
			t.Fatal(e)
		} else if cls.String() != base.Text {
			t.Fatal("unexpected", cls)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &IsKindOf{this, base.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsKindOf{that, base.Text}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(t, &run, &IsKindOf{that, derived.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&IsKindOf{this, derived.Text}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{&KindOf{this}, &EqualTo{}, base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{that}, &NotEqualTo{}, base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{that}, &EqualTo{}, derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{this}, &NotEqualTo{}, derived}); e != nil {
			t.Fatal(e)
		}
	})
}

func named(n string) *Text {
	return &Text{n}
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
