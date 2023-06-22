package cmdcheck

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/debug"
	"github.com/ionous/errutil"
)

// collection of local flags
var checkFlags = struct {
	srcPath   string
	checkOne  string
	responses bool
	logLevel  string
}{}

func buildFlags() (flags flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	spec := debug.LogLevel.Compose()
	levels := strings.Join(spec.Strings, ", ")
	flags.StringVar(&checkFlags.srcPath, "in", inPath, `input file or directory name.`)
	flags.StringVar(&checkFlags.checkOne, "run", "", "run check on a particular test")
	flags.BoolVar(&checkFlags.responses, "responses", false, "print response names instead of values")
	flags.StringVar(&checkFlags.logLevel, "log", "", levels)
	return
}

func getLogLevel(in string) (ret debug.LoggingLevel, err error) {
	if len(in) > 0 {
		spec := ret.Compose()
		if key, idx := spec.IndexOfValue(in); idx < 0 {
			err = errutil.New("Unknown log level")
		} else {
			ret.Str = key
		}
	}
	return
}
