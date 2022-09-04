package markup

import (
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/ionous/errutil"
)

// listen for html-like markup written to the returned stream and write its plain text equivalent into the passed stream.
// ex. writing <b>bold</b> will output **bold** to w.
func ToText(w io.Writer) io.Writer {
	return &fsm{converter{out: w, line: 2}, readingText}
}

type fsm struct {
	c converter
	x state
}

func (a *fsm) Write(src []byte) (ret int, err error) {
	for p := src; len(p) > 0; {
		if q, width := utf8.DecodeRune(p); width == 0 {
			err = errutil.New("error decoding text")
			break
		} else if next, e := a.x(&a.c, q); e != nil {
			err = errutil.New("error parsing text", e)
		} else {
			a.x = next
			// advance to the next
			p = p[width:]
			ret += width
		}
	}
	return
}

// if we return 0, ret, and no error, it tries again.
type state func(c *converter, q rune) (ret state, err error)

func readingText(c *converter, q rune) (ret state, err error) {
	if q == '<' {
		ret = openingTag
		_, err = c.buf.WriteRune(q)
	} else {
		ret = readingText
		c.WriteRune(q)
	}
	return
}

// parse the tag content
func openingTag(c *converter, q rune) (ret state, err error) {
	if _, e := c.buf.WriteRune(q); e != nil {
		err = e
	} else {
		switch {
		case q == '/': // ex. </b>
			if c.tag.Len() > 0 {
				ret, err = rejectTag(c)
			} else {
				ret = closingTag
			}

		case q == '>':
			ret, err = dispatchTag(c, true)

		case unicode.IsLetter(q):
			ret, err = accumTag(c, q, openingTag)

		default:
			ret, err = rejectTag(c)
		}
	}
	return
}

// ex. </b>
func closingTag(c *converter, q rune) (ret state, err error) {
	if _, e := c.buf.WriteRune(q); e != nil {
		err = e
	} else {
		switch {
		case q == '>':
			ret, err = dispatchTag(c, false)

		case unicode.IsLetter(q):
			ret, err = accumTag(c, q, closingTag)

		default: // wasn't a valid tag, revert to readingText
			ret, err = rejectTag(c)
		}
	}
	return
}

// still reading the tag
func accumTag(c *converter, q rune, next state) (ret state, err error) {
	if c.tag.Len() > len("blockquote") {
		ret, err = rejectTag(c)
	} else {
		_, err = c.tag.WriteRune(q)
		ret = next
	}
	return
}

// wasn't a valid tag, revert to readingText
func rejectTag(c *converter) (ret state, err error) {
	c.tag.Reset()
	// flush any pending data
	if _, e := c.buf.WriteTo(c.out); e != nil {
		err = e
	} else {
		ret = readingText
	}
	return
}

// done with the tag
func dispatchTag(c *converter, open bool) (ret state, err error) {
	if tag := c.tag.String(); len(tag) == 0 {
		_, err = c.buf.WriteTo(c.out) // flush
	} else if ok, e := c.dispatchTag(tag, open); e != nil {
		err = e
	} else if !ok {
		// it wasnt a tag after all
		_, err = c.buf.WriteTo(c.out)
	}
	c.tag.Reset()
	c.buf.Reset()
	ret = readingText
	return
}
