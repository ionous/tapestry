package jsn_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"

	"git.sr.ht/~ionous/iffy/jsn/dout"
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
	if e := din.Decode(&dst, iffy.Registry(), []byte(det)); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

func TestCompactEncoder(t *testing.T) {
	if d, e := cout.Encode(debug.FactorialStory); e != nil {
		t.Fatal(e)
	} else if b, e := json.Marshal(d); e != nil {
		t.Fatal(e)
	} else if str := string(b); str != com {
		t.Fatal(str)
	}
}

func TestCompactDecode(t *testing.T) {
	var dst story.Story
	if e := cin.Decode(&dst, []byte(com), iffy.AllSignatures); e != nil {
		pretty.Println(dst)
		t.Fatal(e)
	} else if diff := pretty.Diff(debug.FactorialStory, &dst); len(diff) != 0 {
		pretty.Print(dst)
		t.Fatal(diff)
	}
}

var det = `{"type":"story","value":{"$PARAGRAPH":[{"type":"paragraph","value":{"$STORY_STATEMENT":[{"type":"story_statement","value":{"type":"test_statement","value":{"$TEST":{"type":"testing","value":{"type":"test_output","value":{"$LINES":{"type":"lines","value":"6"}}}},"$TEST_NAME":{"type":"test_name","value":"factorial"}}}},{"type":"story_statement","value":{"type":"test_rule","value":{"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"say_text","value":{"$TEXT":{"type":"text_eval","value":{"type":"print_num","value":{"$NUM":{"type":"number_eval","value":{"type":"call_pattern","value":{"$ARGUMENTS":{"type":"call_args","value":{"$ARGS":[{"type":"call_arg","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":3}}}}}}},"$NAME":{"type":"text","value":"num"}}}]}},"$PATTERN":{"type":"pattern_name","value":"factorial"}}}}}}}}}}]}}}},"$TEST_NAME":{"type":"test_name","value":"factorial"}}}},{"type":"story_statement","value":{"type":"pattern_decl","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$OPTVARS":{"type":"pattern_variables_tail","value":{"$VARIABLE_DECL":[{"type":"variable_decl","value":{"$AN":{"type":"determiner","value":"$A"},"$NAME":{"type":"variable_name","value":"num"},"$TYPE":{"type":"variable_type","value":{"$PRIMITIVE":{"type":"primitive_type","value":"$NUMBER"}}}}}]}},"$TYPE":{"type":"pattern_type","value":"$PATTERNS"}}}},{"type":"story_statement","value":{"type":"pattern_actions","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$PATTERN_RETURN":{"type":"pattern_return","value":{"$RESULT":{"type":"variable_decl","value":{"$AN":{"type":"determiner","value":"$A"},"$NAME":{"type":"variable_name","value":"num"},"$TYPE":{"type":"variable_type","value":{"$PRIMITIVE":{"type":"primitive_type","value":"$NUMBER"}}}}}}},"$PATTERN_RULES":{"type":"pattern_rules","value":{"$PATTERN_RULE":[{"type":"pattern_rule","value":{"$GUARD":{"type":"bool_eval","value":{"type":"always","value":{}}},"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"assign","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"product_of","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"diff_of","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":1}}}}}}}}}}}}},"$VAR":{"type":"variable_name","value":"num"}}}}]}}}}}}]}}}}},{"type":"story_statement","value":{"type":"pattern_actions","value":{"$NAME":{"type":"pattern_name","value":"factorial"},"$PATTERN_RETURN":{"type":"pattern_return","value":{"$RESULT":{"type":"variable_decl","value":{"$AN":{"type":"determiner","value":"$A"},"$NAME":{"type":"variable_name","value":"num"},"$TYPE":{"type":"variable_type","value":{"$PRIMITIVE":{"type":"primitive_type","value":"$NUMBER"}}}}}}},"$PATTERN_RULES":{"type":"pattern_rules","value":{"$PATTERN_RULE":[{"type":"pattern_rule","value":{"$GUARD":{"type":"bool_eval","value":{"type":"compare_num","value":{"$A":{"type":"number_eval","value":{"type":"get_var","value":{"$NAME":{"type":"variable_name","value":"num"}}}},"$B":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":0}}}},"$IS":{"type":"comparator","value":{"type":"equal","value":{}}}}}},"$HOOK":{"type":"program_hook","value":{"$ACTIVITY":{"type":"activity","value":{"$EXE":[{"type":"execute","value":{"type":"assign","value":{"$FROM":{"type":"assignment","value":{"type":"from_num","value":{"$VAL":{"type":"number_eval","value":{"type":"num_value","value":{"$NUM":{"type":"number","value":1}}}}}}},"$VAR":{"type":"variable_name","value":"num"}}}}]}}}}}}]}}}}}]}}]}}`
var com = `{"Story:":[{"Paragraph:":[{"TestStatement:test:":["factorial",{"TestOutput:":"6"}]},{"TestRule:hook activity:":["factorial",{"Act:":[{"Say:":{"Numeral:":{"CallPattern:args:":["factorial",{"Inargs:":[{"Inarg:from:":["num",{"FromNum:":3}]}]}]}}}]}]},{"PatternDecl:name:optvars:":["patterns","factorial",{"PatternVariablesTail:":[{"VariableDecl:name:type primitive:":["a","num","number"]}]}]},{"PatternActions:patternReturn:patternRules:":["factorial",{"PatternReturn:":{"VariableDecl:name:type primitive:":["a","num","number"]}},{"PatternRules:":[{"PatternRule:hook activity:":["Always",{"Act:":[{"Let:be:":["num",{"FromNum:":{"Mul:by:":["@num",{"Dec:by:":["@num",1]}]}}]}]}]}]}]},{"PatternActions:patternReturn:patternRules:":["factorial",{"PatternReturn:":{"VariableDecl:name:type primitive:":["a","num","number"]}},{"PatternRules:":[{"PatternRule:hook activity:":[{"Cmp:is:num:":["@num","Equals",0]},{"Act:":[{"Let:be:":["num",{"FromNum:":1}]}]}]}]}]}]}]}`
