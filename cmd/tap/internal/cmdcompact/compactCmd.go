// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package cmdcompact

import (
	"context"
	"flag"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

// entry point for "tap compact"
func runCompact(ctx context.Context, cmd *base.Command, args []string) (err error) {
	tgtExt := strings.ReplaceAll(compactFlags.outExt, "_", "")
	if len(tgtExt) != 0 && !files.IsValidExtension(tgtExt, allExts) {
		flag.Usage() // exits
	} else if len(compactFlags.inPath) == 0 || len(compactFlags.outPath) == 0 {
		flag.Usage() // exits
	}
	// call process() on every requested file.
	var count int
	if e := files.ReadPaths(
		compactFlags.inPath,
		compactFlags.recurse,
		strings.Split(compactFlags.inExts, ","),
		// per file callback, called with a path
		func(path string) error {
			count++
			return process(path, tgtExt)
		},
	); e != nil {
		err = errutil.New("error processing files", e)
	} else if count == 0 {
		err = errutil.New("no files")
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
	outName := compactFlags.replaceExt(filepath.Base(inFile))
	outFile := filepath.Join(compactFlags.outPath, outName)

	// specs can only become specs:
	inExt := filepath.Ext(inFile)
	if inSpec, outSpec := isSpecExt(inExt), isSpecExt(tgtExt); inSpec != outSpec {
		err = errutil.Fmt("can only change specs to specs. %s vs %s", inExt, tgtExt)
	} else if inSpec {
		var doc spec.TypeSpec
		if e := readSpec(inFile, &doc); e != nil {
			err = errutil.Fmt("%w while reading %s", e, inFile)
		} else if e := writeSpec(outFile, &doc); e != nil {
			err = errutil.Fmt("%w while writing %s", e, outFile)
		}
	} else {
		// pick the reader:
		reader, writer := readError, writeError
		switch inExt {
		case CompactExt, TellStory:
			reader = readStory
		case DetailedExt:
			reader = readDetailed
		case BlockExt:
			reader = readBlock
		}
		// pick the writer:
		switch tgtExt {
		case CompactExt, TellStory:
			writer = writeStory
		case DetailedExt:
			writer = writeDetailed
		case BlockExt:
			writer = writeBlock
		}
		// read and write:
		var doc story.StoryFile
		if e := reader(inFile, &doc); e != nil {
			err = errutil.Fmt("%w while reading %s", e, inFile)
		} else if e := writer(outFile, &doc); e != nil {
			err = errutil.Fmt("%w while writing %s", e, outFile)
		}
	}
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
