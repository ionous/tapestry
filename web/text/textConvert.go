package text

import (
	"bytes"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/ionous/errutil"
)

type converter struct {
	out   io.Writer
	buf   bytes.Buffer
	tag   strings.Builder
	tabs, // desired horizontal spacing
	line int // current vertical spacing ( ie. are we at the start of a new line )
}

func (c *converter) dispatchTag(tag string, open bool) (okay bool, err error) {
	switch tag {
	case "ol", "ul":
		if open {
			c.tabs++
		} else {
			c.tabs--
		}
		okay = true
	default:
		okay, err = Formats.WriteTag(c, tag, open)
	}
	return
}

// watches for newlines, etc.
func (c *converter) Write(b []byte) (ret int, err error) {
	for p := b; len(p) > 0; {
		if q, width := utf8.DecodeRune(p); width == 0 {
			err = errutil.New("error decoding runes")
			break
		} else {
			c.WriteRune(q)
			p = b[width:]
		}
	}
	return
}

// watches for newlines, etc.
func (c *converter) WriteRune(q rune) (ret int, err error) {
	switch q {
	case Newline:
		ret = 1 // we know the size of the newline
		writeRune(c.out, Newline)
		c.line++

	case Paragraph:
		ret = 1 // we know the size of the rune
		for c.line < 2 {
			writeRune(c.out, Newline)
			c.line++
		}

	case Softline:
		ret = 1 // we know the size of the rune
		for c.line < 1 {
			writeRune(c.out, Newline)
			c.line++
		}

	default:
		// we arent writing a line ( caught by the other cases )
		// but we might be writing something just *after* a line
		// in which case, we should write tabs
		if c.line > 0 && c.tabs > 0 {
			for i := 0; i < c.tabs*Tabwidth; i++ {
				c.out.Write([]byte{Space})
			}
		}
		// write the rune
		ret, err = writeRune(c.out, q)
		c.line = 0
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

func writeRune(w io.Writer, q rune) (int, error) {
	var p [utf8.UTFMax]byte
	n := utf8.EncodeRune(p[:], q)
	return w.Write(p[:n])
}
