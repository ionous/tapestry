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
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/jsn/dout"
	"github.com/ionous/errutil"
)

const (
	SpecExt     = ".ifspecs"
	DetailedExt = ".ifx"
	CompactExt  = ".if"
)

func oppositeExt(ext string) (ret string) {
	if ext == CompactExt {
		ret = DetailedExt
	} else if ext == DetailedExt {
		ret = CompactExt
	} else {
		ret = ext
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
// go build compact.go; for f in ../regenspec/out/*.ifspecs; do ./compact -in $f -out $f; done;
//
func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (.if|.ifx|.ifspecs)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.if|.ifx|.ifspecs)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		println("requires an input file")
	} else {
		inExt := filepath.Ext(inFile)
		if !strings.HasPrefix(inExt, CompactExt) {
			println("requires some sort of .if, .ifx, or .ifspecs file")
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
			if inExt == SpecExt {
				// report on results:
				if e := decodeEncodeSpec(inFile, outFile); e != nil {
					println(e.Error())
				} else {
					println("done.")
				}
			} else {
				var x xform
				if inExt == DetailedExt {
					x.decode = detailed.decode
				} else {
					x.decode = compact.decode
				}
				if outExt == DetailedExt {
					x.encode = detailed.encode
				} else {
					x.encode = compact.encode
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
}

type xform struct {
	decode func(jsn.Marshalee, []byte) error
	encode func(jsn.Marshalee) (interface{}, error)
}

func (p *xform) decodeEncode(in, out string) (err error) {
	var dst story.StoryFile
	if b, e := readOne(in); e != nil {
		err = e
	} else if e := p.decode(&dst, b); e != nil {
		err = e
	} else if e := xformStory(&dst); e != nil {
		err = e
	} else if data, e := p.encode(&dst); e != nil {
		err = e
	} else {
		err = writeOut(out, data)
	}
	return
}

func decodeEncodeSpec(in, out string) (err error) {
	var dst spec.TypeSpec
	if b, e := readOne(in); e != nil {
		err = e
	} else if e := cin.Decode(&dst, b, cin.Signatures(tapestry.AllSignatures)); e != nil {
		err = e
	} else if data, e := cout.Encode(&dst, nil); e != nil {
		err = e
	} else {
		err = writeOut(out, data)
	}
	return
}

var compact = xform{
	func(dst jsn.Marshalee, b []byte) error {
		return story.Decode(dst, b, tapestry.AllSignatures)
	},
	func(src jsn.Marshalee) (interface{}, error) {
		return cout.Encode(src, story.CompactEncoder)
	},
}
var detailed = xform{
	func(dst jsn.Marshalee, b []byte) error {
		return din.Decode(dst, tapestry.Registry(), b)
	},
	func(src jsn.Marshalee) (interface{}, error) {
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

// example of migrating one command to another.
func xformStory(tgt jsn.Marshalee) (err error) {
	// ts := chart.MakeEncoder()
	// err = ts.Marshal(tgt, story.Map(&ts, story.BlockMap{
	// 	story.OtherBlocks: story.KeyMap{
	// 		story.BlockStart: func(b jsn.Block, v interface{}) (err error) {
	// 			switch newBlock := b.(type) {
	// 			case jsn.FlowBlock:
	// 				f := newBlock.GetFlow()
	// 				if b, ok := f.(interface{ RewriteActivity() }); ok {
	// 					b.RewriteActivity()
	// 				}
	// 			}
	// 			return
	// 		},
	// 	},
	// }))
	return
}
