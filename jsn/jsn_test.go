package jsn_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/debug"
	"github.com/kr/pretty"

	"git.sr.ht/~ionous/tapestry/jsn/dout"
)

func TestDetailsEncode(t *testing.T) {
	if d, e := dout.Encode(debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else if str := string(b); str != jsnTestIfx {
		t.Fatal(str)
	}
}

func TestDetailsDecode(t *testing.T) {
	var dst story.Story
	if e := din.Decode(&dst, tapestry.Registry(), []byte(jsnTestIfx)); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

func TestCompactEncoder(t *testing.T) {
	if str, e := cout.Marshal(debug.FactorialStory, story.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if str != jsnTestIf {
		t.Fatal(str)
	}
}

func TestCompactDecode(t *testing.T) {
	var dst story.Story
	if e := story.Decode(&dst, []byte(jsnTestIf), tapestry.AllSignatures); e != nil {
		pretty.Println(dst)
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

//go:embed jsnTest.ifx
var jsnTestIfx string

//go:embed jsnTest.if
var jsnTestIf string

// TestAnonymousSwap - unit test for broken parsing case
func TestAnonymousSwap(t *testing.T) {
	var jsnTestIf = `{"Listen kinds:handlers:": ["things",[]]}`
	want := story.EventBlock{
		Target: story.EventTarget{
			Value: &story.PluralKinds{
				Str: "things",
			},
			Choice: story.EventTarget_Kinds_Opt,
		},
		Handlers: make([]story.EventHandler, 0, 0),
	}
	//
	var have story.EventBlock
	if e := story.Decode(&have, []byte(jsnTestIf), tapestry.AllSignatures); e != nil {
		pretty.Println(have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println(have)
		t.Fatal(diff)
	}
}

// TestAnonymousOptional - unit test for broken parsing case
func TestAnonymousOptional(t *testing.T) {
	inputs := []string{
		`{ "NounRelation relation:otherNouns:":["whereabouts",[]]}`,
		`{ "NounRelation areBeing:relation:otherNouns:":["is", "whereabouts",[]]}`,
	}
	wants := []story.NounRelation{{
		AreBeing:   story.AreBeing{},
		Relation:   rel.RelationName{Str: "whereabouts"},
		OtherNouns: []story.NamedNoun{},
	}, {
		AreBeing:   story.AreBeing{Str: story.AreBeing_Is},
		Relation:   rel.RelationName{Str: "whereabouts"},
		OtherNouns: []story.NamedNoun{},
	}}
	for i, in := range inputs {
		var have story.NounRelation
		if e := story.Decode(&have, []byte(in), tapestry.AllSignatures); e != nil {
			pretty.Println("test", i, "got:", have)
			t.Fatal(e)
		} else if diff := pretty.Diff(&wants[i], &have); len(diff) != 0 {
			pretty.Println("test", i, "got:", have)
			t.Fatal(diff)
		}
	}
}

// TestVarAsBool - unit test for broken parsing case
// @requires light double committed
func TestVarAsBool(t *testing.T) {
	in := `{"AllTrue:":["@requires light",{"Get:from:":["is in darkness",{"VarFields:":"actor"}]}]}`
	want := core.AllTrue{Test: []rt.BoolEval{
		&core.GetVar{
			Name: core.VariableName{Str: "requires light"},
		},
		&core.GetAtField{
			Field: "is in darkness",
			From:  &core.FromVar{Var: core.VariableName{Str: "actor"}},
		},
	}}
	var have core.AllTrue
	if e := story.Decode(&have, []byte(in), tapestry.AllSignatures); e != nil {
		pretty.Println("got:", have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println("got:", have)
		t.Fatal(diff)
	}
}

func TestMissingSlot(t *testing.T) {
	in := `{"Join parts:":["one","two","three"]}`
	want := core.Join{Parts: []rt.TextEval{
		core.T("one"), core.T("two"), core.T("three"),
	}}
	var have core.Join
	if e := story.Decode(&have, []byte(in), tapestry.AllSignatures); e != nil {
		pretty.Println("got:", have)
		t.Fatal(e)
	} else if diff := pretty.Diff(&want, &have); len(diff) != 0 {
		pretty.Println("got:", have)
		t.Fatal(diff)
	}
}
