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
	log.Println("Checking", cfg.srcPath, checkOne)
	if lvl, ok := debug.MakeLoggingLevel(cfg.logLevel); !ok {
		err = errutil.New("Unknown log level", cfg.logLevel)
	} else if srcPath, e := filepath.Abs(cfg.srcPath); e != nil {
		err = e
	} else {
		debug.LogLevel = lvl
		opt := qna.NewOptions()
		opt.SetOption(meta.PrintResponseNames, rt.BoolOf(cfg.responses))
		if cnt, e := CheckFile(srcPath, inflect.Normalize(checkOne), opt); e != nil {
			errutil.PrintErrors(e, func(s string) { log.Println(s) })
			if errutil.Panic {
				log.Panic("mismatched")
			}
		} else {
			log.Println("Checked", cnt, cfg.srcPath)
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

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	srcPath   string
	responses bool
	logLevel  string
}{}

func buildFlags() (ret flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	ret.StringVar(&cfg.srcPath, "in", inPath, `input file or directory name.`)
	ret.BoolVar(&cfg.responses, "responses", false, "print response names instead of values")
	ret.StringVar(&cfg.logLevel, "log", debug.C_LoggingLevel_Debug.String(), levels)
	return
}
