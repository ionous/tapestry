package markup

// Format - replacement text for certain html-like tags.
// If a substitution string exists, it replaces the opening and closing tags;
// otherwise the specified rune replaces the opening tag, and Close is ignored.
type Format struct {
	Sub, Close string
	Rune       rune
}

// Formatting - all of the unique replacement strings.
// Which html tags map to which format is hardcoded.
type Formatting struct {
	Bold, Divider, Italic, Item, Strike, Underline,
	Paragraph, Newline, Softline Format
}

// Formats - the default set of replacement string values.
var Formats = Formatting{
	Bold:      Format{Sub: "**"},
	Divider:   Format{Sub: "\r-----------------------------\n"},
	Italic:    Format{Sub: "*"},
	Item:      Format{Sub: "\r- ", Close: "\r"},
	Strike:    Format{Sub: "~~"},
	Underline: Format{Sub: "__"},
	//
	Paragraph: Format{Rune: Paragraph}, // p
	Newline:   Format{Rune: Newline},   // br
	Softline:  Format{Rune: Softline},  // wbr
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

func (fs *Formatting) Select(tag, attr string) (ret Format, okay bool) {
	okay = true // provisionally
	switch tag {
	case "a":
		ret = Format{
			Sub:   attr + ": ",
			Close: " ", // adds a space after a tags. but will get eaten at eol
		}
	case "b":
		ret = fs.Bold
	case "hr":
		ret = fs.Divider
	case "i":
		ret = fs.Italic
	case "li":
		ret = fs.Item
	case "s":
		ret = fs.Strike
	case "u":
		ret = fs.Underline
	//
	case "p":
		ret = fs.Paragraph
	case "br":
		ret = fs.Newline
	case "wbr":
		ret = fs.Softline
	default:
		okay = false
	}
	return
}

// returns true if it recognized the tag
func (fs *Formatting) WriteTag(c *converter, tag, attr string, open bool) (okay bool, err error) {
	if g, ok := Formats.Select(tag, attr); ok {
		if !open && len(g.Close) > 0 {
			_, err = c.WriteString(g.Close)
		} else if len(g.Sub) > 0 {
			_, err = c.WriteString(g.Sub)
		} else if open && g.Rune != 0 {
			_, err = c.WriteRune(g.Rune)
		}
		okay = true
	}
	return
}
