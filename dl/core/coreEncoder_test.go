package core_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/encode"
)

func TestGetTrait(t *testing.T) {
	out := &assign.ObjectRef{
		Name:  core.Variable("noun"),
		Field: core.Variable("trait"),
	}
	if have, e := stringify(out); e != nil {
		t.Fatal(e)
	} else if have != `{"Object:field:":["@noun","@trait"]}` {
		t.Fatal(have)
	}
}

func TestGetVariablePath(t *testing.T) {
	out := core.Variable("pawn", "trait")

	// fix: this should have path syntax ( re: expressions )
	// ex. @pawn.trait
	if have, e := stringify(out); e != nil {
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
	if have, e := stringify(out); e != nil {
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
	if have, e := stringify(out); e != nil {
		t.Fatal(e)
	} else if have != `{"Wbr":true}` {
		t.Fatal(have)
	}
}

func stringify(v jsn.Marshalee) (ret string, err error) {
	var enc encode.Encoder
	if res, e := enc.CustomEncoder(core.CustomEncoder).MarshalFlow(v); e != nil {
		err = e
	} else if b, e := json.Marshal(res); e != nil {
		err = e
	} else {
		ret = string(b)
	}
	return
}
