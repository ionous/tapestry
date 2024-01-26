package cmdplay

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/debug"
)

var cfg = struct {
	inFile, testString, scene  string
	json, debugging, responses bool
	logLevel                   string
}{}

func buildFlags() (fs flag.FlagSet) {
	var inFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}

	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	fs.StringVar(&cfg.inFile, "in", inFile, "input file name (sqlite3)")
	fs.StringVar(&cfg.scene, "scene", "tapestry", "scene to start playing")
	fs.StringVar(&cfg.testString, "test", "", "optional list of commands to run (non-interactive)")
	fs.BoolVar(&cfg.json, "json", false, "expect input/output in json (default is plain text)")
	fs.BoolVar(&cfg.responses, "responses", false, "print response names instead of values")
	fs.StringVar(&cfg.logLevel, "log", "", levels)
	return
}
