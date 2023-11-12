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
func NewMapping(doc *Document, header string, depth int) *Mapping {
	c := &Mapping{doc: doc, depth: depth, values: doc.MakeMap(doc.keepComments)}
	if doc.keepComments {
		c.keepComments = true
		c.comments.WriteString(header)
	}
	return c
}

// a state that can parse one key:value pair
// maybe push the returned thingy
// return doc.PushCallback(depth, STATE, ent.finalizeEntry)
func (c *Mapping) NewEntry() charm.State {
	ent := riftEntry{
		doc:          c.doc,
		depth:        c.depth + 2,
		pendingValue: computedValue{},
		addsValue: func(val any, comment string) (err error) {
			if key, e := c.key.GetSignature(); e != nil {
				err = e
			} else {
				c.values = c.values.Add(key, val)
				c.comments.WriteString(comment)
				c.comments.WriteRune(Record)
			}
			return
		},
	}
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
	return c.doc.PushCallback(ent.depth, next, ent.finalizeEntry)
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
