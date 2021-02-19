package qna

import (
	"bytes"
	"encoding/gob"
	"testing"

	"git.sr.ht/~ionous/iffy/assembly"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

// manually add an assembled pattern to the database, test that it works as expected.
func TestSayMe(t *testing.T) {
	gob.Register((*core.Text)(nil))
	gob.Register((*debug.MatchNumber)(nil))

	db := testdb.Open(t.Name(), testdb.Memory, assembly.SqlCustomDriver)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	}
	m := assembly.NewAssembler(db)
	if e := assembly.WritePattern(m, &debug.SayPattern); e != nil {
		t.Fatal(e)
	} else if e := m.WriteGob("say_me", &debug.SayPattern); e != nil {
		t.Fatal(e)
	}
	//
	if e := tables.CreateRun(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	}
	run := NewRuntime(db)
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
