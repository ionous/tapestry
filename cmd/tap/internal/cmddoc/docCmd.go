package cmddoc

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/lang/doc"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func runDoc(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if outPath, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else {
		err = doc.Build(outPath, []typeinfo.TypeSet{
			story.Z_Types,
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
