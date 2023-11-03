package charmed

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// scans until the matching quote marker is found
func ScanQuote(q rune, useEscapes bool, onDone func(string)) (ret charm.State) {
	const escape = '\\'
	var out strings.Builder
	return charm.Self("findMatchingQuote", func(self charm.State, r rune) (ret charm.State) {
		switch {
		case r == q:
			onDone(out.String())
			ret = charm.Finished("quotes")

		case r == escape && useEscapes:
			ret = charm.Statement("escape", func(r rune) (ret charm.State) {
				if x, ok := escapes[r]; !ok {
					e := errutil.Fmt("unknown escape sequence %q", r)
					ret = charm.Error(e)
				} else {
					out.WriteRune(x)
					ret = self // loop...
				}
				return
			})

		case r == '\n':
			e := errutil.New("unexpected newline")
			ret = charm.Error(e)

		case r == charm.Eof:
			e := errutil.New("unexpected end of file")
			ret = charm.Error(e)

		default:
			out.WriteRune(r)
			ret = self // loop...
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
