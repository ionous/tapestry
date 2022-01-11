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

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/jsn/dout"
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
//
// bulk conversion:
//
// from the generated .if files, generate the .ifx files:
// go build compact.go; for f in ../../stories/*.if; do ./compact -in $f; done;
//
// or, load and rewrite the .if files
// go build compact.go; for f in ../../stories/shared/*.if; do ./compact -in $f -out .if; done;
//
func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (.if|.ifx)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.if|.ifx)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		println("requires an input file")
	} else {
		inExt := filepath.Ext(inFile)
		if inExt != CompactExt && inExt != DetailedExt {
			println("requires an .if or .ifx file")
		} else {
			// determine the output extension
			// ( if nothing was specified, it will be the opposite of in )
			outExt := filepath.Ext(outFile)
			if len(outExt) == 0 {
				outExt = oppositeExt(inExt)
			} else if outExt == outFile {
				outFile = ""
			}
			// create outfile name if needed
			if len(outFile) == 0 {
				outFile = inFile[:len(inFile)-len(inExt)] + outExt
			} else if len(filepath.Ext(outFile)) == 0 {
				// convert directory
				base := filepath.Base(inFile)
				outFile = filepath.Join(outFile, base[:len(base)-len(inExt)]+outExt)
			}
			// transform the files:
			var x xform
			if inExt == CompactExt {
				x.decode = compact.decode
			} else {
				x.decode = detailed.decode
			}
			if outExt == CompactExt {
				x.encode = compact.encode
			} else {
				x.encode = detailed.encode
			}
			// report on results:
			if e := x.decodeEncode(inFile, outFile); e != nil {
				println(e.Error())
			} else {
				println("done.")
			}
		}
	}
}

type xform struct {
	decode func(*story.Story, []byte) error
	encode func(*story.Story) (interface{}, error)
}

func (p *xform) decodeEncode(in, out string) (err error) {
	var dst story.Story
	if b, e := readOne(in); e != nil {
		err = e
	} else if e := p.decode(&dst, b); e != nil {
		err = e
	} else if data, e := p.encode(&dst); e != nil {
		err = e
	} else {
		err = writeOut(out, data)
	}
	return
}

var compact = xform{
	func(dst *story.Story, b []byte) error {
		return story.Decode(dst, b, tapestry.AllSignatures)
	},
	func(src *story.Story) (interface{}, error) {
		return story.Encode(src)
	},
}
var detailed = xform{
	func(dst *story.Story, b []byte) error {
		return din.Decode(dst, tapestry.Registry(), b)
	},
	func(src *story.Story) (interface{}, error) {
		return dout.Encode(src)
	},
}

func writeOut(outPath string, data interface{}) (err error) {
	log.Println("writing", outPath)
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		js := json.NewEncoder(fp)
		js.SetEscapeHTML(false)
		js.SetIndent("", "  ")
		err = js.Encode(data)
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
