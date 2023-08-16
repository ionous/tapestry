package story_test

import (
	"testing"
)

// verify pattern declaration generate the simplest of pattern ephemera
func xTestPatternImport(t *testing.T) {
	// 	patternDecl := &story.DefinePattern{
	// 		PatternName: T("corral"),
	// 	}
	// 	var els []eph.Ephemera
	// 	k := weave.NewCatalog(eph.NewCommandQueue(&els))
	// 	if e := patternDecl.Schedule(k); e != nil {
	// 		t.Fatal(e)
	// 	} else {
	// 		expect := []eph.Ephemera{
	// 			&eph.Patterns{
	// 				PatternName: "corral",
	// 			},
	// 		}
	// 		if diff := pretty.Diff(els, expect); len(diff) > 0 {
	// 			t.Log(diff)
	// 			t.Error(pretty.Sprint(els))
	// 		}
	// 	}
}

// verify pattern parameter declarations generate pattern parameter ephemera
func xTestPatternParameterImport(t *testing.T) {
	// 	patternVariables := &story.DefinePattern{
	// 		PatternName: T("corral"),
	// 		Params: []story.FieldDefinition{&story.TextField{
	// 			Name: "pet",
	// 			Type: "animal",
	// 		}},
	// 	}
	// 	var els []eph.Ephemera
	// 	k := weave.NewCatalog(eph.NewCommandQueue(&els))

	// 	if e := patternVariables.Schedule(k); e != nil {
	// 		t.Log(e)
	// 	} else {
	// 		expect := []eph.Ephemera{
	// 			&eph.Patterns{
	// 				PatternName: "corral",
	// 			},
	// 			&eph.Patterns{
	// 				PatternName: "corral",
	// 				Params: []eph.Params{{
	// 					Affinity: eph.affine.Text},
	// 					Name:     "pet",
	// 					Class:    "animal",
	// 				}},
	// 			},
	// 		}
	// 		if diff := pretty.Diff(els, expect); len(diff) > 0 {
	// 			t.Log(diff)
	// 			t.Error("got:", pretty.Sprint(els))
	// 		}
	// 	}
}

// verify that pattern rules generate appropriate ephemera
// see also: TestFactorialImport which is more extensive
func xTestPatternRuleImport(t *testing.T) {
	// var els []eph.Ephemera
	// k := weave.NewCatalog(eph.NewCommandQueue(&els))
	// prog := story.DefinePattern{
	// 	PatternName: T("example"),
	// 	Rules: []story.PatternRule{{
	// 		Guard: &core.Always{},
	// 		Exe: []rt.Execute{
	// 			&core.PrintText{Text: T("hello")},
	// 			&core.PrintText{Text: T("hello")},
	// 		},
	// 	}},
	// }
	// if e := prog.Schedule(k); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	expect := []eph.Ephemera{
	// 		&eph.Patterns{PatternName: "example"},
	// 		&eph.Rules{
	// 			// the rules are for the named pattern.
	// 			PatternName: "example",
	// 			// "always" was specified as the guard.
	// 			Filter: &core.Always{},
	// 			// this is the default timing:
	// 			When: eph.Timing{eph.Timing_During},
	// 			// exe is exactly what was specified:
	// 			Exe: []rt.Execute{
	// 				&core.PrintText{Text: T("hello")},
	// 				&core.PrintText{Text: T("hello")},
	// 			},
	// 		},
	// 	}
	// 	if diff := pretty.Diff(els, expect); len(diff) > 0 {
	// 		t.Log(diff)
	// 		t.Error("got:", pretty.Sprint(els))
	// 	}
	// }
}
