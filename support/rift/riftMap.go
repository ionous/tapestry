package rift

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
)

type Mapping struct {
	doc      *Document
	depth    int
	sig      Signature
	values   MapValues
	comments strings.Builder
}

// maybe doc is a factory even?
func NewMapping(parent Collection, header string, depth int) *Mapping {
	doc := parent.Document()
	out := &Mapping{doc: doc, depth: depth}
	keepComments := doc.KeepComments
	if keepComments {
		out.values = make(MapValues, 1)
	}
	if len(header) > 0 {
		out.comments.WriteString(header)
	}
	doc.PushCallback(depth, out, func() (err error) {
		// see if there was a value-less key in the pipeline
		// ex. "signature:<eof>"
		if sig, e := out.sig.GetSignature(); e != nil {
			err = e
		} else {
			if len(sig) > 0 {
				out.values.Add(sig, nil)
			}
			// write the comment block
			if keepComments {
				comment := strings.TrimRightFunc(out.comments.String(), unicode.IsSpace)
				out.values[0] = MapValue{Value: comment}
			}
			err = parent.WriteValue(out.values)
		}
		return
	})
	return out
}

func (c *Mapping) Document() *Document {
	return c.doc
}

// fix: return error if already written
func (c *Mapping) WriteValue(val any) (_ error) {
	c.values[len(c.values)-1].Value = val
	return
}

func (c *Mapping) CommentWriter() RuneWriter {
	return &c.comments
}

func (c *Mapping) NewRune(r rune) (ret charm.State) {
	if r == Hash {
		// fix fix: header comments
		panic("not implemented")
	} else {
		ret = charm.RunStep(r, &c.sig, charm.Statement("after sig", func(r rune) (ret charm.State) {
			if sig, e := c.sig.GetSignature(); e != nil {
				ret = charm.Error(e)
			} else {
				// add a nil placeholder value; alt: trigger the pop() write every time
				// rather than just the last time ( ie. pop before returning the indented state )
				c.values.Add(sig, nil)
				ret = CollectionEntries(c, c, c.depth)
			}
			return
		}))
	}
	return
}
