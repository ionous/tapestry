package core_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

func TestCompactEncoder(t *testing.T) {
	out := &assign.CallPattern{
		PatternName: core.P("factorial"),
		Arguments: []assign.Arg{{
			Name: "num",
			/// in the old bits this will generate an error if the assignment is emptyl
			Value: &assign.FromNumber{Value: core.I(1)},
		}}}

	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Factorial num:":{"FromNumber:":1}}` {
		t.Fatal(have)
	}
}

func TestGetTrait(t *testing.T) {
	out := &assign.ObjectRef{
		Name:  core.Variable("noun"),
		Field: core.Variable("trait"),
	}
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Object:field:":["@noun","@trait"]}` {
		t.Fatal(have)
	}
}

func TestGetVariablePath(t *testing.T) {
	out := core.Variable("pawn", "trait")

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Variable:dot:":["pawn",[{"AtField:":"trait"}]]}` {
		t.Fatal(have)
	}
}

// true, not a core test *shrug*
func TestListPush(t *testing.T) {
	out := &list.ListPush{
		Value:  &assign.FromText{Value: core.Variable("obj")},
		AtEdge: core.B(true),
		Target: core.Variable("bounds"),
	}

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Push:into:atFront:":[{"FromText:":"@obj"},"@bounds",true]}` {
		t.Fatal(have)
	}
}

// a unary function:
func TestUnaryEncoding(t *testing.T) {
	out := &core.Softline{}

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Wbr":true}` {
		t.Fatal(have)
	}
}
