package cmdgenerate

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/generate"
	"github.com/ionous/errutil"
	"golang.org/x/tools/imports"
)

var CmdGenerate = &base.Command{
	Run:       runGenerate,
	Flag:      buildFlags(),
	UsageLine: "tap generate [-out ../../dl]",
	Short:     "generate golang serializers from .ifspecs",
	Long: `
Make generates .go language serialization code for reading and writing .if files.`,
}

// FIX: where do keyword specs come from?
func runGenerate(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if path, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else {
		dlLike := strings.HasSuffix(path, string(filepath.Separator)+"dl")
		if e := generate.WriteSpecs(os.DirFS(genFlags.in), func(groupName string, b []byte) (err error) {
			if len(genFlags.dl) == 0 || (genFlags.dl == groupName) {
				// fix: an option to do everything in memory and write to stdout?
				path := path
				if groupName == "rt" && dlLike {
					path = path[:len(path)-2]
				}
				base := filepath.Join(path, groupName)
				os.MkdirAll(base, 0700)
				path = filepath.Join(base, groupName+"_lang.go")
				// uses goimports to add imports and format the source
				if formatted, e := imports.Process(path, b, nil); e != nil {
					fmt.Println(string(b))
					err = errutil.New(e, "while formatting", groupName)
				} else {
					b = formatted
				}
				if fp, e := os.Create(path); e != nil {
					err = e // writing errors take precedence over the formatting error
				} else {
					println("writing", path)
					fp.Write(b)
					fp.Close()
				}
			}
			return
		}); e != nil {
			log.Fatalf("error generating %q from %q %v", path, genFlags.in, e)
		}
	}
	return
}
