package tables

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

const defaultDriver = "sqlite3"
const tapestryDriver = "tapestry" // a pair of databases main, and rt

func init() {
	sql.Register(tapestryDriver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) (err error) {
			_, err = conn.Exec("attach database ':memory:' as rt", nil)
			return
		},
	})
}
