package internal

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

//go:generate templify -p internal -o testData.gen.go testData.sql
func CreateTestData(db *sql.DB) (err error) {
	if e := tables.CreateModel(db); e != nil {
		err = e
	} else if _, e := db.Exec(testDataTemplate()); e != nil {
		err = errutil.New("createTestData", e)
	}
	return
}
