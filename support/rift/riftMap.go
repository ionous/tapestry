package rift

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"git.sr.ht/~ionous/tapestry/support/rift/maps"
	"github.com/ionous/errutil"
)

type Mapping struct {
	doc    *Document
	depth  int
	key    Signature
	values maps.Builder
	CommentBlock
}

// maybe doc is a factory even?
func NewMapping(parent Collection, header string, depth int) *Mapping {
	doc := parent.Document()
	c := &Mapping{doc: doc, depth: depth, values: doc.MakeMap(doc.keepComments)}
	if doc.keepComments {
		c.keepComments = true
		if len(header) > 0 {
			c.comments.WriteString(header)
		}
	}
	return c
}

func (c *Mapping) Document() *Document {
	return c.doc
}

// accept a new value ( errors if the key hasn't yet been determined )
// fix: return error if the key already exists
func (c *Mapping) writeValue(val any) (err error) {
	if key, e := c.key.GetSignature(); e != nil {
		err = e
	} else {
		c.values = c.values.Add(key, val)
	}
	return
}

// a state that can parse one key:value pair
// maybe push the returned thingy
// return doc.PushCallback(depth, STATE, ent.finalizeEntry)
func (c *Mapping) NewEntry() charm.State {
	ent := riftEntry{Collection: c, depth: c.depth + 2, pendingValue: computedValue{}}
	next := charm.Statement("map entry", func(r rune) (ret charm.State) {
		switch r {
		case Hash:
			// fix fix: header comments
			panic("not implemented")

		default:
			// key and after key:
			ret = charm.RunStep(r, &c.key, charm.Statement("after key", func(r rune) charm.State {
				// unlike sequence, we need to hand off the first character that isnt the key
				return ContentsLoop(&ent).NewRune(r)
			}))
			return
		}
	})
	return c.Document().PushCallback(ent.depth, next, ent.finalizeEntry)
}

// used by parent collections to read the completed collection
func (c *Mapping) FinalizeValue() (ret any, err error) {
	if c.key.IsKeyPending() {
		err = errutil.New("signature must end with a colon")
	} else {
		// write the comment block
		if c.keepComments {
			comment := strings.TrimRightFunc(c.comments.String(), unicode.IsSpace)
			c.values = c.values.Add("", comment)
		}
		ret, c.values = c.values.Map(), nil
	}
	return
}
