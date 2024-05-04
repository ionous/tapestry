package rt_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

// access a record via the value interface
func TestRecordValue(t *testing.T) {
	k := rt.NewKind([]string{"Ks"}, []rt.Field{
		{Name: "d", Affinity: affine.Number},
		{Name: "t", Affinity: affine.Text},
		{Name: "a", Affinity: affine.Text, Type: "a"},
	}, []rt.Aspect{{
		Name:   "a",
		Traits: []string{"x", "w", "y"},
	}})

	t.Run("numbers", func(t *testing.T) {
		rv := rt.RecordOf(rt.NewRecord(k))
		if el, e := rv.FieldByName("d"); e != nil {
			t.Fatal(e)
		} else if v := el.Float(); v != 0 {
			t.Fatal("not default", v)
		} else if e := rv.SetFieldByName("d", rt.FloatOf(5)); e != nil {
			t.Fatal(e)
		} else if el, e := rv.FieldByName("d"); e != nil {
			t.Fatal(e)
		} else if v := el.Float(); v != 5 {
			t.Fatal("not changed", v)
		}
	})
	t.Run("text", func(t *testing.T) {
		rv := rt.RecordOf(rt.NewRecord(k))
		if el, e := rv.FieldByName("t"); e != nil {
			t.Fatal(e)
		} else if v := el.String(); len(v) > 0 {
			t.Fatal("not default", v)
		} else if e := rv.SetFieldByName("t", rt.StringOf("xyzzy")); e != nil {
			t.Fatal(e)
		} else if el, e := rv.FieldByName("t"); e != nil {
			t.Fatal(e)
		} else if v := el.String(); v != "xyzzy" {
			t.Fatal("not changed", v)
		}
	})
	t.Run("aspects", func(t *testing.T) {
		rv := rt.RecordOf(rt.NewRecord(k))
		if el, e := rv.FieldByName("x"); e != nil {
			t.Fatal(e)
		} else if v := el.Bool(); !v {
			t.Fatal("not default", v)
		} else if e := rv.SetFieldByName("x", rt.BoolOf(true)); e != nil {
			t.Fatal(e)
		} else if el, e := rv.FieldByName("x"); e != nil {
			t.Fatal(e)
		} else if v := el.Bool(); v != true {
			t.Fatal("not changed", v)
		} else if el, e := rv.FieldByName("a"); e != nil {
			t.Fatal(e)
		} else if v := el.String(); v != "x" {
			t.Fatal(e)
		} else if e := rv.SetFieldByName("a", rt.StringOf("w")); e != nil {
			t.Fatal(e)
		} else if el, e := rv.FieldByName("w"); e != nil {
			t.Fatal(e)
		} else if v := el.Bool(); v != true {
			t.Fatal("aspect not changed")
		}
	})
	t.Run("failures", func(t *testing.T) {
		rv := rt.RecordOf(rt.NewRecord(k))
		if _, e := rv.FieldByName("nope"); e == nil {
			t.Fatal("expected no such field")
		} else if e := rv.SetFieldByName("a", rt.True); e == nil {
			t.Fatal("aspects should be set with strings")
		} else if e := rv.SetFieldByName("x", rt.Empty); e == nil {
			t.Fatal("traits should be set with bools")
		} else if e := rv.SetFieldByName("x", rt.False); e == nil {
			// we dont have support for opposite values right now.
			t.Fatal("traits should be set with true values only")
		}
	})
}
