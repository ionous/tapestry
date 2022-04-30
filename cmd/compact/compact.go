// Package main for 'compact'.
// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/prim"
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
	BlockExt    = ".block"
)

var exts = []string{SpecExt, DetailedExt, CompactExt, BlockExt}

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
// go build compact.go; for f in ../../stories/shared/*.if; do ./compact -pretty -in $f -out .if; done;
// go build compact.go; for f in ../regenspec/out/*.ifspecs; do ./compact -pretty -in $f -out $f; done;
// go build compact.go; for f in ../../stories/shared/*.if; do ./compact -pretty -in $f -out .block; done;
// go build compact.go; for f in ../../stories/shared/*.block; do ./compact -pretty -in $f -out .if; done;
//
// windows:
// for %i in (..\..\stories\shared\*.if) do ( compact -pretty -in %i -out %i )
func main() {
	var inFile, outFile string
	var pretty bool
	flag.StringVar(&inFile, "in", "", "input file name (.if|.ifx|.ifspecs)")
	flag.StringVar(&outFile, "out", "", "optional output file name (.if|.ifx|.ifspecs)")
	flag.BoolVar(&pretty, "pretty", false, "make the output somewhat human readable")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(inFile) == 0 {
		println("requires an input file")
	} else {
		inExt := filepath.Ext(inFile)
		var allExts strings.Builder
		var exists bool
		for _, x := range exts {
			if inExt == x {
				exists = true
				break
			}
			allExts.WriteString(x)
			allExts.WriteRune(' ')
		}
		if !exists {
			println("expected one of the file types:" + allExts.String())
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
			} else if ext := filepath.Ext(outFile); len(ext) == 0 || ext == outFile {
				// convert directory
				base := filepath.Base(inFile)
				outFile = filepath.Join(outFile, base[:len(base)-len(inExt)]+outExt)
			}
			// transform the files:
			if inExt == SpecExt {
				// report on results:
				if e := decodeEncodeSpec(inFile, outFile, pretty); e != nil {
					println(e.Error())
				} else {
					println("done.")
				}
			} else {
				var x xform
				switch inExt {
				case DetailedExt:
					x.decode = detailed.decode
				case CompactExt:
					x.decode = compact.decode
				case BlockExt:
					x.decode = blockly.decode
				}
				switch outExt {
				case DetailedExt:
					x.encode = detailed.encode
				case CompactExt:
					x.encode = compact.encode
				case BlockExt:
					x.encode = blockly.encode
				}
				// report on results:
				if e := x.decodeEncode(inFile, outFile, pretty); e != nil {
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

func (p *xform) decodeEncode(in, out string, pretty bool) (err error) {
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
		err = writeOut(out, data, pretty)
	}
	return
}

func decodeEncodeSpec(in, out string, pretty bool) (err error) {
	var dst spec.TypeSpec
	if b, e := readOne(in); e != nil {
		err = e
	} else if e := cin.Decode(&dst, b, cin.Signatures(tapestry.AllSignatures)); e != nil {
		err = e
	} else if data, e := cout.Encode(&dst, nil); e != nil {
		err = e
	} else {
		err = writeOut(out, data, pretty)
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
var blockly = xform{
	// turn a block file into some dl structures
	func(dst jsn.Marshalee, b []byte) error {
		return unblock.Decode(dst, "story_file", tapestry.Registry(), b)
	},
	// turn some dl structures into a block file
	func(src jsn.Marshalee) (interface{}, error) {
		return block.Convert(src)
	},
}

func writeOut(outPath string, data interface{}, pretty bool) (err error) {
	log.Println("writing", outPath)
	if fp, e := os.Create(outPath); e != nil {
		err = e
	} else {
		defer fp.Close()
		if str, ok := data.(string); ok {
			_, err = fp.Write(prettify(str, pretty))
		} else {
			js := json.NewEncoder(fp)
			js.SetEscapeHTML(false)
			if pretty {
				js.SetIndent("", "  ")
			}
			err = js.Encode(data)
		}
	}
	return
}

func prettify(str string, pretty bool) (ret []byte) {
	ret = []byte(str)
	if pretty {
		var indent bytes.Buffer
		if e := json.Indent(&indent, ret, "", "  "); e != nil {
			log.Println(e)
		} else {
			ret = indent.Bytes()
		}
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
	return
}

// // install a custom encoder.
// func init() {
// 	story.CompactEncoder = func(m jsn.Marshaler, flow jsn.FlowBlock) error {
// 		switch op := flow.GetFlow().(type) {
// 		case *story.AspectProperty:
// 			swap(&op.UserComment, &op.Comment)
// 		case *story.BoolProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.NumberProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.NumListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.RecordProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.RecordListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.TextListProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		case *story.TextProperty:
// 			swap(&op.UserComment, &op.NamedProperty.Comment)
// 		}
// 		return core.CompactEncoder(m, flow)
// 	}
// }

func swap(tgt *string, from *prim.Lines) {
	if len(*tgt) == 0 {
		*tgt = from.Str
		from.Str = ""
	}
}
