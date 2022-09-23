package qna

// FIX: have to remap using the asm table writer...
// no idea where this test should live...
// tests the importation, assembly, and execution of a factorial story.
// doesn't test the *reading* of the story.
// func TestFullFactorial(t *testing.T) {
// 	db := testdb.Open(t.Name(), testdb.Memory, "")
// 	defer db.Close()

// 	// read factorialStory, assemble and run.
// 	var ds reader.Dilemmas
// 	if e := tables.CreateAll(db); e != nil {
// 		t.Fatal("couldn't create tables", e)
// 	} else {
// 		k := imp.NewImporter(dbwriter(db), storyMarshaller)
// 		if e := asm.ImportStory(k, t.Name(), debug.FactorialStory); e != nil {
// 			t.Fatal("couldn't import story", e)
// 		} else if e := asm.AssembleStory(db, "kinds", ds.Add); e != nil {
// 			t.Fatal("couldnt assemble story", e, ds.Err())
// 		} else if len(ds) > 0 {
// 			t.Fatal("issues assembling", ds.Err())
// 		} else if cnt, e := CheckAll(db, "", tapestry.AllSignatures); e != nil {
// 			t.Fatal(e)
// 		} else if cnt != 1 {
// 			t.Fatal("expected one test", cnt)
// 		} else {
// 			t.Log("ok", cnt)
// 		}
// 	}
// }

// func dbwriter(db *sql.DB) story.WriterFun {
// 	cache := tables.NewCache(db)
// 	return func(q string, args ...interface{}) {
// 		cache.Must(q, args...)
// 	}
// }

// func storyMarshaller(m jsn.Marshalee) (string, error) {
// 	return cout.Marshal(m, story.CompactEncoder)
// }
