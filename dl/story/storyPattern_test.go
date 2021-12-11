package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"
)

func TestPatternVars(t *testing.T) {
	patternVariables := &story.PatternVariablesDecl{
		PatternName: core.PatternName{Str: "corral"},
		VariableDecl: []story.VariableDecl{{
			Type: story.VariableType{
				Choice: story.VariableType_Object_Opt,
				Value: &story.ObjectType{
					// An: story.Ana{
					// 	Str: "$AN",
					// },
					Kind: story.SingularKind{
						Str: "animal",
					},
				},
			},
			Name: N("pet"),
		}},
	}
	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)

	if e := patternVariables.ImportPhrase(k); e != nil {
		t.Log(e)
	} else {
		t.Fatal("beep")
		// var buf strings.Builder
		// tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		// tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		// if have, want := buf.String(), lines(
		// 	"corral,pattern",       // 1
		// 	"pet,parameter",        // 2
		// 	"animal,singular_kind", // 3
		// 	"2,3,4,0",              // NewPatternDecl
		// ); have != want {
		// 	t.Fatal("mismatch, have:", have)
		// } else {
		// 	t.Log("ok")
		// }
	}
}

func TestPatternDecl(t *testing.T) {
	patternDecl := &story.PatternDecl{
		Name: P("corral"),
		Type: story.PatternType{
			Str: story.PatternType_Patterns,
		},
	}

	var els []eph.Ephemera
	k := story.NewImporter(collectEphemera(&els), storyMarshaller)
	//
	if e := patternDecl.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		t.Fatal("beep")
		// var buf strings.Builder
		// tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		// tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		// if have, want := buf.String(), lines(
		// 	"corral,pattern", // 1
		// 	"patterns,type",  // 2
		// 	"2,2,3,0",        // NewPatternDecl
		// ); have != want {
		// 	t.Fatal("mismatch", have)
		// } else {
		// 	t.Log("ok")
		// }
	}
}
