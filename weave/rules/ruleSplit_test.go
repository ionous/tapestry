package rules_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/rt"
)

func TestFindTree(t *testing.T) {
	if n := logic.PickTree(branchingBlock); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := logic.PickTree(blockWithDebugLogs); n != tree {
		t.Fatal("expected to find the branching tree")
	} else if n := logic.PickTree(blockWithTail); n != nil {
		t.Fatal("should have rejected the branching tree")
	} else if n := logic.PickTree(blockWithSiblingBranches); n != nil {
		t.Fatal("should have rejected the sibling branches")
	}
}

var filter = &literal.BoolValue{}

var tree = &logic.ChooseBranch{
	Condition: filter,
	Else: &logic.ChooseBranch{
		Condition: &literal.BoolValue{},
		Else:      &logic.ChooseNothingElse{},
	},
}

// a block that consists only of branching statements
var branchingBlock = []rt.Execute{
	tree,
}

// a block with some extraneous debug logs
var blockWithDebugLogs = []rt.Execute{
	&debug.LogValue{},
	tree,
	&debug.LogValue{},
}

// the branching block but has some non-branching trailing statements
// ( so its not really a pure set of rules )
var blockWithTail = append(blockWithDebugLogs, &format.PrintText{})

// two sibling blocks -- these should be considered as non-branching.
var blockWithSiblingBranches = []rt.Execute{
	tree,
	tree,
}
