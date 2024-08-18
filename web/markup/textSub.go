package markup

import (
	"errors"
	"strconv"
)

// FormatWriter - generate replacement text for certain html-like tags.
type FormatWriter interface {
	WriteOpen(c *converter, attr string) error
	WriteClose(c *converter) error
}

// Format - generic implementation of FormatWriter
// If a substitution string exists, it replaces the opening and closing tags;
// otherwise the specified rune replaces the opening tag, and Close is ignored.
type Format struct {
	Sub   string // replaces the opening tag
	Close string // optionally, replaces the closing tag ( but only if sub is specified )
	Rune  rune   // when sub is empty
}

// Formatting - all of the unique replacement strings.
// Which html tags map to which format is hardcoded.
type Formatting map[string]FormatWriter

// Formats - the default set of replacement string values.
var Formats = Formatting{
	TagLink:      FormatFunc(linkFormat),
	TagBold:      Format{Sub: "**"},
	TagDivider:   Format{Sub: "\r-----------------------------\n"},
	TagItalic:    Format{Sub: "*"},
	TagStrike:    Format{Sub: "~~"},
	TagUnderline: Format{Sub: "__"},
	//
	TagParagraph: Format{Rune: Paragraph}, // p
	TagNewline:   Format{Rune: Newline},   // br
	TagSoftline:  Format{Rune: Softline},  // wbr
	//
	TagItem:          FormatFunc(listItem),
	TagOrderedList:   listFormat(true),
	TagUnorderedList: listFormat(false),
}

// Tabwidth - how many spaces should list indentation generate
var Tabwidth = 2

const (
	// Starts a new line of text.
	Newline = '\n'
	// Conditionally prints a single line of blank text;
	// doesn't write a blank line if there already is one.
	Paragraph = '\v'
	// Starts a new line only if the output isnt already at the start of a newline.
	Softline = '\r'
	Space    = ' '
)

const (
	TagLink          = "a"
	TagBold          = "b"
	TagDivider       = "hr"
	TagItalic        = "i"
	TagItem          = "li"
	TagStrike        = "s"
	TagUnderline     = "u"
	TagParagraph     = "p"
	TagNewline       = "br"
	TagSoftline      = "wbr"
	TagOrderedList   = "ol"
	TagUnorderedList = "ul"
)

// returns true if it recognized the tag
func (fs *Formatting) WriteTag(c *converter, tag, attr string) (okay bool, err error) {
	opening := tag[0] != '/'
	if !opening {
		tag = tag[1:]
	}
	if g, ok := Formats[tag]; ok {
		if opening {
			err = g.WriteOpen(c, attr)
		} else {
			err = g.WriteClose(c)
		}
		okay = true
	}
	return
}

func (g Format) WriteClose(c *converter) (err error) {
	if len(g.Close) > 0 {
		_, err = c.WriteString(g.Close)
	} else if len(g.Sub) > 0 {
		_, err = c.WriteString(g.Sub)
	}
	return
}

func (g Format) WriteOpen(c *converter, attr string) (_ error) {
	if len(g.Sub) > 0 {
		c.WriteString(g.Sub)
	} else if g.Rune != 0 {
		c.WriteRune(g.Rune)
	}
	return
}

// function based implementation of FormatWriter interface
type FormatFunc func(c *converter, attr string, open bool) error

func (fn FormatFunc) WriteOpen(c *converter, attr string) error {
	return fn(c, attr, true)
}

func (fn FormatFunc) WriteClose(c *converter) error {
	return fn(c, "", false)
}

// writer for <a="label"> links
// in the absence of some sort of style rules:
// writes the label unless the direct parent is an ordered list.
// ( the assumption is that the list index is enough to identify the link )
func linkFormat(c *converter, attr string, open bool) (_ error) {
	inOrderedList := func(c *converter) (okay bool) {
		if cnt := len(c.lists); cnt > 0 {
			okay = c.lists[cnt-1].ordered
		}
		return
	}
	if !open {
		c.WriteString(" ")
	} else if len(attr) > 0 && !inOrderedList(c) {
		c.WriteString(attr)
		c.WriteString(": ")
	}
	return
}

func listFormat(ordered bool) FormatFunc {
	return func(c *converter, attr string, open bool) (err error) {
		if open {
			list := listHistory{ordered: ordered}
			c.lists = append(c.lists, list)
		} else if cnt := len(c.lists); cnt > 0 {
			c.lists = c.lists[:cnt-1]
		} else {
			err = errors.New("too many list closing tags")
		}
		return
	}
}

func listItem(c *converter, attr string, open bool) (err error) {
	if cnt := len(c.lists); cnt == 0 {
		err = errors.New("li, but there are no lists")
	} else if !open {
		c.WriteString("\r")
	} else if el := &c.lists[cnt-1]; !el.ordered {
		c.WriteString("\r- ")
	} else {
		el.itemCount++
		c.WriteRune('\r')
		c.WriteString(strconv.Itoa(el.itemCount))
		c.WriteString(". ")
	}
	return
}
