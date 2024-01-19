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

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/generate"
	"git.sr.ht/~ionous/tapestry/support/files"
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
				err = e
			} else if !d.IsDir() { // the first dir we get is "."
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
			return
		}); e != nil {
			err = e
		} else {
			// write files
			for g := generate.MakeGenerator(groups); g.Next(); {
				groupName := g.Name()
				if len(genFlags.dl) == 0 || (genFlags.dl == groupName) {
					// tbd: an option to do everything in memory and write to stdout?
					var b bytes.Buffer
					if e := g.Write(&b); e != nil {
						err = e
						break
					} else if b, e := format.Source(b.Bytes()); e != nil {
						err = e
						break
					} else {
						path := filepath.Join(outPath, groupName)
						// are we in the dl directory, and writing to rt?
						// fix: since there are no rt commands, only slots:
						// writing to dl/rt ( or rtypes ) should be fine
						// maybe even rename tells to "rtypes", "rtslots", "rti", to simplify
						// dlLike := strings.HasSuffix(path, string(filepath.Separator)+"dl")
						// if groupName == "rt" && dlLike {
						// 	path = path[:len(path)-2]
						// }
						os.MkdirAll(path, 0700)
						path = filepath.Join(path, groupName+"_types.go")
						println("writing", path)
						if e := os.WriteFile(path, b, 0666); e != nil {
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
