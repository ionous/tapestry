// Package main for 'compact'.
// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/din"
	"git.sr.ht/~ionous/iffy/jsn/dout"
	"github.com/ionous/errutil"
)

const (
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

func oppositeExt(ext string) (ret string) {
	if ext == DetailedExt {
		ret = CompactExt
	} else {
		ret = DetailedExt
	}
	return
}

// ex. go run compact.go -in ../../stories/blank.ifx [-out ../../stories/]
func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (.if|.ifx)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.if|.ifx)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		println("requires an input file")
	} else {
		inext := filepath.Ext(inFile)
		if inext != CompactExt && inext != DetailedExt {
			println("requires an .if or .ifx file")
		} else {
			// create outfile name if needed
			if len(outFile) == 0 {
				outFile = inFile[:len(inFile)-len(inext)] + oppositeExt(inext)
			}
			// transform the files:
			var err error
			if outext := filepath.Ext(outFile); outext == inext {
				err = errutil.New("requires one file to be compact and the other detailed")
			} else if outext == CompactExt {
				err = compact(inFile, outFile)
			} else {
				err = expand(inFile, outFile)
			}
			// report on results:
			if err != nil {
				println(err.Error())
			} else {
				println("done.")
			}
		}
	}
}

func compact(inDetails, outCompact string) (err error) {
	var dst story.Story
	if b, e := readOne(inDetails); e != nil {
		err = e
	} else if e := din.Decode(&dst, iffy.Registry(), b); e != nil {
		err = e
	} else if data, e := cout.Encode(&dst); e != nil {
		err = e
	} else {
		err = writeOut(outCompact, data)
	}
	return
}

func expand(inCompact, outDetails string) (err error) {
	var dst story.Story
	if b, e := readOne(inCompact); e != nil {
		err = e
	} else if e := cin.Decode(&dst, b, iffy.AllSignatures); e != nil {
		err = e
	} else if data, e := dout.Encode(&dst); e != nil {
		err = e
	} else {
		err = writeOut(outDetails, data)
	}
	return
}

func writeOut(outPath string, data interface{}) (err error) {
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		out := json.NewEncoder(fp)
		out.SetIndent("", "  ")
		err = out.Encode(data)
	}
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
