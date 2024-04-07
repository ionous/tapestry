// Runs tests on an existing story database.
package cmdcheck

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func runCheck(ctx context.Context, cmd *base.Command, args []string) (err error) {
	checkOne := strings.Join(args, " ")
	log.Println("Checking", checkFlags.srcPath, checkOne)
	if lvl, ok := debug.MakeLoggingLevel(checkFlags.logLevel); !ok {
		err = errutil.New("Unknown log level", checkFlags.logLevel)
	} else {
		debug.LogLevel = lvl
		opt := qna.NewOptions()
		opt.SetOption(meta.PrintResponseNames, generic.BoolOf(checkFlags.responses))
		if cnt, e := checkFile(checkFlags.srcPath, inflect.Normalize(checkOne), opt); e != nil {
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

var CmdCheck = &base.Command{
	Run:       runCheck,
	Flag:      buildFlags(),
	UsageLine: "tap check [-in path] [name]",
	Short:     "run tests on existing stories",
	Long: `Loads an playable database and runs (one or more) test scripts that it contains.

Runs all unit tests by default, can specify a name to run a specific one.
`,
}

// collection of local flags
var checkFlags = struct {
	srcPath   string
	responses bool
	logLevel  string
}{}

func buildFlags() (flags flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	flags.StringVar(&checkFlags.srcPath, "in", inPath, `input file or directory name.`)
	flags.BoolVar(&checkFlags.responses, "responses", false, "print response names instead of values")
	flags.StringVar(&checkFlags.logLevel, "log", debug.C_LoggingLevel_Debug.String(), levels)
	return
}
