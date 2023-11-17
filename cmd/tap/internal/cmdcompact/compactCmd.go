// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package cmdcompact

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/jsn/dout"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

const (
	SpecExt     = ".ifspecs"
	DetailedExt = ".ifx"
	CompactExt  = ".if"
	BlockExt    = ".block"
	//
	TellSpec  = ".tells"
	TellStory = ".tell"
)

var allExts = []string{SpecExt, DetailedExt, CompactExt, BlockExt, TellSpec, TellStory}

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

func runCompact(ctx context.Context, cmd *base.Command, args []string) (err error) {
	tgtExt := strings.ReplaceAll(compactFlags.outExt, "_", "")
	if len(tgtExt) != 0 && !files.IsValidExtension(tgtExt, allExts) {
		flag.Usage() // exits
	} else if len(compactFlags.inPath) == 0 || len(compactFlags.outPath) == 0 {
		flag.Usage() // exits
	}

	// call the  process function on every requested file.
	if e := files.ReadPaths(
		compactFlags.inPath,
		compactFlags.recurse,
		strings.Split(compactFlags.inExts, ","),
		// per file callback, called with a path
		func(path string) error { return process(path, tgtExt) },
	); e != nil {
		err = errutil.New("error processing files", e)
	}
	return
}

// called for every input file:
func process(inFile, tgtExt string) (err error) {
	// skip files we cant handle
	if !files.IsValidExtension(inFile, allExts) {
		return
	}
	// convert the filename
	var outName string
	if fileName := filepath.Base(inFile); len(compactFlags.outExt) == 0 {
		outName = fileName
	} else {
		fileExt := filepath.Ext(fileName)
		outName = fileName[:len(fileName)-len(fileExt)] + compactFlags.outExt
	}
	outFile := filepath.Join(compactFlags.outPath, outName)
	// specs can only become specs:
	if inExt := filepath.Ext(inFile); inExt == SpecExt {
		if tgtExt == SpecExt {
			if e := decodeEncodeSpec(inFile, outFile, compactFlags.format()); e != nil {
				err = errutil.New("couldnt process", inFile, "=>", outFile, e)
			}
		}
	} else {
		// story files can be converted from one format to another
		var x xform
		switch inExt {
		case DetailedExt:
			x.decode = detailed.decode
			x.encode = detailed.encode
		case CompactExt:
			x.decode = compact.decode
			x.encode = compact.encode
		case BlockExt:
			x.decode = blockly.decode
			x.encode = blockly.encode
		}
		// override the default encoding if any specified
		switch tgtExt {
		case DetailedExt:
			x.encode = detailed.encode
		case CompactExt:
			x.encode = compact.encode
		case BlockExt:
			x.encode = blockly.encode
		}
		// report on results:
		if e := x.decodeEncode(inFile, outFile, compactFlags.format()); e != nil {
			err = errutil.New("couldnt process", inFile, "=>", outFile, e)
		}
	}
	return // done processing
}

// decode and encode are not mirror images
type xform struct {
	// reads a byte stream
	decode func(jsn.Marshalee, []byte) error
	// generates a serializable object
	encode func(jsn.Marshalee) (any, error)
}

type format int

const (
	unindentedJson format = iota
	indentedJson
	useTellFormat
)

func (f format) write(out string, data any) (err error) {
	switch f {
	case unindentedJson:
		err = files.WriteJson(out, data, false)
	case indentedJson:
		err = files.WriteJson(out, data, true)
	case useTellFormat:
		err = files.WriteTell(out, data)
	default:
		panic("unknown format")
	}
	return
}

func (p *xform) decodeEncode(in, out string, format format) (err error) {
	var dst story.StoryFile
	if b, e := files.ReadFile(in); e != nil {
		err = e
	} else if e := p.decode(&dst, b); e != nil {
		err = e
	} else if e := customMigration(&dst); e != nil {
		err = e
	} else if data, e := p.encode(&dst); e != nil {
		err = e
	} else {
		err = format.write(out, data)
	}
	return
}

func decodeEncodeSpec(in, out string, format format) (err error) {
	if b, e := files.ReadFile(in); e != nil {
		err = e
	} else {
		var msg map[string]any
		if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else {
			var dst spec.TypeSpec
			if e := cin.Decode(&dst, msg, cin.Signatures(story.AllSignatures)); e != nil {
				err = e
			} else if data, e := cout.Encode(&dst, customSpecEncoder); e != nil {
				err = e
			} else {
				err = format.write(out, data)
			}
		}
	}
	return
}

var compact = xform{
	decode: func(dst jsn.Marshalee, b []byte) (err error) {
		var msg map[string]any
		if e := json.Unmarshal(b, &msg); e != nil {
			err = e
		} else {
			err = story.Decode(dst, msg, story.AllSignatures)
		}
		return
	},
	encode: func(src jsn.Marshalee) (any, error) {
		return cout.CustomEncode(src, cout.Handlers{
			Flow: customStoryFlow,
			Slot: customStorySlot,
		})
	},
}
var tellCompact = xform{
	decode: compact.decode,
	encode: compact.encode,
}
var detailed = xform{
	decode: func(dst jsn.Marshalee, b []byte) error {
		return din.Decode(dst, story.Registry(), b)
	},
	encode: func(src jsn.Marshalee) (any, error) {
		return dout.Encode(src)
	},
}
var blockly = xform{
	// turn a block file into some dl structures
	decode: func(dst jsn.Marshalee, b []byte) error {
		return unblock.Decode(dst, "story_file", tapestry.Registry(), b)
	},
	// turn some dl structures into a block file
	encode: func(src jsn.Marshalee) (ret any, err error) {
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

// example of migrating one command to another.
func customMigration(tgt jsn.Marshalee) (err error) {
	return
}

var CmdCompact = &base.Command{
	Run:       runCompact,
	Flag:      buildFlags(),
	UsageLine: "tap compact [-in path] [-out path]",
	Short:     "reformat story files",
	Long: `Transform and reformat various Tapestry file formats.

Known file types include:

	.if: default tapestry story file format
	.ifx: an older, more verbose story file format
	.ifspecs: idl files
	.block: blockly formatted story files

Example, reads the .if file and saves it out again after formatting.
	tap compact -in ../../stories/blank.if

Bulk conversion example:
	tap compact -in ../../stories/shared -out ../../stories/shared
`,
}
