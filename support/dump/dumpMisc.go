package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/tables"
)

// fix? this is a little different than the way the normal player works
// this includes kinds
func QueryPlurals(db *sql.DB, scene string) (ret []Plural, err error) {
	var p Plural
	if rows, e := db.Query(must("plurals"), scene); e != nil {
		err = fmt.Errorf("%w while querying plurals", e)
	} else {
		err = tables.ScanAll(rows, func() (_ error) {
			ret = append(ret, p)
			return
		}, &p.One, &p.Other)
	}
	return
}
