package flex

import (
	"errors"
	"io"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
)

type Document struct {
	cnt int
	k   Section
}

// overwrites the story statements in the passed file
func ReadStory(in Unreader, out *story.StoryFile) (err error) {
	var a accum
	if k := MakeSection(in); k.NextSection() {
		if e := a.readHeader(&k); e != nil {
			err = e
		} else if e := a.readBody(&k); e != nil {
			err = e
		} else {
			out.StoryStatements = a
		}
	}
	return
}

type accum []story.StoryStatement

// header can only contain comments
func (a *accum) readHeader(in io.RuneReader) (err error) {
	if lines, e := ReadComments(in); e != nil {
		err = e
	} else {
		(*a) = append((*a), &story.Comment{
			// hrm. why isnt lines ... lines?
			Lines: strings.Join(lines, "\n"),
		})
	}
	return
}

// the body alternates between story and plain text
func (a *accum) readBody(k *Section) (err error) {
	for story := true; err == nil && k.NextSection(); story = !story {
		if story {
			err = a.readStory(k)
		} else {
			err = a.readText(k)
		}
	}
	return
}

// read and record a story section
func (a *accum) readStory(in io.RuneReader) (err error) {
	if slots, e := decodeStorySection(in); e != nil {
		err = e
	} else {
		(*a) = append((*a), slots...)
	}
	return
}

// read and record a plain text section
func (a *accum) readText(in io.RuneReader) (err error) {
	if lines, e := ReadText(in); e != nil {
		err = e
	} else {
		// fix: record the line number1
		op := &story.DeclareStatement{
			Text: &literal.TextValue{
				Value: strings.Join(lines, "\n"),
			},
		}
		(*a) = append((*a), op)
	}
	return
}

// read a story file ( which just happens to actually be a section of a file )
func decodeStorySection(in io.RuneReader) (ret []story.StoryStatement, err error) {
	var slots story.StoryStatement_Slots
	dec := story.NewDecoder() // fix:  reusable?
	if msg, e := readTellSection(in); e != nil {
		err = e
	} else if e := dec.Decode(&slots, msg); e != nil {
		err = e
	} else {
		ret = slots
	}
	return
}

// read the tell data, and normalize it into a series of statements
func readTellSection(in io.RuneReader) (ret []any, err error) {
	if d, e := files.ReadTellRunes(in); e != nil {
		err = e
	} else {
		switch content := d.(type) {
		case map[string]any:
			// one tell block
			ret = []any{content}
		case []any:
			// possibly a series of tell statements
			ret = content
		default:
			err = errors.New("expected one or more tell statements")
		}
	}
	return
}
