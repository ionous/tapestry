package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// verify pattern declaration generate the simplest of pattern ephemera
func TestPatternImport(t *testing.T) {
	patternDecl := &story.PatternDecl{
		Name: P("corral"),
		Type: story.PatternType{
			Str: story.PatternType_Patterns,
		},
	}
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := patternDecl.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphPatterns{
				Name: "corral",
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Log(diff)
			t.Error(pretty.Sprint(els))
		}
	}
}

// verify pattern parameter declarations generate pattern parameter ephemera
func TestPatternParameterImport(t *testing.T) {
	patternVariables := &story.PatternVariablesDecl{
		PatternName: core.PatternName{Str: "corral"},
		Props: []story.PropertySlot{&story.TextProperty{story.NamedProperty{
			Name: "pet",
			Type: "animal",
		}}},
	}
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := patternVariables.ImportPhrase(k); e != nil {
		t.Log(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphPatterns{
				Name: "corral",
				Params: []eph.EphParams{{
					Affinity: eph.Affinity{eph.Affinity_Text},
					Name:     "pet",
					Class:    "animal",
				}},
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Log(diff)
			t.Error(pretty.Sprint(els))
		}
	}
}

// verify that pattern rules generate appropriate ephemera
// see also: TestFactorialImport which is more extensive
func TestPatternRuleImport(t *testing.T) {
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)

	var prog story.PatternActions
	if b, e := json.Marshal(_pattern_actions); e != nil {
		t.Fatal(e)
	} else if e := din.Decode(&prog, tapestry.Registry(), b); e != nil {
		t.Fatal(e)
	} else if e := prog.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		expect := []eph.Ephemera{
			// one pattern, no parameters, no locals, no return value.
			&eph.EphPatterns{
				Name: "example",
			},
			&eph.EphRules{
				// the rules are for the named pattern.
				Name: "example",
				// "always" was specified as the guard.
				Filter: &core.Always{},
				// this is the default timing:
				When: eph.EphTiming{eph.EphTiming_During},
				// exe is exactly what was specified:
				Exe: &core.Activity{[]rt.Execute{
					&core.SayText{T("hello")},
					&core.SayText{T("hello")},
				}},
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Log(diff)
			t.Error(pretty.Sprint(els))
		}
	}
}

// see also: TestStoryActivity which runs the hook specified by these rules
var _pattern_actions = map[string]interface{}{
	"type": "pattern_actions",
	"value": map[string]interface{}{
		"$NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "example",
		},
		"$PATTERN_RULES": map[string]interface{}{
			"type": "pattern_rules",
			"value": map[string]interface{}{
				"$PATTERN_RULE": []interface{}{
					map[string]interface{}{
						"type": "pattern_rule",
						"value": map[string]interface{}{
							"$GUARD": map[string]interface{}{
								"type": "bool_eval",
								"value": map[string]interface{}{
									"type":  "always",
									"value": map[string]interface{}{},
								}},
							"$HOOK": map[string]interface{}{
								"type": "program_hook",
								"value": map[string]interface{}{
									"$ACTIVITY": _pattern_activity,
								}}}}}}}},
}
