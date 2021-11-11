package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"

	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func TestImportStory(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else {
		var curr story.Story
		if e := din.Decode(&curr, iffy.Registry(), []byte(debug.Blob)); e != nil {
			t.Fatal(e)
		} else {
			k := story.NewImporter(db, cout.Marshal)
			if e := k.ImportStory(t.Name(), &curr); e != nil {
				t.Fatal("import", e)
			} else {
				t.Log("ok")
			}
		}
	}
}

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func F(n float64) rt.NumberEval     { return &core.NumValue{n} }
func I(n int) rt.NumberEval         { return &core.NumValue{float64(n)} }
func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func T(s string) *core.TextValue    { return &core.TextValue{W(s)} }
func V(i string) *core.GetVar       { return &core.GetVar{N(i)} }
func W(v string) string             { return v }
