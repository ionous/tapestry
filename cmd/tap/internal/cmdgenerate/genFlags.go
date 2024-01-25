package cmdgenerate

import (
	"flag"
	"os"
	"path/filepath"
)

// collection of local flags
var genFlags = struct {
	dl     string // filter by group
	in     string // input path
	out    string // output directory
	useDB  bool
	dbPath string // output file or path
}{}

func buildFlags() (fs flag.FlagSet) {
	var inPath string
	if home, e := os.UserHomeDir(); e == nil {
		inPath = filepath.Join(home, "Documents", "Tapestry", "idl")
	}
	var dbPath string
	if home, e := os.UserHomeDir(); e == nil {
		dbPath = filepath.Join(home, "Documents", "Tapestry", "build", "idl.db")
	}
	fs.StringVar(&genFlags.dl, "dl", "", "limit to which groups")
	fs.StringVar(&genFlags.in, "in", inPath, "input directory containing one or more spec files")
	fs.StringVar(&genFlags.out, "out", "../../dl", "output directory")
	fs.BoolVar(&genFlags.useDB, "db", false, "generate a sqlite representation")
	fs.StringVar(&genFlags.dbPath, "dbFile", dbPath, "sqlite output file")
	return
}
