package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func TestUpdateTracker(t *testing.T) {
	var up rules.UpdateTracker
	direct := &core.CallTrigger{}
	embedded := &core.AllTrue{Test: []rt.BoolEval{direct}}
	directArg := assign.Arg{Value: &assign.FromBool{Value: direct}}
	embeddedArg := assign.Arg{Value: &assign.FromBool{Value: embedded}}
	negative := &core.IsEmpty{}
	negativeArg := assign.Arg{Value: &assign.FromBool{Value: negative}}
	if b := up.CheckFilter(nil); b != false {
		t.Fatal("expected a blank check to return false")
	} else if b := up.HasUpdate(); b != false {
		t.Fatal("has update should be false")
	} else if b := up.CheckArgs(nil); b != false {
		t.Fatal("expected a blank check to return false")
	} else if b := up.CheckFilter(direct); b != true {
		t.Fatal("expected a a direct filter to return true")
	} else if b := up.CheckFilter(embedded); b != true {
		t.Fatal("expected an embedded filter to return true")
	} else if b := up.CheckArgs([]assign.Arg{directArg}); b != true {
		t.Fatal("expected a a direct arg to return true")
	} else if b := up.CheckArgs([]assign.Arg{embeddedArg}); b != true {
		t.Fatal("expected a a direct filter to return true")
	} else if b := up.CheckFilter(negative); b != false {
		t.Fatal("expected a a no counter filter to return false")
	} else if b := up.CheckArgs([]assign.Arg{negativeArg}); b != false {
		t.Fatal("expected a a no counter arg to return false")
	} else if b := up.HasUpdate(); b != true {
		t.Fatal("has update should be true")
	}
}
