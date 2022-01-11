package jsn_test

import (
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

	"git.sr.ht/~ionous/tapestry/jsn/dout"
	"github.com/kr/pretty"
)

func TestDetailsEncode(t *testing.T) {
	if d, e := dout.Encode(debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else if str := string(b); str != det {
		t.Fatal(str)
	}
}

func TestDetailsDecode(t *testing.T) {
	var dst story.Story
	if e := din.Decode(&dst, tapestry.Registry(), []byte(det)); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

func TestCompactEncoder(t *testing.T) {
	if str, e := cout.Marshal(debug.FactorialStory, story.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if str != com {
		t.Fatal(str)
	}
}

func TestCompactDecode(t *testing.T) {
	var dst story.Story
	if e := story.Decode(&dst, []byte(com), tapestry.AllSignatures); e != nil {
		pretty.Println(dst)
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

var det = `{"type":"story","value":{"$PARAGRAPH":[{"type":"paragraph","value":{"$STORY_STATEMENT":[{"type":"story_statement","value":{"type":"test_statement","value":{"$TEST":{"type":"testing","value":{"type":"test_output","value":{"$LINES":{"type":"lines","value":"6"}}}},"$TEST_NAME":{"type":"test_name","value":"factorial"}}}},{"type":"story_statement","value":{"type":"test_rule","value":{"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"say_text","value":{"$TEXT":{"type":"text_eval","value":{"type":"print_num","value":{"$NUM":{"type":"number_eval","value":{"type":"call_pattern","value":{"$ARGUMENTS":{"type":"call_args","value":{"$ARGS":[{"type":"call_arg","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":3}}}}}}},"$NAME":{"type":"text","value":"num"}}}]}},"$PATTERN":{"type":"pattern_name","value":"factorial"}}}}}}}}}}]}}}},"$TEST_NAME":{"type":"test_name","value":"factorial"}}}},{"type":"story_statement","value":{"type":"pattern_decl","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$OPTVARS":{"type":"pattern_variables_tail","value":{"$PROPS":[{"type":"property_slot","value":{"type":"number_property","value":{"$NAME":{"type":"text","value":"num"}}}}]}},"$TYPE":{"type":"pattern_type","value":"$PATTERNS"}}}},{"type":"story_statement","value":{"type":"pattern_actions","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$PATTERN_RETURN":{"type":"pattern_return","value":{"$RESULT":{"type":"property_slot","value":{"type":"number_property","value":{"$NAME":{"type":"text","value":"num"}}}}}},"$PATTERN_RULES":{"type":"pattern_rules","value":{"$PATTERN_RULE":[{"type":"pattern_rule","value":{"$GUARD":{"type":"bool_eval","value":{"type":"always","value":{}}},"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"assign","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"product_of","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"diff_of","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":1}}}}}}}}}}}}},"$VAR":{"type":"variable_name","value":"num"}}}}]}}}}}}]}}}}},{"type":"story_statement","value":{"type":"pattern_actions","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$PATTERN_RETURN":{"type":"pattern_return","value":{"$RESULT":{"type":"property_slot","value":{"type":"number_property","value":{"$NAME":{"type":"text","value":"num"}}}}}},"$PATTERN_RULES":{"type":"pattern_rules","value":{"$PATTERN_RULE":[{"type":"pattern_rule","value":{"$GUARD":{"type":"bool_eval","value":{"type":"compare_num","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":0}}}},"$IS":{"type":"comparator","value":{"type":"equal","value":{}}}}}},"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"assign","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":1}}}}}}},"$VAR":{"type":"variable_name","value":"num"}}}}]}}}}}}]}}}}}]}}]}}`
var com = `{"Story:":[{"Paragraph:":[{"TestStatement:test:":["factorial",{"TestOutput:":"6"}]},{"TestRule:hook activity:":["factorial",{"Act:":[{"Say:":{"Numeral:":{"CallPattern:args:":["factorial",{"Inargs:":[{"Inarg:from:":["num",{"FromNum:":3}]}]}]}}}]}]},{"PatternDecl:name:optvars:":["patterns","factorial",{"PatternVariablesTail:":[{"Number named:":"num"}]}]},{"PatternActions:patternReturn:patternRules:":["factorial",{"PatternReturn:":{"Number named:":"num"}},{"PatternRules:":[{"PatternRule:hook activity:":["Always",{"Act:":[{"Let:be:":["num",{"FromNum:":{"Mul:by:":["@num",{"Dec:by:":["@num",1]}]}}]}]}]}]}]},{"PatternActions:patternReturn:patternRules:":["factorial",{"PatternReturn:":{"Number named:":"num"}},{"PatternRules:":[{"PatternRule:hook activity:":[{"Cmp:is:num:":["@num","Equals",0]},{"Act:":[{"Let:be:":["num",{"FromNum:":1}]}]}]}]}]}]}]}`

// TestAnonymousSwap - unit test for broken parsing case
func TestAnonymousSwap(t *testing.T) {
	var com = `{"EventBlock kinds:handlers:": ["things",[]]}`
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
	if e := story.Decode(&have, []byte(com), tapestry.AllSignatures); e != nil {
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
	want := core.AllTrue{[]rt.BoolEval{
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
