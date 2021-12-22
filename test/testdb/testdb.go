package testdb

import (
	"database/sql"
	"io"
	"os/user"
	"path"
	"strings"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
	"github.com/ionous/sliceOf"
)

const Memory = "file:test.db?cache=shared&mode=memory"

// from testing.T.Name() return a local db testing path
func PathFromName(name string) (ret string, err error) {
	rest := strings.Replace(name, "/", ".", -1) + ".db"
	if user, e := user.Current(); e != nil {
		err = errutil.New(e, "for", name)
	} else {
		ret = path.Join(user.HomeDir, rest)
	}
	return
}

// given a table name, and list of column names return a super-secret helper.
// ( used for Ins, and WriteCsv )
func TableCols(table_cols ...string) []string {
	return table_cols
}

// Opens a (sqlite) database in memory or on disk, panicking on error.
// To create a memory table: pass "Memory" as the path.
// If path is empty, uses the users's home directory.
// If driver is empty, assumes a sqlite database.
// If the db/file does not already exist, it will be created.
func Open(name, path, driver string) (ret *sql.DB) {
	if len(driver) == 0 {
		driver = tables.DefaultDriver
	}
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := PathFromName(name); e != nil {
		panic(e)
	} else {
		source = p
	}
	if db, e := sql.Open(driver, source); e != nil {
		panic(e)
	} else {
		ret = db
	}
	return
}

// insert an arbitrary number of rows into the passed db.
// tablecols holds the names of the table and columns to query,
// els can hold multiple rows of data, each containing the number of cols specified by tablecols.
func Ins(db tables.Executer, tablecols []string, els ...interface{}) (err error) {
	ins, width := tables.Insert(tablecols[0], tablecols[1:]...), len(tablecols)-1
	for i, cnt := 0, len(els); i < cnt; i += width {
		if _, e := db.Exec(ins, els[i:i+width]...); e != nil {
			err = e
			break
		}
	}
	return
}

// query the passed db and write the results to w --
// builds the query from "tablecols" which holds the names of the table and columns to query;
// "where" can filter that data. ( see also: tables.WriteCsv. )
func WriteCsv(db tables.Querier, w io.Writer, tablecols []string, where string) (err error) {
	table, cols := tablecols[0], strings.Join(tablecols[1:], ", ")
	q := strings.Join(sliceOf.String("select", cols, "from", table, where, "order by", cols), " ")
	return tables.WriteCsv(db, w, q, len(tablecols)-1)
}
