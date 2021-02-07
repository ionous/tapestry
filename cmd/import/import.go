// Generates ephemera from a story file.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// Import reads a json file (from the composer editor)
// and creates a new sqlite database of "ephemera".
// It uses package export's list of commands for parsing program statements.

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
			var ds reader.Dilemmas
			defer db.Close()
			if e := tables.CreateEphemera(db); e != nil {
				err = errutil.New("couldn't create tables", outFile, e)
			} else {
				if storyFiles, e := readPath(inFile); e != nil {
					err = errutil.New("couldn't import  file", inFile, e)
				} else {
					k := story.NewImporter(db, ds.Report)
					for _, storyFile := range storyFiles {
						if v, e := storyFile.importStory(k); e != nil {
							err = errutil.Append(err, e)
						} else {
							ret = append(ret, v)
						}
					}

					reader.PrintDilemmas(log.Writer(), ds)
				}
			}
		}
	}
	return
}

// read a comma-separated list of files and directories
func readPath(filePaths string) (ret []storyFile, err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = e
		} else {
			if !info.IsDir() {
				if one, e := readOne(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, one)
				}
			} else {
				if many, e := readMany(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, many...)
				}
			}
		}
	}
	return
}

func readMany(path string) (ret []storyFile, err error) {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() && filepath.Ext(path) == ".if" {
			if one, e := readOne(path); e != nil {
				err = errutil.New("error reading", path, e)
			} else {
				ret = append(ret, one)
			}
		}
		return
	})
	return
}

func readOne(filePath string) (ret storyFile, err error) {
	log.Println("reading", filePath)
	if f, e := os.Open(filePath); e != nil {
		err = e
	} else {
		defer f.Close()
		var one reader.Map
		if e := json.NewDecoder(f).Decode(&one); e != nil && e != io.EOF {
			err = e
		} else {
			ret = storyFile{filePath, one}
		}
	}
	return
}

type storyFile struct {
	path string
	data reader.Map
}

func (fp *storyFile) importStory(k *story.Importer) (ret *story.Story, err error) {
	log.Println("importing", fp.path)
	return k.ImportStory(fp.path, fp.data)
}
