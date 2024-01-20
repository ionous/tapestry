package cmdgenerate

import (
	"flag"
	"os"
	"path/filepath"
)

// collection of local flags
var genFlags = struct {
	dl  string // filter by group
	in  string // input path
	out string // output directory
}{}

func buildFlags() (flags flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "idl")
	}
	flags.StringVar(&genFlags.dl, "dl", "", "limit to which groups")
	flags.StringVar(&genFlags.in, "in", inPath, "input directory containing one or more spec files")
	flags.StringVar(&genFlags.out, "out", "../../dl", "output directory")
	return
}
