// Generates a playable database from a story file.
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strconv"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/asm"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/files"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// ex. go run asm.go -in /Users/ionous/Documents/Tapestry/stories/shared -out /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	var srcPath, outFile string
	var check bool
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&check, "check", false, "run check after importing?")
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
	if e := importStoryFiles(k, srcPath); e != nil {
		log.Fatalln(e)
	} else if writeErr != nil {
		log.Fatalln(writeErr)
	} else {
		log.Println("assembling....")
		if e := Assemble(&cat, outFile); e != nil {
			log.Fatalln(e)
		} else if check {
			if cnt, e := checkFile(outFile); e != nil {
				log.Fatalln(e)
			} else {
				log.Println("Checked", cnt, outFile)
			}
		}
	}
}

func checkFile(inFile string) (ret int, err error) {
	if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't open db", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else {
			filter := ""
			opt := qna.NewOptions()
			ret, err = qna.CheckAll(db, filter, opt, tapestry.AllSignatures)
		}
	}
	return
}

func Assemble(cat *eph.Catalog, outFile string) (err error) {
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
			err = asm.Assemble(cat, db)
		}
	}
	return
}

// fix: this is probably only core marshal here.
func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(cat *eph.Catalog, out *error) story.WriterFun {
	// fix: needs to be more clever eventually...
	if e := cat.AddEphemera(
		"asm",
		&eph.EphBeginDomain{Name: "entire_game"}); e != nil {
		panic(e)
	}
	// built in kinds -- see ephKinds.go
	for _, k := range kindsOf.DefaultKinds {
		pk := k.Parent()
		if e := cat.AddEphemera(
			"built in kinds",
			&eph.EphKinds{Kinds: k.String(), From: pk.String()}); e != nil {
			panic(e)
		}
	}
	var i int
	return func(el eph.Ephemera) {
		if e := cat.AddEphemera(strconv.Itoa(i), el); e != nil {
			*out = errutil.Append(*out, e)
		}
		i++ // temp
	}
}

// read a comma-separated list of files and directories
func importStoryFiles(k *story.Importer, srcPath string) (err error) {
	if e := files.ReadPaths(srcPath,
		[]string{CompactExt, DetailedExt}, func(p string) error {
			return readOne(k, p)
		}); e != nil {
		err = errutil.New("couldn't read file", srcPath, e)
	} else {
		k.Flush()
	}
	return
}

func readOne(k *story.Importer, path string) (err error) {
	log.Println("reading", path)
	if b, e := files.ReadFile(path); e != nil {
		err = e
	} else if script, e := decodeStory(path, b); e != nil {
		err = errutil.New("couldn't decode", path, "b/c", e)
	} else if e := k.ImportStory(path, script); e != nil {
		err = errutil.New("couldn't import", path, "b/c", e)
	}
	return
}

func decodeStory(path string, b []byte) (ret *story.Story, err error) {
	var curr story.Story
	switch ext := filepath.Ext(path); ext {
	case CompactExt:
		if e := story.Decode(&curr, b, tapestry.AllSignatures); e != nil {
			err = e
		} else {
			ret = &curr
		}
	case DetailedExt:
		if e := din.Decode(&curr, tapestry.Registry(), b); e != nil {
			err = e
		} else {
			ret = &curr
		}
	default:
		err = errutil.Fmt("unknown file type %q", ext)
	}
	return
}
