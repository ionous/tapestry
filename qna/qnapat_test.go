package qna

import (
	"bytes"
	"testing"

	"git.sr.ht/~ionous/iffy/assembly"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
	"github.com/ionous/errutil"
)

// manually add an assembled pattern to the database, test that it works as expected.
func TestSayMe(t *testing.T) {
	errutil.Panic = true
	db := testdb.Open(t.Name(), testdb.Memory, assembly.SqlCustomDriver)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	}
	m := assembly.NewAssembler(db)
	// FIX: change to write sub pattern functions
	src := debug.SayPattern
	frag := assembly.PatternFrag{
		Name: src.Name, Return: src.Return, Labels: src.Labels, Fields: src.Fields,
	}
	if e := assembly.WriteRules(m, src.Name, "", "", src.Rules); e != nil {
		t.Fatal(e)
	} else if e := assembly.WriteFragment(m, "patterns", &frag); e != nil {
		t.Fatal(e)
	}
	//
	if e := tables.CreateRun(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	}
	run := NewRuntime(db, []map[uint64]interface{}{
		rt.Signatures,
		core.Signatures, {
			cin.Hash("SayMe:"):       (*debug.SayMe)(nil),
			cin.Hash("MatchNumber:"): (*debug.MatchNumber)(nil),
		}})
	for i, expect := range []string{"One!", "Two!", "Three!", "Not between 1 and 3."} {
		var buf bytes.Buffer
		run.SetWriter(&buf)
		if e := debug.DetermineSay(i + 1).Execute(run); e != nil {
			t.Fatal(e)
		} else if text := buf.String(); expect != text {
			t.Fatal(i+1, text)
		} else {
			t.Log(text)
		}
	}
}
