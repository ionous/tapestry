package qna

// FIX: have to remap using the asm table writer...
// no idea where this test should live...
// tests the importation, assembly, and execution of a factorial story.
// doesn't test the *reading* of the story.
// func TestFullFactorial(t *testing.T) {
// 	db := testdb.Create(t.Name())
// 	defer db.Close()

// 	// read factorialStory, assemble and run.
// 	var ds reader.Dilemmas
// 	if e := tables.CreateAll(db); e != nil {
// 		t.Fatal("couldn't create tables", e)
// 	} else {
// 		k := weave.NewCatalog(dbwriter(db))
// 		if e := story.ImportStory(k, t.Name(), &debug.FactorialStory); e != nil {
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

// type WriterFun func(eph Ephemera)

// func (m mdl.Modeler Fun) WriteEphemera(op Ephemera) {
// 	w(op)
// }
