package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/doc"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func main() {
	flags := buildFlags()
	flags.Parse(os.Args)
	if e := runDoc(); e != nil {
		log.Fatal(e)
	}
}

func runDoc() (err error) {
	if outPath, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		err = e
	} else {
		err = doc.Build(outPath, []typeinfo.TypeSet{
			story.Z_Types,
		})
	}
	return
}

// collection of local flags
var genFlags = struct {
	out string // output directory
}{}

func buildFlags() (fs flag.FlagSet) {
	var outPath string
	if home, e := os.UserHomeDir(); e == nil {
		outPath = filepath.Join(home, "Documents", "Tapestry", "doc")
	}
	fs.StringVar(&genFlags.out, "out", outPath, "output directory")
	return
}
