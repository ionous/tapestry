package cmdgenerate

import (
	"flag"
	"os"
	"path/filepath"
)

const idlRelativeDir = `src/git.sr.ht/~ionous/tapestry/idl`

// collection of local flags
var genFlags = struct {
	dl  string // filter by group
	in  string // input path
	out string // output directory: defaults to "_temp"
}{}

func buildFlags() (flags flag.FlagSet) {
	defaultIn := filepath.Join(os.Getenv("GOPATH"), idlRelativeDir)
	flags.StringVar(&genFlags.dl, "dl", "", "limit to which groups")
	flags.StringVar(&genFlags.in, "in", defaultIn, "input directory containing one or more .ifspecs")
	flags.StringVar(&genFlags.out, "out", "./_temp", "output directory")
	return
}
