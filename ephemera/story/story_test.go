package story_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func TestImportStory(t *testing.T) {
	db := testdb.Open(t.Name(), testdb.Memory, "")
	defer db.Close()
	//
	var in reader.Map
	if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
		t.Fatal("read json", e)
	} else if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create tables", e)
	} else {
		k := story.NewImporter(db, func(pos reader.Position, err error) {
			t.Errorf("%s at %s", err, pos)
		})
		if _, e := k.ImportStory(t.Name(), in); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}

func B(b bool) rt.BoolEval          { return &core.BoolValue{b} }
func F(n float64) rt.NumberEval     { return &core.NumValue{n} }
func I(n int) rt.NumberEval         { return &core.NumValue{float64(n)} }
func P(p string) value.PatternName  { return value.PatternName{Str: p} }
func N(v string) value.VariableName { return value.VariableName{Str: v} }
func T(s string) *core.TextValue    { return &core.TextValue{W(s)} }
func V(i string) *core.Var          { return &core.Var{N(i)} }
func W(v string) value.Text         { return value.Text{Str: v} }
