package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func TestUpdateTracker(t *testing.T) {
	var up rules.UpdateTracker
	updates := &call.CallTrigger{}
	embedded := &logic.IsAll{Test: []rt.BoolEval{updates}}
	updatesArg := call.Arg{Value: &call.FromBool{Value: updates}}
	embeddedArg := call.Arg{Value: &call.FromBool{Value: embedded}}
	//
	negative := &text.IsEmpty{}
	negativeArg := call.Arg{Value: &call.FromBool{Value: negative}}
	//
	if b := up.CheckFilter(nil); b != false {
		t.Fatal("expected a blank check to return false")
	} else if b := up.HasUpdate(); b != false {
		t.Fatal("has update should be false")
	} else if b := up.CheckArgs(nil); b != false {
		t.Fatal("expected a blank check to return false")
	} else if b := up.CheckFilter(updates); b != true {
		t.Fatal("expected a a updates filter to return true")
	} else if b := up.CheckFilter(embedded); b != true {
		t.Fatal("expected an embedded filter to return true")
	} else if b := up.CheckArgs([]call.Arg{updatesArg}); b != true {
		t.Fatal("expected a a updates arg to return true")
	} else if b := up.CheckArgs([]call.Arg{embeddedArg}); b != true {
		t.Fatal("expected a a updates filter to return true")
	} else if b := up.CheckFilter(negative); b != false {
		t.Fatal("expected a a no counter filter to return false")
	} else if b := up.CheckArgs([]call.Arg{negativeArg}); b != false {
		t.Fatal("expected a a no counter arg to return false")
	} else if b := up.HasUpdate(); b != true {
		t.Fatal("has update should be true")
	}
}
