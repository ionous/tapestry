package tables

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

func LoadFile(dst *sql.DB, fromFile string) (err error) {
	if src, e := sql.Open(defaultDriver, fromFile); e != nil {
		err = e
	} else {
		err = copyDB(dst, src, rtTables)
		src.Close()
	}
	return
}

func SaveFile(toFile string, src *sql.DB) (err error) {
	if dst, e := sql.Open(defaultDriver, toFile); e != nil {
		err = e
	} else {
		err = copyDB(dst, src, rtTables)
		dst.Close()
	}
	return
}

const rtTables = "rt"

// ported from https://www.sqlite.org/backup.html
//
// load the contents of a file into an open database
// or save the contents of an open db into a file.
// in either case, the destination is completely overwritten.
func copyDB(dst *sql.DB, src *sql.DB, dbName string) error {
	return SqliteConn(dst, func(to *sqlite3.SQLiteConn) error {
		return SqliteConn(src, func(from *sqlite3.SQLiteConn) (err error) {
			if bk, e := to.Backup(dbName, from, dbName); e != nil {
				err = e
			} else if ok, e := bk.Step(-1); e != nil {
				err = e // ^ -1 copies everything all at once
			} else if !ok {
				err = errors.New("unknown error in backup")
			} else if e := bk.Finish(); e != nil {
				err = e
			}
			return
		})
	})
}

// automatically closes the connection after the callback is done.
func SqliteConn(db *sql.DB, cb func(conn *sqlite3.SQLiteConn) error) (err error) {
	if conn, e := db.Conn(context.Background()); e != nil {
		err = e
	} else {
		err = conn.Raw(func(raw any) error {
			return cb(raw.(*sqlite3.SQLiteConn))
		})
		conn.Close()
	}
	return
}
