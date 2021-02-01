package assembly

import (
	"database/sql"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/tables"
	"git.sr.ht/~ionous/iffy/test/testdb"
	"github.com/ionous/errutil"
)

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest, err error) {
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		err = e
	} else {
		source = p
	}
	if err == nil {
		if db, e := sql.Open(SqlCustomDriver, source); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateEphemera(db); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateAssembly(db); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateModel(db); e != nil {
			err = errutil.New(e, "for", source)
		} else {
			var ds reader.Dilemmas
			rec := ephemera.NewRecorder(t.Name(), db)
			mdl := NewAssemblerReporter(db, ds.Add)
			ret = &assemblyTest{
				T:         t,
				db:        db,
				rec:       rec,
				assembler: mdl,
				dilemmas:  &ds,
			}
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
