package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/play"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
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
			assign.Z_Types,
			core.Z_Types,
			debug.Z_Types,
			frame.Z_Types,
			game.Z_Types,
			grammar.Z_Types,
			jess.Z_Types,
			list.Z_Types,
			literal.Z_Types,
			play.Z_Types,
			prim.Z_Types,
			rel.Z_Types,
			render.Z_Types,
			rtti.Z_Types,
			story.Z_Types,
			// testdl.Z_Types,
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
		outPath = filepath.Join(home, "Documents", "Tapestry", "api")
	}
	fs.StringVar(&genFlags.out, "out", outPath, "output directory")
	return
}
