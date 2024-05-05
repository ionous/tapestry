package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

func TestFindTree(t *testing.T) {
	if n := core.PickTree(branchingBlock); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := core.PickTree(blockWithDebugLogs); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := core.PickTree(blockWithTail); n != nil {
		t.Fatal("should have rejected the branching tree")
	} else if n := core.PickTree(blockWithSiblingBranches); n != nil {
		t.Fatal("should have rejected the sibling branches")
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

// a block that consists only of branching statements
var branchingBlock = []rt.Execute{
	tree,
}

// a block with some extraneous debug logs
var blockWithDebugLogs = []rt.Execute{
	&debug.DebugLog{},
	tree,
	&debug.DebugLog{},
}

// the branching block but has some non-branching trailing statements
// ( so its not really a pure set of rules )
var blockWithTail = append(blockWithDebugLogs, &core.PrintText{})

// two sibling blocks -- these should be considered as non-branching.
var blockWithSiblingBranches = []rt.Execute{
	tree,
	tree,
}
