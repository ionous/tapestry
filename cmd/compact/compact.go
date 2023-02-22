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
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/jsn/dout"
	"git.sr.ht/~ionous/tapestry/web/files"
	"github.com/ionous/errutil"
)

const (
	SpecExt     = ".ifspecs"
	DetailedExt = ".ifx"
	CompactExt  = ".if"
	BlockExt    = ".block"
)

var allExts = []string{SpecExt, DetailedExt, CompactExt, BlockExt}

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

func main() {
	var inPath, outPath, inExts, outExt string
	var pretty bool
	flag.StringVar(&outPath, "out", "", "output directory; required.")
	flag.StringVar(&inPath, "in", "", "input file(s) or paths(s) (comma separated)")
	flag.StringVar(&inExts, "filter", ".if",
		`extension(s) for directory scanning.
ignored if 'in' refers to a specific file`)
	flag.StringVar(&outExt, "convert", "",
		`an optional file extension to force a story format conversion (.if|.ifx|.block)
underscores are allowed to avoid copying over the original files. (._if, .if_, etc.)
( ex. if the in and out directories are the same.
if no extension is specified, the output format is the same as the import format.`)
	flag.BoolVar(&pretty, "pretty", false, "make the output somewhat human readable")
	flag.BoolVar(&errutil.Panic, "panic", false, "pa_nic on error?")
	flag.Parse()

	tgtExt := strings.ReplaceAll(outExt, "_", "")
	if len(tgtExt) != 0 && !files.IsValidExtension(tgtExt, allExts) {
		flag.Usage() // exits
	} else if len(inPath) == 0 || len(outPath) == 0 {
		flag.Usage() // exits
	}
	process := func(inFile string) (err error) {
		// skip files we cant handle
		if !files.IsValidExtension(inFile, allExts) {
			return
		}
		// convert the filename
		var outName string
		if fileName := filepath.Base(inFile); len(outExt) == 0 {
			outName = fileName
		} else {
			fileExt := filepath.Ext(fileName)
			outName = fileName[:len(fileName)-len(fileExt)] + outExt
		}
		outFile := filepath.Join(outPath, outName)
		// specs can only become specs:
		if inExt := filepath.Ext(inFile); inExt == SpecExt {
			if tgtExt == SpecExt {
				if e := decodeEncodeSpec(inFile, outFile, pretty); e != nil {
					err = errutil.New("couldnt process", inFile, "=>", outFile, e)
				}
			}
		} else {
			// story files can be converted from one format to another
			var x xform
			switch inExt {
			case DetailedExt:
				x.decode = detailed.decode
			case CompactExt:
				x.decode = compact.decode
			case BlockExt:
				x.decode = blockly.decode
			}
			switch tgtExt {
			case DetailedExt:
				x.encode = detailed.encode
			case CompactExt:
				x.encode = compact.encode
			case BlockExt:
				x.encode = blockly.encode
			}
			// report on results:
			if e := x.decodeEncode(inFile, outFile, pretty); e != nil {
				err = errutil.New("couldnt process", inFile, "=>", outFile, e)
			}
		}
		return // done processing
	}
	if e := files.ReadPaths(inPath, strings.Split(inExts, ","), process); e != nil {
		log.Fatal("error processing files", e)
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
	} else if data, e := cout.Encode(&dst, customSpecEncoder); e != nil {
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
		return cout.Encode(src, customStoryEncoder)
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
	func(src jsn.Marshalee) (ret interface{}, err error) {
		// load the typespecs on demand then cache them
		if ptypes == nil {
			if ts, e := rs.FromSpecs(idl.Specs); e != nil {
				err = e
			} else {
				ptypes = &ts
			}
		}
		if err == nil {
			ret, err = block.Convert(ptypes, src)
		}
		return
	},
}
var ptypes *rs.TypeSpecs // cache of loaded typespecs

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

// define  a custom spec encoder.
var customSpecEncoder cout.CustomFlow = nil
var customStoryEncoder = story.CompactEncoder

// example removing "trim" for underscore names
// func init() {
// 	customSpecEncoder = func(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
// 		switch op := flow.GetFlow().(type) {
// 		case *spec.FlowSpec:
// 			if op.Trim {
// 				if len(op.Terms) == 0 {
// 					panic("empty terms " + op.Name)
// 				}
// 				if op.Terms[0].Name != "" {
// 					panic("unexpected name " + op.Name + " " + op.Terms[0].Name)
// 				}
// 				if op.Terms[0].Key == "" {
// 					panic("unexpected key " + op.Name)
// 				} else {
// 					op.Terms[0].Name = op.Terms[0].Key
// 					op.Terms[0].Key = "_"
// 				}
// 				op.Trim = false
// 			}
// 		}
// 		// we haven't serialized it -- just poked at its memory
// 		return chart.Unhandled("no custom encoder")
// 	}
// }

// install a custom encoder to rewrite things
// func init() {
// 	customStoryEncoder = func(m jsn.Marshaler, flow jsn.FlowBlock) error {
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

const Description = //
`Transforms detailed format json files to their compact format, and back again.
`
const Example = `
ex. go run compact.go -in ../../stories/blank.ifx [-out ../../stories/]

bulk conversion examples:

from the generated .if files, generate the .ifx files:
	go build compact.go; for f in ../../stories/*.if; do ./compact -in $f; done;

	or, load and rewrite the .if files
	go build compact.go; for f in ../../stories/shared/*.if; do ./compact -pretty -in $f -out .if; done;
	go build compact.go; for f in ../regenspec/out/*.ifspecs; do ./compact -pretty -in $f -out $f; done;
	go build compact.go; for f in ../../stories/shared/*.if; do ./compact -pretty -in $f -out .block; done;
	go build compact.go; for f in ../../stories/shared/*.block; do ./compact -pretty -in $f -out .if; done;

windows:
	for %i in (..\..\stories\shared\*.if) do ( compact -pretty -in %i -out %i )
	for %i in (..\..\idl\*.ifspecs) do ( compact -pretty -in %i -out %i )
`

func init() {
	flag.Usage = func() {
		println(Description)
		flag.PrintDefaults()
		println(Example)
	}
}
