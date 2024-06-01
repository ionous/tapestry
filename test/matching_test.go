package test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/scope"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/object"
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
	kargs := kinds.Kind("match groups")
	locals := kinds.NewRecord("locals")
	//
	lt := testpat.Runtime{
		Map: testpat.Map{
			"match groups": &matchGroups,
		},
		Runtime: testutil.Runtime{
			Kinds: &kinds,
		},
	}
	lt.Chain = scope.MakeChain(scope.FromRecord(&kinds, locals))

	if a, e := lt.FieldByName("a"); e != nil {
		t.Fatal(e)
	} else if b, e := lt.FieldByName("b"); e != nil {
		t.Fatal(e)
	} else {
		a, b := a.Record(), b.Record()
		runMatching := &call.CallPattern{
			PatternName: (kargs.Name()), Arguments: core.MakeArgs(
				&assign.FromRecord{Value: object.Variable("a")},
				&assign.FromRecord{Value: object.Variable("b")},
			)}
		// default should match
		{
			if ok, e := runMatching.GetBool(&lt); e != nil {
				t.Fatal(e)
			} else if res := ok.Bool(); res != true {
				t.Fatal("matched", res)
			}
		}
		// different labels shouldnt match
		{
			if e := testutil.SetRecord(a, "label", "beep"); e != nil {
				t.Fatal(e)
			} else if ok, e := runMatching.GetBool(&lt); e != nil {
				t.Fatal(e)
			} else if res := ok.Bool(); res != false {
				t.Fatal("different labels shouldn't match")
			}
		}
		// same labels should match
		{
			if e := testutil.SetRecord(b, "label", "beep"); e != nil {
				t.Fatal(e)
			} else if ok, e := runMatching.GetBool(&lt); e != nil {
				t.Fatal(e)
			} else if res := ok.Bool(); res != true {
				t.Fatal("same labels should match")
			}
		}
		// many fields should match
		{
			if e := testutil.SetRecord(a, "group options", "objects with articles"); e != nil {
				t.Fatal(e)
			} else if e := testutil.SetRecord(b, "group options", "objects with articles"); e != nil {
				t.Fatal(e)
			} else if ok, e := runMatching.GetBool(&lt); e != nil {
				t.Fatal(e)
			} else if res := ok.Bool(); res != true {
				t.Fatal("many fields should match")
			}
		}
		// names shouldnt be involved
		{
			if e := testutil.SetRecord(a, "name", "hola"); e != nil {
				t.Fatal(e)
			} else if ok, e := runMatching.GetBool(&lt); e != nil {
				t.Fatal(e)
			} else if res := ok.Bool(); res != true {
				t.Fatal("names shouldnt matter")
			}
		}
	}
}
