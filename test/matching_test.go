package test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/testpat"
	"git.sr.ht/~ionous/iffy/test/testutil"
	test "git.sr.ht/~ionous/iffy/test/testutil"
)

func TestMatching(t *testing.T) {
	var kinds testutil.Kinds
	type Things struct{}
	kinds.AddKinds((*Things)(nil), (*GroupSettings)(nil), (*MatchGroups)(nil))
	k := kinds.Kind("group_settings")

	//
	lt := testpat.Runtime{
		pattern.Map{
			"match_groups": &matchGroups,
		},
		testutil.Runtime{
			Kinds: &kinds,
		},
	}

	a, b := k.NewRecord(), k.NewRecord()
	runMatching := &pattern.Determine{
		Pattern: "match_groups", Arguments: core.Args(
			&core.FromValue{g.RecordOf(a)},
			&core.FromValue{g.RecordOf(b)},
		)}
	// default should match
	{
		if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// different labels shouldnt match
	{
		if e := test.SetRecord(a, "label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != false {
			t.Fatal(e)
		}
	}
	// same labels should match
	{
		if e := test.SetRecord(b, "label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// many fields should match
	{
		if e := test.SetRecord(a, "innumerable", "is_innumerable"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "is_innumerable", true); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(a, "group_options", "objects_with_articles"); e != nil {
			t.Fatal(e)
		} else if e := test.SetRecord(b, "group_options", "objects_with_articles"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// names shouldnt be involved
	{
		if e := test.SetRecord(a, "name", "hola"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
}
