package cmdweave

import (
	"flag"
	"os"
	"path/filepath"
)

// collection of local flags
var weaveFlags = struct {
	srcPath, outFile string
	checkAll         bool
	checkOne         string
}{}

func buildFlags() (flags flag.FlagSet) {
	var inPath string
	var outPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "stories", "shared")
		outPath = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
	}
	flags.StringVar(&weaveFlags.srcPath, "in", inPath, `input file or directory name.`)
	flags.StringVar(&weaveFlags.outFile, "out", outPath, "optional output filename (sqlite3)")
	flags.BoolVar(&weaveFlags.checkAll, "check", false, "run check after importing?")
	flags.StringVar(&weaveFlags.checkOne, "run", "", "run check on one test after importing?")
	return
}
