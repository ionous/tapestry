package core_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

func TestCompactEncoder(t *testing.T) {
	out := &core.CallPattern{
		Pattern: core.P("factorial"),
		Arguments: []core.Arg{{
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
	out := &core.GetValue{
		Source: core.Address{
			Choice: core.Address_Object_Opt,
			Value: &core.ObjectRef{
				Name:  core.GetVariable("noun"),
				Field: core.GetVariable("trait"),
			},
		}}

	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Get object:":{"Object:field:":["@noun","@trait"]}}` {
		t.Fatal(have)
	}
}

func TestGetVariablePath(t *testing.T) {
	out := &core.GetValue{
		Source: core.Address{
			Choice: core.Address_Variable_Opt,
			Value: &core.VariableRef{
				Name: core.T("pawn"),
				Dot: []core.Dot{
					&core.AtField{Field: core.T("trait")},
				},
			},
		}}

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Get variable:":{"Variable:dot:":["pawn",[{"AtField:":"trait"}]]}}` {
		t.Fatal(have)
	}
}

// true, not a core test *shrug*
func TestListPush(t *testing.T) {
	out := &list.ListPush{
		Value:  &assign.FromText{Value: core.GetVariable("obj")},
		AtEdge: core.B(true),
		Target: core.Address{
			Choice: core.Address_Variable_Opt,
			Value: &core.VariableRef{
				Name: core.T("bounds"),
			},
		}}

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Push:into variable:atFront:":[{"FromText:":"@obj"},"@bounds",true]}` {
		t.Fatal(have)
	}
}
