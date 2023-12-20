// Convert between various tapestry file formats.
package cmdcompact

import (
	"context"
	"flag"
	"os"
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
	if len(tgtExt) != 0 && isValidExtension(tgtExt) {
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

func isValidExtension(path string) bool {
	return files.Ext(path) > 0
}

// called for every input file:
func process(inFile, tgtExt string) (err error) {
	// skip files we cant handle
	if !isValidExtension(inFile) {
		return
	}
	// convert the filename
	fileSystem := os.DirFS(filepath.Dir(inFile))
	outName := compactFlags.replaceExt(filepath.Base(inFile))
	outFile := filepath.Join(compactFlags.outPath, outName)

	// specs can only become specs:
	if inExt, tgtExt := files.Ext(inFile), files.Ext(tgtExt); inExt.Spec() != tgtExt.Spec() {
		err = errutil.Fmt("can only change specs to specs. %s vs %s", inExt, tgtExt)
	} else if inExt.Spec() {
		var doc spec.TypeSpec
		if e := readSpec(fileSystem, inFile, &doc); e != nil {
			err = errutil.Fmt("%w while reading %s", e, inFile)
		} else if e := writeSpec(outFile, &doc); e != nil {
			err = errutil.Fmt("%w while writing %s", e, outFile)
		}
	} else {
		// pick the reader:
		reader, writer := readError, writeError
		switch inExt {
		case files.CompactExt, files.TellStory:
			reader = readStory
		case files.DetailedExt:
			reader = readDetailed
		case files.BlockExt:
			reader = readBlock
		}
		// pick the writer:
		switch tgtExt {
		case files.CompactExt, files.TellStory:
			writer = writeStory
		case files.DetailedExt:
			writer = writeDetailed
		case files.BlockExt:
			writer = writeBlock
		}
		// read and write:
		var doc story.StoryFile
		if e := reader(fileSystem, inFile, &doc); e != nil {
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

	.if: json formatted story files
	.ifx: an older, more verbose story file format
	.tell: a compact, yaml like story format
	.block: blockly formatted story files
	.ifspecs: json formatted command descriptions.
	.tells: tell formatted command descriptions.

To reformat an existing file:
	tap compact -in ../../stories/blank.tell

To bulk convert files:
	tap compact -in ../../stories/shared -out ../../stories/shared
`,
}
