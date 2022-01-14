package core

import (
	"strings"
	"testing"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// test some simple functionality of the object commands using a mock runtime
func TestObjects(t *testing.T) {
	this, that, nothing := T("this"), T("that"), T("nothing")
	base, derived := T("base"), T("derived")

	run := modelTest{objClass: map[string]string{
		// objects:
		"this": base.Text,
		"that": derived.Text,
		// hierarchy:
		"base":    base.Text,
		"derived": base.Text + "," + derived.Text,
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
		} else if cls.String() != base.Text {
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
		if e := testTrue(t, &run, &Not{Test: &IsKindOf{Object: this, Kind: derived.Text}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{A: &KindOf{Object: this}, Is: &Equal{}, B: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{Object: that}, Is: &Unequal{}, B: base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{Object: that}, Is: &Equal{}, B: derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{A: &KindOf{Object: this}, Is: &Unequal{}, B: derived}); e != nil {
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
		case meta.ObjectId:
			ret = g.StringOf(field)

		case meta.ObjectKind:
			ret = g.StringOf(cls)

		case meta.ObjectKinds:
			if path, ok := m.objClass[cls]; !ok {
				err = errutil.New("modelTest: unknown class", cls)
			} else {
				ret = g.StringsOf(strings.Split(path, ","))
			}

		default:
			err = g.UnknownField(target, field)
		}
	}
	return
}
