package markup

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

type converter struct {
	out          io.Writer
	buf          strings.Builder
	newLineCount int // current vertical spacing ( ie. are we at the start of a new line )
	spaces       int // leading spaces
	// for lists
	lists []listHistory
}

type listHistory struct {
	ordered   bool
	itemCount int // item count
}

// return the contents of buf as a string, and reset it.
func flush(buf *strings.Builder) string {
	str := buf.String()
	buf.Reset()
	return str
}

// closing tags start with a forward slash (/)
func (c *converter) dispatchTag(tag, attr string) (okay bool, err error) {
	return Formats.WriteTag(c, tag, attr)
}

// watches for newlines, etc.
func (c *converter) Write(b []byte) (ret int, err error) {
	for p := b; len(p) > 0; {
		if q, width := utf8.DecodeRune(p); width == 0 {
			err = errors.New("error decoding runes")
			break
		} else {
			c.WriteRune(q)
			p = b[width:]
		}
	}
	return
}

// watches for newlines, etc.
func (c *converter) WriteString(s string) (ret int, err error) {
	for _, q := range s {
		if n, e := c.WriteRune(q); e != nil {
			err = e
			break
		} else {
			ret += n
		}
	}
	return
}

// watches for newlines, etc.
func (c *converter) WriteRune(q rune) (ret int, err error) {
	ret = 1 // provisionally
	switch q {
	case Newline:
		c.spaces = 0
		writer.WriteRune(c.out, Newline)
		c.newLineCount++

	case Paragraph:
		c.spaces = 0
		for c.newLineCount < 2 {
			writer.WriteRune(c.out, Newline)
			c.newLineCount++
		}

	case Softline:
		c.spaces = 0
		for c.newLineCount < 1 {
			writer.WriteRune(c.out, Newline)
			c.newLineCount++
		}

	case Space:
		c.spaces++

	default:
		c.writeIndent()
		ret, err = writer.WriteRune(c.out, q)
	}
	return
}

func (c *converter) writeIndent() {
	// we arent writing a line ( caught by the other cases )
	// but we might be writing something just *after* a line
	// in which case, we should write any pending indentation
	// ( and note: the very first line starts at line count 2 )
	if c.newLineCount > 0 {
		listIndent := len(c.lists)
		for range listIndent * Tabwidth {
			c.out.Write([]byte{Space})
		}
		c.newLineCount = 0
	} else if c.spaces > 0 {
		for range c.spaces {
			c.out.Write([]byte{Space})
		}
	}
	c.newLineCount = 0
	c.spaces = 0
}

func (c *converter) flushOut(eatTrailing bool) (err error) {
	str := flush(&c.buf)
	if !eatTrailing || len(str) > 0 {
		c.writeIndent()
	}
	_, err = io.WriteString(c.out, str)
	return
}
