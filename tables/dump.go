package tables

import (
	"database/sql"
	"io"
)

// query data from the passed db, and write the comma separated results to the passed writer.
// q holds the query string, cols holds the number of expected columns in the output.
// ex. WriteCsv(db, os.Stdout, "select col1, col2 from table", 2)
func WriteCsv(db Querier, w io.Writer, q string, cols int) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = e
	} else {
		c := make([]sql.NullString, cols)
		cp := make([]any, cols)
		for i := 0; i < cols; i++ {
			cp[i] = &c[i]
		}
		err = ScanAll(rows, func() (err error) {
			for i, col := range c {
				if i > 0 {
					io.WriteString(w, ",")
				}
				if !col.Valid {
					io.WriteString(w, "NULL")
				} else {
					io.WriteString(w, col.String)
				}
			}
			io.WriteString(w, "\n")
			return
		}, cp...) // pass the pointers to the column strings
	}
	return
}

// where each row is one string.
func QueryStrings(db Querier, q string, args ...any) (ret []string, err error) {
	if rows, e := db.Query(q, args...); e != nil {
		err = e
	} else {
		ret, err = ScanStrings(rows)
	}
	return
}

func ScanStrings(rows *sql.Rows) (ret []string, err error) {
	var str sql.NullString
	err = ScanAll(rows, func() (_ error) {
		if !str.Valid {
			str.String = "NULL"
		}
		ret = append(ret, str.String)
		return
	}, &str)
	return
}
