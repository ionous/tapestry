package tables

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	_ "github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

// name without extension
func createTables(db *sql.DB, names ...string) (err error) {
	for _, name := range names {
		if b, e := fs.ReadFile(sqlFs, "sql/"+name+".sql"); e != nil {
			err = e
			break
		} else if _, e := db.Exec(string(b)); e != nil {
			err = fmt.Errorf("couldn't create %s because %v", name, e)
			break
		}
	}
	return
}

//go:embed sql
var sqlFs embed.FS
