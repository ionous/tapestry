package cmddoc

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/play"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/lang/doc"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func runDoc(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if outPath, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else {
		err = doc.Build(outPath, []typeinfo.TypeSet{
			call.Z_Types,
			debug.Z_Types,
			frame.Z_Types,
			game.Z_Types,
			grammar.Z_Types,
			jess.Z_Types,
			list.Z_Types,
			literal.Z_Types,
			logic.Z_Types,
			math.Z_Types,
			object.Z_Types,
			play.Z_Types,
			prim.Z_Types,
			format.Z_Types,
			rel.Z_Types,
			render.Z_Types,
			rtti.Z_Types,
			story.Z_Types,
			// testdl.Z_Types,

			text.Z_Types,
		})
	}
	return
}

var CmdDoc = &base.Command{
	Run:       runDoc,
	Flag:      buildFlags(),
	UsageLine: "tap doc [-out]",
	Short:     "generate html documentation",
	Long: `
Transforms generated idl files into api documentation.`,
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
