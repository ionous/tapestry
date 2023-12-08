package cmdcompact

import (
	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/jsn/din"
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

func readSpec(path string, out *spec.TypeSpec) (err error) {
	if msg, e := formatOf(path).read(path); e != nil {
		err = e
	} else if e := cin.Decode(out, msg, cin.Signatures(story.AllSignatures)); e != nil {
		err = e
	}
	return
}

func readError(path string, _ *story.StoryFile) error {
	return errutil.New("unhandled read")
}

func readStory(path string, out *story.StoryFile) (err error) {
	if msg, e := formatOf(path).read(path); e != nil {
		err = e
	} else {
		err = story.Decode(out, msg, story.AllSignatures)
	}
	return
}

func readDetailed(path string, out *story.StoryFile) (err error) {
	if b, e := files.ReadFile(path); e != nil {
		err = e
	} else {
		err = din.Decode(out, story.Registry(), b)
	}
	return
}

func readBlock(path string, out *story.StoryFile) (err error) {
	if b, e := files.ReadFile(path); e != nil {
		err = e
	} else {
		err = unblock.Decode(out, "story_file", tapestry.Registry(), b)
	}
	return
}
