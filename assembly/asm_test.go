package assembly

import (
	"database/sql"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
)

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest) {
	db := testdb.Open(t.Name(), path, SqlCustomDriver)
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateAssembly(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else {
		var ds reader.Dilemmas
		rec := ephemera.NewRecorder(db).SetSource(t.Name())
		mdl := NewAssemblerReporter(db, ds.Add)
		ret = &assemblyTest{
			T:         t,
			db:        db,
			rec:       rec,
			assembler: mdl,
			dilemmas:  &ds,
		}
	}
	return
}

type assemblyTest struct {
	*testing.T
	db        *sql.DB
	rec       *ephemera.Recorder
	assembler *Assembler
	dilemmas  *reader.Dilemmas
}
