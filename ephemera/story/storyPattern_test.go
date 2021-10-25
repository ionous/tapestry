package story_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func TestPatternVars(t *testing.T) {
	patternVariables := &story.PatternVariablesDecl{
		PatternName: value.PatternName{Str: "corral"},
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
	k, db := newImporter(t, testdb.Memory)
	defer db.Close()

	if e := patternVariables.ImportPhrase(k); e != nil {
		t.Log(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern",       // 1
			"pet,parameter",        // 2
			"animal,singular_kind", // 3
			"2,3,4,0",              // NewPatternDecl
		); have != want {
			t.Fatal("mismatch, have:", have)
		} else {
			t.Log("ok")
		}
	}
}

func TestPatternDecl(t *testing.T) {
	patternDecl := &story.PatternDecl{
		Name: P("corral"),
		Type: story.PatternType{
			Str: story.PatternType_Patterns,
		},
	}

	k, db := newImporter(t, testdb.Memory)
	defer db.Close()
	if e := patternDecl.ImportPhrase(k); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
		tables.WriteCsv(db, &buf, "select idNamedPattern,idNamedParam,idNamedType,idProg from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern", // 1
			"patterns,type",  // 2
			"2,2,3,0",        // NewPatternDecl
		); have != want {
			t.Fatal("mismatch", have)
		} else {
			t.Log("ok")
		}
	}
}
