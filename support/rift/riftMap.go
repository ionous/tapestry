package rift

import (
	"strings"
	"unicode"

	"github.com/ionous/errutil"
)

type Mapping struct {
	doc    *Document
	depth  int
	key    Signature
	values MapValues
	CommentBlock
}

// maybe doc is a factory even?
func NewMapping(parent Collection, header string, depth int) *Mapping {
	doc := parent.Document()
	c := &Mapping{doc: doc, depth: depth}
	if doc.keepComments {
		c.keepComments = true
		c.values = make(MapValues, 1)
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
func (c *Mapping) WriteValue(val any) (err error) {
	if key, e := c.key.GetSignature(); e != nil {
		err = e
	} else {
		c.values.Add(key, val)
	}
	return
}

// used by parent collections to read the completed collection
func (c *Mapping) FinalizeValue() (ret any, err error) {
	if c.key.IsKeyPending() {
		err = errutil.New("signature must end with a colon")
	} else {
		// write the comment block
		if c.keepComments {
			comment := strings.TrimRightFunc(c.comments.String(), unicode.IsSpace)
			c.values[0] = MapValue{Value: comment}
		}
		ret, c.values = c.values, nil
	}
	return
}
