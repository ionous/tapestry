package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestMatching(t *testing.T) {
	var kinds testutil.Kinds
	type Things struct{}
	kinds.AddKinds((*Things)(nil), (*GroupSettings)(nil), (*MatchGroups)(nil))
	k := kinds.Kind("group_settings")

	//
	lt := testpat.Runtime{
		testpat.Map{
			"match_groups": &matchGroups,
		},
		testutil.Runtime{
			Kinds: &kinds,
		},
	}

	a, b := k.NewRecord(), k.NewRecord()
	runMatching := &core.CallPattern{
		Pattern: P("match_groups"), Arguments: core.Args(
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
		if e := testutil.SetRecord(a, "label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != false {
			t.Fatal(e)
		}
	}
	// same labels should match
	{
		if e := testutil.SetRecord(b, "label", "beep"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// many fields should match
	{
		if e := testutil.SetRecord(a, "innumerable", "is_innumerable"); e != nil {
			t.Fatal(e)
		} else if e := testutil.SetRecord(b, "is_innumerable", true); e != nil {
			t.Fatal(e)
		} else if e := testutil.SetRecord(a, "group_options", "objects_with_articles"); e != nil {
			t.Fatal(e)
		} else if e := testutil.SetRecord(b, "group_options", "objects_with_articles"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
	// names shouldnt be involved
	{
		if e := testutil.SetRecord(a, "name", "hola"); e != nil {
			t.Fatal(e)
		} else if ok, e := runMatching.GetBool(&lt); e != nil {
			t.Fatal(e)
		} else if ok.Bool() != true {
			t.Fatal(e)
		}
	}
}
