// Runs tests on an existing story database.
package cmdcheck

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

func runCheck(ctx context.Context, cmd *base.Command, args []string) (err error) {
	checkOne := strings.Join(args, " ")
	log.Println("Checking", checkFlags.srcPath, checkOne)
	if lvl, ok := debug.MakeLoggingLevel(checkFlags.logLevel); !ok {
		err = errutil.New("Unknown log level", checkFlags.logLevel)
	} else if srcPath, e := filepath.Abs(checkFlags.srcPath); e != nil {
		err = e
	} else {
		debug.LogLevel = lvl
		opt := qna.NewOptions()
		opt.SetOption(meta.PrintResponseNames, rt.BoolOf(checkFlags.responses))
		if cnt, e := CheckFile(srcPath, inflect.Normalize(checkOne), opt); e != nil {
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
