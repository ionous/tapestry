package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// verify pattern declaration generate the simplest of pattern ephemera
func TestPatternImport(t *testing.T) {
	patternDecl := &story.DefinePattern{
		PatternName: T("corral"),
	}
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := patternDecl.PostImport(k); e != nil {
		t.Fatal(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphPatterns{
				PatternName: "corral",
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
	patternVariables := &story.DefinePattern{
		PatternName: T("corral"),
		Params: []story.FieldDefinition{&story.TextField{
			Name: "pet",
			Type: "animal",
		}},
	}
	var els []eph.Ephemera
	k := imp.NewImporter(collectEphemera(&els), storyMarshaller)
	if e := patternVariables.PostImport(k); e != nil {
		t.Log(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphPatterns{
				PatternName: "corral",
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
	k := imp.NewImporter(collectEphemera(&els), storyMarshaller)
	prog := story.DefinePattern{
		PatternName: T("example"),
		Rules: []story.PatternRule{{
			Guard: &core.Always{},
			Does: []rt.Execute{
				&core.PrintText{Text: T("hello")},
				&core.PrintText{Text: T("hello")},
			},
		}},
	}
	if e := prog.PostImport(k); e != nil {
		t.Fatal(e)
	} else {
		expect := []eph.Ephemera{
			&eph.EphPatterns{PatternName: "example"},
			&eph.EphRules{
				// the rules are for the named pattern.
				PatternName: "example",
				// "always" was specified as the guard.
				Filter: &core.Always{},
				// this is the default timing:
				When: eph.EphTiming{eph.EphTiming_During},
				// exe is exactly what was specified:
				Exe: []rt.Execute{
					&core.PrintText{Text: T("hello")},
					&core.PrintText{Text: T("hello")},
				},
			},
		}
		if diff := pretty.Diff(els, expect); len(diff) > 0 {
			t.Log(diff)
			t.Error(pretty.Sprint(els))
		}
	}
}
