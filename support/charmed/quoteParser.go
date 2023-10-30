package charmed

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// parses the rhs of a single line quoted string
// user has to implement NewRune
type QuoteParser struct {
	runes charm.Runes
	err   error
}

func (p *QuoteParser) StateName() string {
	return "quotes"
}

// GetString returns the text including its surrounding quote markers.
func (p *QuoteParser) GetString() (ret string, err error) {
	if p.err != nil {
		err = p.err
	} else {
		ret = p.runes.String()
	}
	return
}

// scans until the matching quote marker is found
func (p *QuoteParser) ScanQuote(q rune) (ret charm.State) {
	const escape = '\\'
	return charm.Self("findMatchingQuote", func(self charm.State, r rune) (ret charm.State) {
		switch {
		case r == q:
			// eats the ending quote
			ret = charm.Terminal

		case r == escape:
			ret = charm.Statement("escape", func(r rune) (ret charm.State) {
				if x, ok := escapes[r]; !ok {
					p.err = errutil.Fmt("unknown escape sequence %q", r)
				} else {
					ret = p.runes.Accept(x, self)
				}
				return
			})
		case r != Eof:
			ret = p.runes.Accept(r, self) // loop...
		}
		return
	})
}

var escapes = map[rune]rune{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
}
