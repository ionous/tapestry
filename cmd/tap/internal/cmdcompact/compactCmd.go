// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package cmdcompact

import (
	"context"
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

func runCompact(ctx context.Context, cmd *base.Command, args []string) (err error) {

	tgtExt := strings.ReplaceAll(compactFlags.outExt, "_", "")
	if len(tgtExt) != 0 && !files.IsValidExtension(tgtExt, allExts) {
		flag.Usage() // exits
	} else if len(compactFlags.inPath) == 0 || len(compactFlags.outPath) == 0 {
		flag.Usage() // exits
	}
	process := func(inFile string) (err error) {
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
				if e := decodeEncodeSpec(inFile, outFile, compactFlags.pretty); e != nil {
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
			if e := x.decodeEncode(inFile, outFile, compactFlags.pretty); e != nil {
				err = errutil.New("couldnt process", inFile, "=>", outFile, e)
			}
		}
		return // done processing
	}
	if e := files.ReadPaths(compactFlags.inPath, compactFlags.recurse, strings.Split(compactFlags.inExts, ","), process); e != nil {
		err = errutil.New("error processing files", e)
	}
	return
}

type xform struct {
	decode func(jsn.Marshalee, []byte) error
	encode func(jsn.Marshalee) (interface{}, error)
}

func (p *xform) decodeEncode(in, out string, pretty bool) (err error) {
	var dst story.StoryFile
	if b, e := files.ReadFile(in); e != nil {
		err = e
	} else if e := p.decode(&dst, b); e != nil {
		err = e
	} else if e := xformStory(&dst); e != nil {
		err = e
	} else if data, e := p.encode(&dst); e != nil {
		err = e
	} else {
		err = files.WriteJson(out, data, pretty)
	}
	return
}

func decodeEncodeSpec(in, out string, pretty bool) (err error) {
	var dst spec.TypeSpec
	if b, e := files.ReadFile(in); e != nil {
		err = e
	} else if e := cin.Decode(&dst, b, cin.Signatures(story.AllSignatures)); e != nil {
		err = e
	} else if data, e := cout.Encode(&dst, customSpecEncoder); e != nil {
		err = e
	} else {
		err = files.WriteJson(out, data, pretty)
	}
	return
}

var compact = xform{
	func(dst jsn.Marshalee, b []byte) error {
		return story.DecodeJson(dst, b, story.AllSignatures)
	},
	func(src jsn.Marshalee) (interface{}, error) {
		return cout.CustomEncode(src, cout.Handlers{
			Flow: customStoryFlow,
			Slot: customStorySlot,
		})
	},
}
var detailed = xform{
	func(dst jsn.Marshalee, b []byte) error {
		return din.Decode(dst, story.Registry(), b)
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

// example of migrating one command to another.
func xformStory(tgt jsn.Marshalee) (err error) {
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

Example:
	tap compact -in ../../stories/blank.if

Bulk conversion example:
	tap compact -in ../../stories/shared -out ../../stories/shared
`,
}
