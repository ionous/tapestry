// Generates a playable database from a story file.
package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// ex. go run import.go -in /Users/ionous/Documents/Iffy/stories/shared -out /Users/ionous/Documents/Iffy/scratch/shared/play.db
func main() {
	var srcPath, outFile string
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	// var printStories bool
	// printStories:= flag.Bool("log", false, "write imported stories to console")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if len(outFile) == 0 {
		dir, _ := filepath.Split(srcPath)
		outFile = filepath.Join(dir, "play.db")
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	// lets do this in the dumbest of ways for now.
	var cat eph.Catalog // fix: capture "Dilemmas" and LogWarning?
	var writeErr error  // fix: this seems less than ideal; maybe writer should return err.
	k := story.NewImporter(collectEphemera(&cat, &writeErr), storyMarshaller)
	if e := distill(k, srcPath); e != nil {
		log.Fatalln(e)
	} else if writeErr != nil {
		log.Fatalln(writeErr)
	} else {
		log.Println("assembling....")
		if e := Assemble(&cat, outFile); e != nil {
			log.Fatalln(e)
		}
	}
}

// a terrible way to optimize database writes
type qel struct {
	tgt  string
	args []interface{}
}
type qels []qel

func (q *qels) Write(tgt string, args ...interface{}) (err error) {
	(*q) = append(*q, qel{tgt, args})
	return
}

func Assemble(cat *eph.Catalog, outFile string) (err error) {
	var queue qels
	var w eph.Writer = &queue

	// go process all of the ephemera
	if e := cat.AssembleCatalog(eph.PhaseActions{
		eph.AncestryPhase: eph.AncestryPhaseActions,
		eph.NounPhase:     eph.NounPhaseActions,
	}); e != nil {
		err = e
	} else if e := cat.WriteAspects(w); e != nil {
		err = e
	} else if e := cat.WriteDomains(w); e != nil {
		err = e
	} else if e := cat.WriteFields(w); e != nil {
		err = e
	} else if e := cat.WriteLocals(w); e != nil {
		err = e
	} else if e := cat.WriteKinds(w); e != nil {
		err = e
	} else if e := cat.WriteNouns(w); e != nil {
		err = e
	} else if e := cat.WriteNames(w); e != nil {
		err = e
	} else if e := cat.WritePatterns(w); e != nil {
		err = e
	} else if e := cat.WritePlurals(w); e != nil {
		err = e
	} else if e := cat.WriteDirectives(w); e != nil {
		err = e
	} else if e := cat.WriteRelations(w); e != nil {
		err = e
	} else if e := cat.WritePairs(w); e != nil {
		err = e
	} else if e := cat.WriteRules(w); e != nil {
		err = e
	} else if e := cat.WriteValues(w); e != nil {
		err = e
	} else {
		log.Println("writing", len(queue), "entries")
		if outFile, e := filepath.Abs(outFile); e != nil {
			err = e
		} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
			err = errutil.New("couldn't clean output file", outFile, e)
		} else {
			// 0755 -> readable by all but only writable by the user
			// 0700 -> read/writable by user
			// 0777 -> ModePerm ... read/writable by all
			os.MkdirAll(path.Dir(outFile), os.ModePerm)
			if db, e := sql.Open(tables.DefaultDriver, outFile); e != nil {
				err = errutil.New("couldn't create output file", outFile, e)
			} else {
				defer db.Close()
				if e := tables.CreateModel(db); e != nil {
					err = errutil.New("couldnt create model", e)
				} else if tx, e := db.Begin(); e != nil {
					err = errutil.New("couldnt create transaction", e)
				} else {
					for _, q := range queue {
						if _, e := tx.Exec(q.tgt, q.args...); e != nil {
							tx.Rollback()
							err = errutil.New("couldnt write to", q.tgt, e)
							break
						}
					}
					if err == nil {
						if e := tx.Commit(); e != nil {
							err = errutil.New("couldnt commit", e)
						}
					}
				}
			}
		}
	}
	return
}

func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(cat *eph.Catalog, out *error) story.WriterFun {
	// fix: needs to be more clever eventually...
	if e := cat.AddEphemera(
		eph.EphAt{
			At:  "asm",
			Eph: &eph.EphBeginDomain{Name: "entire_game"}}); e != nil {
		panic(e)
	}
	// built in kinds -- see ephKinds.go
	// fix? move to an .if file?
	kinds := []string{
		eph.KindsOfAction, eph.KindsOfPattern,
		eph.KindsOfAspect, "",
		eph.KindsOfEvent, eph.KindsOfPattern,
		eph.KindsOfKind, "",
		eph.KindsOfPattern, "",
		eph.KindsOfRecord, "",
		eph.KindsOfRelation, "",
	}
	for i := 0; i < len(kinds); i += 2 {
		k, p := kinds[i], kinds[i+1]
		if e := cat.AddEphemera(
			eph.EphAt{
				At:  "built in kinds",
				Eph: &eph.EphKinds{Kinds: k, From: p}}); e != nil {
			panic(e)
		}
	}
	var i int
	return func(el eph.Ephemera) {
		if e := cat.AddEphemera(eph.EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			*out = errutil.Append(*out, e)
		}
		i++ // temp
	}
}

func distill(k *story.Importer, srcPath string) (err error) {
	if srcPath, e := filepath.Abs(srcPath); e != nil {
		err = e
	} else if e := readPaths(k, srcPath); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	} else {
		k.Flush()
	}
	return
}

// read a comma-separated list of files and directories
func readPaths(k *story.Importer, filePaths string) (err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = errutil.Append(err, e)
		} else {
			which := readOne
			if info.IsDir() {
				which = readMany
			}
			if e := which(k, filePath); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func readMany(k *story.Importer, path string) error {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	return filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() {
			if ext := filepath.Ext(path); ext == CompactExt || ext == DetailedExt {
				if e := readOne(k, path); e != nil {
					err = errutil.New("error reading", path, e)
				}
			}
		}
		return
	})
}

func readOne(k *story.Importer, path string) (err error) {
	log.Println("reading", path)
	if fp, e := os.Open(path); e != nil {
		err = e
	} else {
		defer fp.Close()
		if b, e := io.ReadAll(fp); e != nil {
			err = e
		} else if script, e := decodeStory(path, b); e != nil {
			err = errutil.New("couldn't decode", path, "b/c", e)
		} else if e := k.ImportStory(path, script); e != nil {
			err = errutil.New("couldn't import", path, "b/c", e)
		}
	}
	return
}

func decodeStory(path string, b []byte) (ret *story.Story, err error) {
	var curr story.Story
	switch ext := filepath.Ext(path); ext {
	case CompactExt:
		if e := cin.Decode(&curr, b, iffy.AllSignatures); e != nil {
			err = e
		} else {
			ret = &curr
		}
	case DetailedExt:
		if e := din.Decode(&curr, iffy.Registry(), b); e != nil {
			err = e
		} else {
			ret = &curr
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
