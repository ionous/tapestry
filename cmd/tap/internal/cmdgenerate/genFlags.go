package cmdgenerate

import (
	"flag"
)

// collection of local flags
var genfLags = struct {
	dl  string // filter by group
	out string // output directory: defaults to "_temp"
}{}

func buildFlags() (flags flag.FlagSet) {
	flag.StringVar(&genfLags.dl, "dl", "", "limit to which groups")
	flag.StringVar(&genfLags.out, "out", "./_temp", "output directory")
	return
}
