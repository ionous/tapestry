// Package flex reads tell files that are sectioned into alternating
// blocks of structured and plain text sections.
// The plain text sections are wrapped with commands and
// merged into the structured sections.
// The plain text sections can also "jump out" into structured sections
// on lines ending with colons.
package flex

import (
	"git.sr.ht/~ionous/tapestry/dl/story"
)

func ReadStory(in Unreader) (ret []story.StoryStatement, err error) {
	return ReadStorySource("", in)
}

func ReadStorySource(source string, in Unreader) (ret []story.StoryStatement, err error) {
	var els accum
	k := MakeSection(source, in)
	if e := els.readBody(&k); e != nil {
		err = e
	} else {
		ret = els
	}
	return
}
