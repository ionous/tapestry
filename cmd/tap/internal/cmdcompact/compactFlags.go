package cmdcompact

import "flag"

// collection of local flags
var compactFlags = struct {
	inPath, outPath, inExts, outExt string
	recurse, pretty                 bool
}{}

func buildFlags() (flags flag.FlagSet) {
	flags.StringVar(&compactFlags.outPath, "out", "./_temp", "output directory")
	flags.StringVar(&compactFlags.inPath, "in", "", "input file(s) or paths(s) (comma separated)")
	flags.BoolVar(&compactFlags.recurse, "recurse", false, "scan input sub-directories")
	flags.StringVar(&compactFlags.inExts, "filter", ".if",
		`extension(s) for directory scanning.
ignored if 'in' refers to a specific file`)
	flags.StringVar(&compactFlags.outExt, "convert", "",
		`an optional file extension to force a story format conversion (.if|.ifx|.block)
underscores are allowed to avoid copying over the original files. (._if, .if_, etc.)
( ex. if the in and out directories are the same.
if no extension is specified, the output format is the same as the import format.`)
	flags.BoolVar(&compactFlags.pretty, "pretty", false, "make the output somewhat human readable")
	return
}
