package tables

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

// creates a db for the tables listed in idl.sql
func CreateIdl(idlFile string) (ret *sql.DB, err error) {
	if db, e := open(defaultDriver, idlFile); e != nil {
		err = fmt.Errorf("couldn't open db %s because %v", idlFile, e)
	} else {
		err = createTables(db, "idl")
	}
	return
}

func CreateBuildTime(mdlFile string) (ret *sql.DB, err error) {
	if db, e := open(tapestryDriver, mdlFile); e != nil {
		err = e
	} else {
		if e := createTables(db,
			"model", "modelView",
			"runTables", "runViews"); e != nil {
			err = e
		} else {
			ret = db
		}
		if err != nil {
			db.Close()
		}
	}
	return
}

// reuse the model file as the runtime file
// ( useful sometimes for debugging )
func ShareRunTime(mdlFile string) (ret *sql.DB, err error) {
	if db, e := open(defaultDriver, mdlFile); e != nil {
		err = fmt.Errorf("couldn't open db %s because %v", mdlFile, e)
	} else {
		if e := shareTables(db,
			"runTables", "runViews"); e != nil {
			err = e
		} else {
			ret = db
		}
		if err != nil {
			db.Close()
		}
	}
	return
}

// open the model as read-only, and
// create a memory writable database for the runtime
func CreateRunTime(mdlFile string) (ret *sql.DB, err error) {
	if db, e := open(tapestryDriver, mdlFile+"?mode=ro"); e != nil {
		err = e
	} else {
		if e := createTables(db,
			"runTables", "runViews"); e != nil {
			err = e
		} else {
			ret = db
		}
		if err != nil {
			db.Close()
		}
	}
	return
}

func OpenModel(mdlFile string) (ret *sql.DB, err error) {
	return open(tapestryDriver, mdlFile+"?mode=ro")
}
