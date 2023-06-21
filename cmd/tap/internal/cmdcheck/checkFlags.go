package cmdcheck

import (
	"flag"
	"os"
	"path/filepath"
)

// collection of local flags
var checkFlags = struct {
	srcPath   string
	checkOne  string
	responses bool
}{}

func buildFlags() (flags flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	flags.StringVar(&checkFlags.srcPath, "in", inPath, `input file or directory name.`)
	flags.StringVar(&checkFlags.checkOne, "run", "", "run check on a particular test")
	flags.BoolVar(&checkFlags.responses, "responses", false, "print response names instead of values")
	return
}
