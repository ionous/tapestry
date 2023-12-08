package cmdcompact

import (
	"flag"
	"path/filepath"
)

// collection of local flags
var compactFlags flags

type flags struct {
	inPath, outPath, inExts, outExt string
	recurse, pretty                 bool
}

func buildFlags() (fs flag.FlagSet) {
	fs.StringVar(&compactFlags.outPath, "out", "./_temp", "output directory")
	fs.StringVar(&compactFlags.inPath, "in", "", "input file(s) or paths(s) (comma separated)")
	fs.BoolVar(&compactFlags.recurse, "recurse", false, "scan input sub-directories")
	fs.StringVar(&compactFlags.inExts, "filter", ".if",
		`extension(s) for directory scanning.
ignored if 'in' refers to a specific file`)
	fs.StringVar(&compactFlags.outExt, "convert", "",
		`an optional file extension to force a story format conversion (.if|.ifx|.block|.tell|.tells)
underscores are allowed to avoid copying over the original files. (._if, .if_, etc.)
( ex. if the in and out directories are the same.
if no extension is specified, the output format is the same as the import format.`)
	fs.BoolVar(&compactFlags.pretty, "pretty", false, "make the output somewhat human readable")
	return
}

// given an input filename (without path, but with extension)
// return the desired output filename.
func (f flags) replaceExt(name string) (ret string) {
	if len(compactFlags.outExt) == 0 {
		ret = name
	} else {
		fileExt := filepath.Ext(name)
		ret = name[:len(name)-len(fileExt)] + compactFlags.outExt
	}
	return
}
