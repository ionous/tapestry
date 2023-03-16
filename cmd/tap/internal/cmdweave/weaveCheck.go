package cmdweave

import (
	"database/sql"
	check "git.sr.ht/~ionous/tapestry/cmd/tap/internal/cmdcheck"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func CheckOutput(inFile, testName string) (ret int, err error) {
	if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't open db", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else {
			opt := qna.NewOptions()
			ret, err = check.CheckAll(db, testName, opt, tapestry.AllSignatures)
		}
	}
	return
}
