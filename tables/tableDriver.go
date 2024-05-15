package tables

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

const defaultDriver = "sqlite3"
const tapestryDriver = "tapestry" // a pair of databases main, and rt
const memory = ":memory:"

func init() {
	sql.Register(tapestryDriver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) (err error) {
			_, err = conn.Exec("attach database ':memory:' as rt", nil)
			return
		},
	})
}

func open(driver, uri string) (ret *sql.DB, err error) {
	if db, e := sql.Open(driver, uri); e != nil {
		err = e
	} else {
		// trying to avoid potential issues with connections creating
		// multiple memory databases, multiple temp databases, etc.
		db.SetMaxOpenConns(1)
		ret = db
	}
	return
}
