package cmdcompact

import "flag"

// collection of local flags
var compactFlags flags

type flags struct {
	inPath, outPath, inExts, outExt string
	recurse, pretty, tell           bool
}

func buildFlags() (fs flag.FlagSet) {
	fs.StringVar(&compactFlags.outPath, "out", "./_temp", "output directory")
	fs.StringVar(&compactFlags.inPath, "in", "", "input file(s) or paths(s) (comma separated)")
	fs.BoolVar(&compactFlags.recurse, "recurse", false, "scan input sub-directories")
	fs.StringVar(&compactFlags.inExts, "filter", ".if",
		`extension(s) for directory scanning.
ignored if 'in' refers to a specific file`)
	fs.StringVar(&compactFlags.outExt, "convert", "",
		`an optional file extension to force a story format conversion (.if|.ifx|.block)
underscores are allowed to avoid copying over the original files. (._if, .if_, etc.)
( ex. if the in and out directories are the same.
if no extension is specified, the output format is the same as the import format.`)
	fs.BoolVar(&compactFlags.pretty, "pretty", false, "make the output somewhat human readable")
	fs.BoolVar(&compactFlags.tell, "tell", false, "use the yaml-like tell format. otherwise, use json")
	return
}

// output format style
func (f flags) format() (ret format) {
	if f.tell {
		ret = useTellFormat
	} else if f.pretty {
		ret = indentedJson
	} else {
		ret = unindentedJson
	}
	return
}
