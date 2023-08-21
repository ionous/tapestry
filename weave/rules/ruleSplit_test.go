package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

func TestFindTree(t *testing.T) {
	if n := core.FindBranch(debugTree); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := core.FindBranch(impureTree); n != nil {
		t.Fatal("should 1have rejected the branching tree")
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
