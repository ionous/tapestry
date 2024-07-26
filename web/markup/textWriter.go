package markup

import (
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// listen for html-like markup written to the returned stream and write its plain text equivalent into the passed stream.
// ex. writing <b>bold</b> will output **bold** to w.
func ToText(w io.Writer) io.Writer {
	return &fsm{converter{out: w, newLineCount: 2}, readingText}
}

type fsm struct {
	c converter
	x state
}

func (a *fsm) Write(src []byte) (ret int, err error) {
	for p := src; len(p) > 0; {
		if q, width := utf8.DecodeRune(p); width == 0 {
			err = errors.New("error decoding text")
			break
		} else if next, e := a.x(&a.c, q); e != nil {
			err = fmt.Errorf("error parsing text %w", e)
		} else {
			a.x = next
			// advance to the next
			p = p[width:]
			ret += width
		}
	}
	if err == nil {
		err = a.c.flushOut(true)
	}
	return
}

// if we return 0, ret, and no error, it tries again.
type state func(c *converter, q rune) (ret state, err error)

func readingText(c *converter, q rune) (ret state, err error) {
	if q == '<' {
		ret = openingTag
		_, err = c.buf.WriteRune(q) // buffer the text
	} else {
		ret = readingText
		c.WriteRune(q) // write directly out
	}
	return
}

// parse the opening tag content
func openingTag(c *converter, q rune) (ret state, err error) {
	if _, e := c.buf.WriteRune(q); e != nil {
		err = e
	} else {
		switch {
		case q == '/': // ex. </b>
			if c.buf.Len() > 2 {
				// if there was text between the less and the slash...
				ret, err = rejectTag(c)
			} else {
				ret = closingTag
			}

		case q == '>':
			if c.buf.Len() <= 2 { // only <>?
				ret, err = rejectTag(c)
			} else {
				ret, err = dispatchTag(c, -1)
			}

		case q == '=':
			ret = tagAttr(c)

		case unicode.IsLetter(q):
			ret, err = accumTag(c, openingTag)

		default:
			ret, err = rejectTag(c)
		}
	}
	return
}

// starting at the equal sign, read until an ending quote.
// to keep this simple, no escaping, etc. supported
// ex. <a="something">
func tagAttr(c *converter) state {
	var sub int
	const (
		start = iota
		mid
		end
	)
	attr := c.buf.Len()
	var loop state
	loop = func(c *converter, q rune) (ret state, err error) {
		if _, e := c.buf.WriteRune(q); e != nil {
			err = e
		} else {
			switch sub {
			case start:
				if q != '"' {
					ret, err = rejectTag(c)
				} else {
					sub, ret = mid, loop
				}
			case mid:
				if q == '\\' {
					ret, err = rejectTag(c)
				} else if q != '"' {
					ret = loop
				} else {
					sub, ret = end, loop
				}
			case end:
				if q != '>' {
					ret, err = rejectTag(c)
				} else {
					ret, err = dispatchTag(c, attr)
				}
			}
		}
		return
	}
	return loop
}

// ex. </b> -- called after the leading slash
func closingTag(c *converter, q rune) (ret state, err error) {
	if _, e := c.buf.WriteRune(q); e != nil {
		err = e
	} else {
		switch {
		case q == '>':
			ret, err = dispatchTag(c, -1)

		case unicode.IsLetter(q): // ie. </
			ret, err = accumTag(c, closingTag)

		default: // wasn't a valid tag, revert to readingText
			ret, err = rejectTag(c)
		}
	}
	return
}

// still reading the tag?
func accumTag(c *converter, next state) (ret state, err error) {
	// if longer than the longest closing tag; its not a tag.
	if c.buf.Len() > len("</blockquote") {
		ret, err = rejectTag(c)
	} else {
		ret = next
	}
	return
}

// wasn't a valid tag, revert to readingText
func rejectTag(c *converter) (ret state, err error) {
	c.writeIndent()
	if _, e := io.WriteString(c.out, flush(&c.buf)); e != nil {
		err = e
	} else {
		ret = readingText
	}
	return
}

// done with the tag ( that's been accumulated into the converter )
// attrOfs is the offset in buffer of the equal sign
func dispatchTag(c *converter, attrOfs int) (ret state, err error) {
	var tag, attr string
	str := flush(&c.buf)
	if attrOfs < 0 {
		tag = str[1 : len(str)-1] // skip the bracket, but not the slash.
	} else {
		tag = str[1 : attrOfs-1]
		attr = str[attrOfs+1 : len(str)-2]
	}
	if ok, e := c.dispatchTag(tag, attr); e != nil {
		err = e
	} else if !ok {
		// it wasnt a tag after all
		c.writeIndent()
		_, err = io.WriteString(c.out, str)
	}
	ret = readingText
	return
}
