package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/rules"
	"github.com/kr/pretty"
)

func TestFindTree(t *testing.T) {
	if n := rules.FindTree(debugTree); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := rules.FindTree(impureTree); n != nil {
		t.Fatal("should 1have rejected the branching tree")
	}
}

func TestSplit(t *testing.T) {
	const ruleName = "some rule"
	if rs := rules.Split(ruleName, 0, impureTree); len(rs) != 1 {
		t.Fatal("expected one rule")
	} else if diff := pretty.Diff(rs[0], mdl.Rule{
		Name: ruleName,
		Prog: assign.Prog{
			Exe:        impureTree,
			Terminates: true,
		},
	}); len(diff) > 0 {
		t.Fatal("should have returned the passed block")
	} else if rs, updates := rules.SplitTree(ruleName, 0, terminalTree); len(rs) != 3 {
		t.Fatal("expected three rules")
	} else if updates != false {
		t.Fatal("expected no update")
	} else if diff := pretty.Diff(rs, []mdl.Rule{{
		Name: ruleName + " (3)",
		Prog: assign.Prog{
			Terminates: true,
		},
	}, {
		Name: ruleName + " (2)",
		Prog: assign.Prog{
			Filter: filter,
		},
	}, {
		Name: ruleName + " (1)",
		Prog: assign.Prog{
			Filter: filter,
		},
	}}); len(diff) > 0 {
		t.Fatal("should have returned the passed block")
	}
}

var filter = &literal.BoolValue{}

var tree = &core.ChooseBranch{
	If: filter,
	Else: &core.ChooseBranch{
		If:   &literal.BoolValue{},
		Else: &core.ChooseNothingElse{},
	},
}

var terminalTree = &core.ChooseBranch{
	If: filter,
	Else: &core.ChooseBranch{
		If:   filter,
		Else: &core.ChooseNothingElse{},
	},
}

var debugTree = []rt.Execute{
	&debug.DebugLog{},
	tree,
	&debug.DebugLog{},
}

var impureTree = append(debugTree, &core.PrintText{})
