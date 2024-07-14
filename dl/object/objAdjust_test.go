package object_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/test"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestAdjustTraits(t *testing.T) {
	var kinds testutil.Kinds
	var objs testutil.Objects

	kinds.AddKinds(
		(*test.Messages)(nil),
	)
	objs.AddObjects(kinds.Kind("messages"), "msg")
	lt := testutil.Runtime{
		Kinds:   &kinds,
		Objects: objs,
	}

	if v, e := safe.GetText(&lt, &object.IncrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "scuffed" {
		t.Fatal(str)
	} else if v, e := safe.GetText(&lt, &object.IncrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
		Step:       literal.I(2),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "neat" {
		t.Fatal(str)
	} else if v, e := safe.GetText(&lt, &object.IncrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
		Step:       literal.I(5),
		Clamp:      literal.B(true),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "trampled" {
		t.Fatal(str)
	} else if v, e := safe.GetText(&lt, &object.DecrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "scuffed" {
		t.Fatal(str)
	} else if v, e := safe.GetText(&lt, &object.DecrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
		Step:       literal.I(2),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "trampled" {
		t.Fatal(str)
	} else if v, e := safe.GetText(&lt, &object.DecrementAspect{
		Target:     object.Object("msg"),
		AspectName: literal.T("neatness"),
		Step:       literal.I(5),
		Clamp:      literal.B(true),
	}); e != nil {
		t.Fatal(e)
	} else if str := v.String(); str != "neat" {
		t.Fatal(str)
	}
}
