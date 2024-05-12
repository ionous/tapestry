package cmdgenerate

import (
	"bytes"
	"context"
	"flag"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/generate"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func runGenerate(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if outPath, e := filepath.Abs(genFlags.out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else if db, e := createDB(genFlags.useDB, genFlags.dbPath); e != nil {
		err = e
	} else {
		defer db.Close()
		if groups, e := readGroups(genFlags.in); e != nil {
			err = e
		} else if g, e := generate.MakeGenerator(groups); e != nil {
			err = e
		} else if e := writeGroups(g, db, outPath); e != nil {
			err = e
		} else if genFlags.useDB {
			for g.Next() {
				g.WriteReferences(db)
			}
		}
	}
	return
}

var CmdGenerate = &base.Command{
	Run:       runGenerate,
	Flag:      buildFlags(),
	UsageLine: "tap gen [-in ../../idl] [-out ../../dl] [-db -dbFile]",
	Short:     "extend tapestry with new golang code",
	Long: `
Generates .go source code for reading and writing story files from .tells files.`,
}

// collection of local flags
var genFlags = struct {
	dl     string // filter by group
	in     string // input path
	out    string // output directory
	useDB  bool
	dbPath string // output file or path
}{}

// fix: maybe this should match go generate style (ish)
// where it operates on a directory ( or file )
// and move idl files into the their own directories
// you could break out the /dl if needed or as appropriate
// out would default to the same directory as in
// ( either see if that directory code is portable
// | or handle the more frequent patterns: ".", "...", and a directory )
func buildFlags() (fs flag.FlagSet) {
	var dbPath string
	if home, e := os.UserHomeDir(); e == nil {
		dbPath = filepath.Join(home, "Documents", "Tapestry", "build", "idl.db")
	}
	fs.StringVar(&genFlags.dl, "dl", "", "limit to which groups")
	fs.StringVar(&genFlags.in, "in", "../../idl", "input directory containing one or more spec files")
	fs.StringVar(&genFlags.out, "out", "../../dl", "output directory")
	fs.BoolVar(&genFlags.useDB, "db", false, "generate a sqlite representation")
	fs.StringVar(&genFlags.dbPath, "dbFile", dbPath, "sqlite output file")
	return
}

func writeGroups(g generate.Generator, db modelWriter, outPath string) (err error) {
	for g.Next() {
		group := g.Name()
		if len(genFlags.dl) == 0 || (genFlags.dl == group) {
			if genFlags.useDB {
				if e := g.WriteTable(db); e != nil {
					err = e
					break
				}
			}
			// write to buffers first so that if there's an error we don't disrupt the existing source.
			var u bytes.Buffer
			if e := g.WriteSource(&u); e != nil {
				err = e
				break
			} else if f, e := format.Source(u.Bytes()); e != nil {
				log.Println(u.String())
				err = e
				break
			} else {
				path := filepath.Join(outPath, group)
				os.MkdirAll(path, os.ModePerm)

				if e := writeFile(path, group+"_types.go", f); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func writeFile(dir, file string, b []byte) (err error) {
	path := filepath.Join(dir, file)
	log.Println("writing", path)
	return os.WriteFile(path, b, 0666)
}

func readGroups(path string) (groups []generate.Group, err error) {
	inDir := os.DirFS(path)
	err = fs.WalkDir(inDir, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			err = e
		} else if n := d.Name(); !d.IsDir() {
			ext := files.Ext(n) // fix: ignores specs for now ( it has swaps still )
			if ext.Spec() && !strings.HasPrefix(n, "spec.") {
				if fp, e := inDir.Open(path); e != nil {
					err = e
				} else {
					defer fp.Close()
					var m map[string]any
					if e := files.FormattedRead(fp, ext, &m); e != nil {
						err = e
					} else if msg, e := decode.DecodeMessage(m); e != nil {
						err = e
					} else if g, e := generate.ReadSpec(msg); e != nil {
						err = e
					} else {
						groups = append(groups, g)
					}
				}
			}
		}
		if err != nil {
			err = errutil.New(err, "in", path)
		}
		return
	})
	return
}

func createDB(create bool, outFile string) (ret modelWriter, err error) {
	if create {
		if outFile, e := filepath.Abs(outFile); e != nil {
			err = e
		} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
			err = errutil.New("couldn't clean output file", outFile, e)
		} else {
			log.Println("generating", outFile)
			os.MkdirAll(path.Dir(outFile), os.ModePerm) // 0777 -> ModePerm ... read/writable by all

			if db, e := tables.CreateIdl(outFile); e != nil {
				err = e
			} else {
				if tx, e := db.Begin(); e != nil {
					err = errutil.New("couldnt create transaction", e)
				} else {
					ret = modelWriter{db, tx}
				}
				if err != nil {
					db.Close()
				}
			}
		}
	}
	return
}
