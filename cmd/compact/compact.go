package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/export"
	"git.sr.ht/~ionous/iffy/export/compact"
	"github.com/ionous/errutil"
)

// ex. go run compact.go -in /Users/ionous/Dev/go/src/git.sr.ht/~ionous/iffy/stories/shared/light.if
func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (json)")
	// flag.StringVar(&outFile, "out", "", "output file name (json)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		inFile = "/Users/ionous/Dev/go/src/git.sr.ht/~ionous/iffy/stories/shared/light.if"
	}
	if e := CompactFile(inFile); e != nil {
		log.Fatalln(e)
	} else {
		log.Println("compacted", inFile, "into", outFile)
	}
}

func CompactFile(inpaths string) (err error) {
	var ds reader.Dilemmas
	dec := compact.NewCompacterReporter(ds.Report)
	for _, slats := range iffy.AllSlats {
		dec.AddTypes(slats)
	}
	dec.AddTypes(story.Slats)
	//
	if fps, e := export.ReadPaths(inpaths); e != nil {
		err = e
	} else {
		for _, fp := range fps {
			dec.SetSource(fp.Path)
			if i, e := dec.Compact(fp.Data); e != nil {
				err = e
				break
			} else {
				b, err := json.MarshalIndent(i, "", "  ")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Print(string(b))
			}
		}
	}
	//
	reader.PrintDilemmas(log.Writer(), ds)
	return
}
