package qna

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"testing"

// 	"git.sr.ht/~ionous/iffy/assembly"
// 	"git.sr.ht/~ionous/iffy/dl/core"
// 	"git.sr.ht/~ionous/iffy/ephemera/debug"
// 	"git.sr.ht/~ionous/iffy/tables"
// 	"git.sr.ht/~ionous/iffy/test/testdb"
// )

// // how to test
// // - make sure that you can query multiple hooks for a noun or noun pairing
// // - make sure that you can build a list of hooks for a hierarchy
// // - maybe some raw query hooks
// // --- who gets the definition of an event hook?

// // - maybe instead of decoding gobs for now we just record number values
// func TestEventHooks(t *testing.T) {
// 	gob.Register((*core.Text)(nil))
// 	gob.Register((*debug.MatchNumber)(nil))

// 	db := newQnaDB(t, testdb.Memory)
// 	defer db.Close()
// 	if e := tables.CreateModel(db); e != nil {
// 		t.Fatal(e)
// 	}
// 	m := assembly.NewAssembler(db)
// 	if e := m.WriteGob("say_me", &debug.SayPattern); e != nil {
// 		t.Fatal(e)
// 	}
// 	//
// 	if e := tables.CreateRun(db); e != nil {
// 		t.Fatal(e)
// 	} else if e := tables.CreateRunViews(db); e != nil {
// 		t.Fatal(e)
// 	}
// 	run := NewRuntime(db)
// 	for i, expect := range []string{"One!", "Two!", "Three!", "Not between 1 and 3."} {
// 		var buf bytes.Buffer
// 		run.SetWriter(&buf)
// 		if e := debug.DetermineSay(i + 1).Execute(run); e != nil {
// 			t.Fatal(e)
// 		} else if text := buf.String(); expect != text {
// 			t.Fatal(i, text)
// 		} else {
// 			t.Log(text)
// 		}
// 	}
// }
