package cmdcompact

import (
	"io"
	"os"

	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func writeSpec(path string, src *spec.TypeSpec) (err error) {
	var enc encode.Encoder
	if plainData, e := enc.Encode(src); e != nil {
		err = e
	} else {
		err = files.FormattedSave(path, plainData, compactFlags.pretty)
	}
	return
}

func writeError(path string, _ *story.StoryFile) error {
	return errutil.New("unhandled write")
}

func writeStory(path string, src *story.StoryFile) (err error) {
	if plainData, e := story.Encode(src); e != nil {
		err = e
	} else {
		err = files.FormattedSave(path, plainData, compactFlags.pretty)
	}
	return
}

func writeBlock(path string, src *story.StoryFile) (err error) {
	// load the typespecs on demand then cache them
	if blockTypes == nil {
		if ts, e := rs.FromSpecs(idl.Specs); e != nil {
			err = e
		} else {
			blockTypes = &ts
		}
	}
	if err == nil {
		if str, e := block.Convert(blockTypes, src); e != nil {
			err = e
		} else if fp, e := os.Create(path); e != nil {
			err = e
		} else {
			// blocks are always in their own json format
			// fix: really should hand it a stream
			_, err = io.WriteString(fp, str)
		}
	}
	return
}

var blockTypes *rs.TypeSpecs // cache of loaded typespecs
