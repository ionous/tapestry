// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package cmdcheck

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
)

func runCheck(ctx context.Context, cmd *base.Command, args []string) (err error) {
	log.Println("Checking", checkFlags.srcPath)
	if lvl, e := getLogLevel(checkFlags.logLevel); e != nil {
		err = e
	} else {
		debug.LogLevel = lvl
		opt := qna.NewOptions()
		opt.SetOption(meta.PrintResponseNames, generic.BoolOf(checkFlags.responses))
		if cnt, e := checkFile(checkFlags.srcPath, lang.Normalize(checkFlags.checkOne), opt); e != nil {
			errutil.PrintErrors(e, func(s string) { log.Println(s) })
			if errutil.Panic {
				log.Panic("mismatched")
			}
		} else {
			log.Println("Checked", cnt, checkFlags.srcPath)
		}
	}
	return
}

var CmdCheck = &base.Command{
	Run:       runCheck,
	Flag:      buildFlags(),
	UsageLine: "tap check [-in path]",
	Short:     "tests a playable story",
	Long: `Loads an assembled database and runs (one or more) test scripts that it contains.

Runs all unit tests by default, use '-run=<name>' to run a specific one.
`,
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func checkFile(inFile, testName string, opt qna.Options) (ret int, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't open db", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else {
			ret, err = CheckAll(db, testName, opt, tapestry.AllSignatures)
		}
	}
	return
}
