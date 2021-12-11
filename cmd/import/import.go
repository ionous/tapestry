// Generates ephemera from a story file.
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
	"strings"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// Import reads a json file (from the composer editor)
// and creates a new sqlite database of "ephemera".
// It uses package export's list of commands for parsing program statements.

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

// ex. go run import.go -in /Users/ionous/Documents/Iffy/stories/shared -out /Users/ionous/Documents/Iffy/scratch/shared/ephemera.db
func main() {
	var inFile, outFile string
	var printStories bool
	flag.StringVar(&inFile, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.BoolVar(&printStories, "log", false, "write imported stories to console")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if len(outFile) == 0 {
		dir, _ := filepath.Split(inFile)
		outFile = filepath.Join(dir, "ephemera.db")
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if xs, e := distill(outFile, inFile); e != nil {
		log.Fatalln(e)
	} else {
		if printStories {
			for _, x := range xs {
				pretty.Println(x)
			}
		}
		log.Println("Imported", outFile)
	}
}

func distill(outFile, inFile string) (ret []*story.Story, err error) {
	// fix: write to temp db file then copy the file on success?
	// currently stray files are left hanging around
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if outFile, e := filepath.Abs(outFile); e != nil {
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
			if e := tables.CreateEphemera(db); e != nil {
				err = errutil.New("couldn't create tables", outFile, e)
			} else {
				fs := make(Files)
				if e := fs.ReadPaths(inFile); e != nil {
					err = errutil.New("couldn't read file", inFile, e)
				} else {
					type qel struct {
						tgt  string
						args []interface{}
					}
					var queue []qel
					writer := func(q string, args ...interface{}) {
						queue = append(queue, qel{q, args})
					}
					k := story.NewImporter(writer, cout.Marshal)
					for path, data := range fs {
						log.Println("importing", path)
						if sptr, e := decodeStory(path, data); e != nil {
							err = errutil.Append(err, errutil.New("couldnt decode", path, "b/c", e))
						} else if e := k.ImportStory(path, sptr); e != nil {
							err = errutil.Append(err, errutil.New("couldnt import", path, "b/c", e))
						} else {
							ret = append(ret, sptr)
						}
					}
					// if this gets too big, could accumulate and flush every set size
					// or -- do some sort of fancy channel that does the same in the background
					// ( maybe half accumulating , half writing )
					log.Println("writing", len(queue), "entries")
					if tx, e := db.Begin(); e != nil {
						err = e
					} else {
						for _, q := range queue {
							if _, e := tx.Exec(q.tgt, q.args...); e != nil {
								tx.Rollback()
								err = e
								break
							}
						}
						if err == nil {
							if e := tx.Commit(); e != nil {
								err = e
							}
						}
					}
				}
			}
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

type Files map[string][]byte

// read a comma-separated list of files and directories
func (fs *Files) ReadPaths(filePaths string) (err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = errutil.Append(err, e)
		} else {
			if info.IsDir() {
				if e := fs.readMany(filePath); e != nil {
					err = errutil.Append(err, e)
				}
			} else {
				if b, e := readOne(filePath); e != nil {
					err = errutil.Append(err, e)
				} else {
					(*fs)[filePath] = b
				}
			}
		}
	}
	return
}

func (fs *Files) readMany(path string) (err error) {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() && filepath.Ext(path) == ".if" {
			if b, e := readOne(path); e != nil {
				err = errutil.New("error reading", path, e)
			} else {
				(*fs)[path] = b
			}
		}
		return
	})
	return
}

func readOne(filePath string) (ret []byte, err error) {
	log.Println("reading", filePath)
	if fp, e := os.Open(filePath); e != nil {
		err = e
	} else {
		ret, err = io.ReadAll(fp)
		fp.Close()
	}
	return
}
