package tables

import (
	"database/sql"
	_ "embed"

	"github.com/ionous/errutil"
	_ "github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

const DefaultDriver = "sqlite3"

// CreateModel creates the tables listed in model.sql
func CreateModel(db *sql.DB) (err error) {
	if _, e := db.Exec(modelSql); e != nil {
		err = errutil.New("Couldn't create model tables", e)
	}
	return
}

//go:embed model.sql
var modelSql string

// CreateRun creates the tables listed in run.sql
func CreateRun(db *sql.DB) (err error) {
	if _, e := db.Exec(runSql); e != nil {
		err = errutil.New("Couldn't create run tables", e)
	}
	return
}

//go:embed run.sql
var runSql string

// CreateAll tables listed in the various .sql files.
func CreateAll(db *sql.DB) (err error) {
	if e := CreateModel(db); e != nil {
		err = e
	} else if e := CreateRun(db); e != nil {
		err = e
	}
	return
}
