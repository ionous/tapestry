package tables

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
)

// Query used for QueryAll to hide differences b/t tables.Cache and sql.DB
type Querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// Query used for QueryAll to hide differences b/t tables.Cache and sql.DB
type QueryRow interface {
	QueryRow(query string, args ...any) *sql.Row
}

// compatible with sql.DB for use with caches, etc.
type Executer interface {
	Exec(q string, args ...any) (sql.Result, error)
}

func Must(db *sql.DB, q string, args ...any) {
	if _, e := db.Exec(q, args...); e != nil {
		panic(e)
	}
}

func RowsAffected(res sql.Result) (ret int) {
	if cnt, e := res.RowsAffected(); e != nil {
		ret = -1
	} else {
		ret = int(cnt)
	}
	return
}

// QueryAll queries for one or more rows.
// For each row, it writes the row to the 'dest' args and calls 'cb' for processing.
func QueryAll(db Querier, q string, cb func() error, dest ...any) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = errutil.New("QueryAll error:", e, "for", q)
	} else {
		err = ScanAll(rows, cb, dest...)
	}
	return
}

// ScanAll writes each row to the 'dest' args and calls 'cb' for processing.
// It closes rows before returning.
func ScanAll(rows *sql.Rows, cb func() error, dest ...any) (err error) {
	for rows.Next() {
		if e := rows.Scan(dest...); e != nil {
			err = errutil.New("ScanAll error:", e)
			break
		} else if e := cb(); e != nil {
			err = e
			break
		}
	}
	if e := rows.Err(); e != nil {
		err = errutil.Append(err, e)
	}
	rows.Close()
	return
}

// Insert creates a sqlite friendly insert statement.
// For example: "insert into foo(col1, col2, ...) values(?, ?, ...)"
func Insert(table string, keys ...string) string {
	return InsertWith(table, "", keys...)
}

// InsertWith allows the specification of on conflict directives
func InsertWith(table string, rest string, keys ...string) string {
	vals := "?"
	if kcnt := len(keys) - 1; kcnt > 0 {
		vals += strings.Repeat(",?", kcnt)
	}
	return "INSERT into " + table +
		"(" + strings.Join(keys, ", ") + ")" +
		" values " + "(" + vals + ")" + rest + ";"
}

// insert an arbitrary number of rows into the passed db.
// tablecols holds the names of the table and columns to query,
// els can hold multiple rows of data, each containing the number of cols specified by tablecols.
func Ins(db Executer, tablecols []string, els ...interface{}) (err error) {
	ins, width := Insert(tablecols[0], tablecols[1:]...), len(tablecols)-1
	for i, cnt := 0, len(els); i < cnt; i += width {
		if _, e := db.Exec(ins, els[i:i+width]...); e != nil {
			err = e
			break
		}
	}
	return
}
