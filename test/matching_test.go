package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
	"github.com/ionous/errutil"
)

func TestMatching(t *testing.T) {
	errutil.Panic = true
	var kinds testutil.Kinds
	type Things struct{}
	type Locals struct {
		A, B GroupSettings
	}
	kinds.AddKinds((*Things)(nil), (*GroupSettings)(nil), (*MatchGroups)(nil), (*Locals)(nil))
	kargs := kinds.Kind("match_groups")
	locals := kinds.NewRecord("locals")
	//
	lt := testpat.Runtime{
		testpat.Map{
			"match_groups": &matchGroups,
		},
		testutil.Runtime{
			Kinds: &kinds,
			Stack: []rt.Scope{
				g.RecordOf(locals),
			},
		},
	}

	if a, e := locals.GetNamedField("a"); e != nil {
		t.Fatal(e)
	} else if b, e := locals.GetNamedField("b"); e != nil {
		t.Fatal(e)
	} else {
		a, b := a.Record(), b.Record()

		runMatching := &core.CallPattern{
			Pattern: P(kargs.Name()), Arguments: core.MakeArgs(
				core.AssignFromRecord(core.GetVariable("a")),
				core.AssignFromRecord(core.GetVariable("b")),
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
}
