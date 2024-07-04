package dump

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryAllScenes(db *sql.DB) (ret []raw.SceneData, err error) {
	if rows, e := db.Query(must("domains")); e != nil {
		err = e
	} else {
		ret, err = scanScenes(rows)
	}
	return
}

func scanScenes(rows *sql.Rows) (ret []raw.SceneData, err error) {
	var last string
	var pending []string
	flush := func() {
		if len(last) > 0 {
			ret = append(ret, raw.SceneData{
				Scene:    last,
				Requires: pending,
			})
			pending = nil
		}
	}
	var domain, requires string
	if e := tables.ScanAll(rows, func() (_ error) {
		if domain != last {
			flush()
			last = domain
		}
		if len(requires) > 0 {
			pending = append(pending, requires)
		}
		return
	}, &domain, &requires); e != nil {
		err = e
	} else {
		flush()
	}
	return
}
