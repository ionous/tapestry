package cmdgenerate

import (
	"bytes"
	"context"
	"flag"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/generate"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

var CmdGenerate = &base.Command{
	Run:       runGenerate,
	Flag:      buildFlags(),
	UsageLine: "tap gen [-out ../../dl]",
	Short:     "gen golang serializers from .tells",
	Long: `
Generates .go source code for reading and writing story files.`,
}

// FIX: where do keyword specs come from?
func runGenerate(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if outPath, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else {
		// read files:
		var groups []generate.Group
		inDir := os.DirFS(genFlags.in)
		if e := fs.WalkDir(inDir, ".", func(path string, d fs.DirEntry, e error) (err error) {
			if e != nil {
				err = e // fix: ignores specs for now ( it has swaps still )
			} else if n := d.Name(); !d.IsDir() && strings.HasSuffix(n, ".tells") && n != "spec.tells" {
				if fp, e := inDir.Open(path); e != nil {
					err = e
				} else if raw, e := files.ReadRawTell(fp); e != nil {
					err = e
				} else if msg, e := decode.ParseMessage(raw); e != nil {
					err = e
				} else if g, e := generate.ReadSpec(msg); e != nil {
					err = e
				} else {
					groups = append(groups, g)
				}
			}
			if err != nil {
				err = errutil.New(e, "in", path)
			}
			return
		}); e != nil {
			err = e
		} else if g, e := generate.MakeGenerator(groups); e != nil {
			err = e
		} else {
			// write files
			for g.Next() {
				groupName := g.Name() // optionally, limit to one particular group
				if len(genFlags.dl) == 0 || (genFlags.dl == groupName) {
					// tbd: an option to do everything in memory and write to stdout?
					var u bytes.Buffer
					if e := g.Write(&u); e != nil {
						err = e
						break
					} else if f, e := format.Source(u.Bytes()); e != nil {
						log.Println(u.String())
						err = e
						break
					} else {
						path := filepath.Join(outPath, groupName)
						os.MkdirAll(path, 0700)
						path = filepath.Join(path, groupName+"_types.go")
						println("writing", path)
						if e := os.WriteFile(path, f, 0666); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
	}
	return
}
