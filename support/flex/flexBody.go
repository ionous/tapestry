package flex

import (
	"errors"
	"io"

	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/files"
)

type accum []story.StoryStatement

// the body alternates between plain text and structured story ops
func (a *accum) readBody(k *Section) (err error) {
	for story := false; err == nil && k.NextSection(); story = !story {
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

// read a story data
// fix: maybe it'd be nice if the structured sections
// could be of any uniform type ( same with plain text )
// some sort of callback collector instead of specifically story/jess.
func decodeStorySection(in io.RuneReader) (ret []story.StoryStatement, err error) {
	if msg, e := readTellSection(in); e != nil {
		err = e
	} else {
		var slots story.StoryStatement_Slots
		dec := story.NewDecoder() // fix:  reusable?
		if e := dec.Decode(&slots, msg); e != nil {
			err = e
		} else {
			ret = slots
		}
	}
	return
}

// read one or more values; presumably mappings.
func readTellSection(in io.RuneReader) (ret []any, err error) {
	if d, e := files.ReadTellRunes(in, true); e != nil {
		err = e
	} else {
		switch content := d.(type) {
		case map[string]any:
			// normalize one mapping into a series of values.
			ret = []any{content}
		case []any:
			// a series of tell values
			ret = content
		default:
			err = errors.New("expected one or more tell statements")
		}
	}
	return
}

// decode execute
func decodeExecute(msgs []any) (ret []rt.Execute, err error) {
	var slots rtti.Execute_Slots
	dec := story.NewDecoder() // fix:  reusable?
	if e := dec.Decode(&slots, msgs); e != nil {
		err = e
	} else {
		ret = slots
	}
	return
}

// decode execute
func decodeAssignment(msg map[string]any) (ret rt.Assignment, err error) {
	var out rtti.Assignment_Slot
	dec := story.NewDecoder() // fix:  reusable?
	if e := dec.Decode(&out, msg); e != nil {
		err = e
	} else {
		ret = out.Value
	}
	return
}
