// Compile story files into a playable database.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	"git.sr.ht/~ionous/tapestry/support/asm"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

const ExampleUsage = /**/
"go run asm.go -in /Users/ionous/Documents/Tapestry/stories/shared -out /Users/ionous/Documents/Tapestry/build/play.db -check"

func main() {
	usage := flag.Usage
	flag.Usage = func() {
		usage()
		println("ex.", ExampleUsage)
	}
	var srcPath, outFile string
	var checkAll bool
	var checkOne string
	flag.StringVar(&srcPath, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&checkAll, "check", false, "run check after importing?")
	flag.StringVar(&checkOne, "run", "", "run check on one test after importing?")
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

	if e := asm.AssembleFolder(srcPath, outFile); e != nil {
		log.Fatalln(e)
	} else if checkAll || len(checkOne) > 0 {
		if cnt, e := asm.CheckOutput(outFile, checkOne); e != nil {
			log.Fatalln(e)
		} else {
			log.Println("Checked", cnt, outFile)
		}
	}
}
