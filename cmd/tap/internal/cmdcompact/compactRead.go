package cmdcompact

import (
	"io/fs"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

func readSpec(fsys fs.FS, path string, out *spec.TypeSpec) (err error) {
	if msg, e := files.FormattedRead(fsys, path); e != nil {
		err = e
	} else if e := cin.Decode(out, msg, cin.Signatures(story.AllSignatures)); e != nil {
		err = e
	}
	return
}

func readError(fsys fs.FS, path string, _ *story.StoryFile) error {
	return errutil.New("unhandled read")
}

func readStory(fsys fs.FS, path string, out *story.StoryFile) (err error) {
	if msg, e := files.FormattedRead(fsys, path); e != nil {
		err = e
	} else {
		err = story.DecodeMessage(out, msg)
	}
	return
}

func readBlock(fsys fs.FS, path string, out *story.StoryFile) (err error) {
	if b, e := fs.ReadFile(fsys, path); e != nil {
		err = e
	} else {
		err = unblock.Decode(out, "story_file", tapestry.Registry(), b)
	}
	return
}
