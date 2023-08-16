package cmdplay

import (
	"flag"
	"os"
	"path/filepath"
)

// collection of local flags
var playFlags = struct {
	inFile, testString, scene string
	json, debugging           bool
}{}

func buildFlags() (flags flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}

	flags.StringVar(&playFlags.inFile, "in", inFile, "input file name (sqlite3)")
	flags.StringVar(&playFlags.scene, "scene", "tapestry", "scene to start playing")
	flags.StringVar(&playFlags.testString, "test", "", "optional list of commands to run (non-interactive)")
	flags.BoolVar(&playFlags.json, "json", false, "expect input/output in json (default is plain text)")
	flags.BoolVar(&playFlags.debugging, "debug", false, "extra debugging output?")
	return
}
