package cmdplay

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ionous/errutil"
)

// collection of local flags
var playFlags = struct {
	inFile, testString, domain string
	json, debugging            bool
}{}

func buildFlags() (flags flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}

	flag.StringVar(&playFlags.inFile, "in", inFile, "input file name (sqlite3)")
	flag.StringVar(&playFlags.domain, "scene", "tapestry", "scene to start playing")
	flag.StringVar(&playFlags.testString, "test", "", "optional list of commands to run (non-interactive)")
	flag.BoolVar(&playFlags.json, "json", false, "expect input/output in json (default is plain text)")
	flag.BoolVar(&playFlags.debugging, "debug", false, "extra debugging output?")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	return
}
