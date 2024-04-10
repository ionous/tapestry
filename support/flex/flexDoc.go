package flex

import (
	"errors"
	"io"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/support/files"
)

type Document struct {
	cnt int
	k   Section
}

const DefaultScene = "Tapestry"

// overwrites the story statements in the passed file
func ReadStory(name string, in Unreader, out *story.StoryFile) (err error) {
	var els accum
	if k := MakeSection(in); k.NextSection() {
		if e := els.readHeader(&k); e != nil {
			err = e
		} else if e := els.readBody(&k); e != nil {
			err = e
		} else {
			// lhs will have the leading comments, rhs everything else
			// if rhs is empty -- everything is comments ( and/or there's nothing )
			if lhs, rhs := splitStatements(els); len(rhs) == 0 {
				out.Statements = lhs
			} else {
				// ensure the file has a top level scene;
				// and move the right hand statements into it.
				var scene *story.DefineScene
				if n, ok := rhs[0].(*story.DefineScene); ok {
					n.Statements = append(n.Statements, rhs[1:]...)
					scene = n
				} else {
					scene = &story.DefineScene{
						Scene: &literal.TextValue{Value: name},
						RequireScenes: &literal.TextValues{Values: []string{
							DefaultScene,
						}},
						Statements: rhs,
					}
				}
				out.Statements = append(lhs, scene)
			}
		}
	}
	return
}

// split so that all leading comments are on the lhs;
// the first statement and everything after are on the rhs
func splitStatements(els []story.StoryStatement) (lhs, rhs []story.StoryStatement) {
	lhs = els
	for i, el := range els {
		if _, ok := el.(*story.Comment); !ok {
			lhs, rhs = els[:i], els[i:]
			break
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
			Lines: lines,
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
	if ops, e := ReadText(in); e != nil {
		err = e
	} else {
		(*a) = append((*a), ops...)
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
