package cmdweave

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cmdcheck"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/tables"
)

func CheckOutput(inFile, testName string) (ret int, err error) {
	if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = fmt.Errorf("couldn't open db %q because %s", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else {
			opt := qna.NewOptions()
			ret, err = cmdcheck.CheckAll(db, testName, opt, tapestry.AllSignatures)
		}
	}
	return
}
