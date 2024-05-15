package tables

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

// reads from the passed file path,
// overwriting the dynamic parts of the dst database.
// ( the dynamic portions are all in the rt schema )
func LoadFile(dst *sql.DB, fromFile string) (err error) {
	if src, e := open(defaultDriver, fromFile); e != nil {
		err = e
	} else {
		err = copyDB(dst, rtData, src, mainData)
		src.Close()
	}
	return
}

// writes the dynamic portions of the src database to the passed file path.
// removes the file if it created
func SaveFile(toFile string, force bool, src *sql.DB) (err error) {
	if _, e := os.Stat(toFile); e == nil && !force {
		err = fmt.Errorf("can't save to %s, it already exists", toFile)
	} else if e != nil && !errors.Is(e, fs.ErrNotExist) {
		err = e
	} else if dst, e := open(defaultDriver, toFile); e != nil {
		err = e
	} else {
		err = copyDB(dst, mainData, src, rtData)
		dst.Close()
		if err != nil {
			os.Remove(toFile)
		}
	}
	return
}

const (
	rtData   = "rt"
	mainData = "main"
)

// ported from https://www.sqlite.org/backup.html
//
// load the contents of a file into an open database
// or save the contents of an open db into a file.
// in either case, the destination is completely overwritten.
func copyDB(dst *sql.DB, dstName string, src *sql.DB, srcName string) error {
	return SqliteConn(dst, func(to *sqlite3.SQLiteConn) error {
		return SqliteConn(src, func(from *sqlite3.SQLiteConn) (err error) {
			if bk, e := to.Backup(dstName, from, srcName); e != nil {
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

// open a connection to the db and calls the passed cb.
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
