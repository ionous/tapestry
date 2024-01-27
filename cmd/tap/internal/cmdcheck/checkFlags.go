package cmdcheck

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/debug"
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
	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	flags.StringVar(&checkFlags.srcPath, "in", inPath, `input file or directory name.`)
	flags.StringVar(&checkFlags.checkOne, "run", "", "run check on a particular test")
	flags.BoolVar(&checkFlags.responses, "responses", false, "print response names instead of values")
	flags.StringVar(&checkFlags.logLevel, "log", debug.C_LoggingLevel_Note.String(), levels)
	return
}
