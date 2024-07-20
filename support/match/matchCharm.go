package match

import (
	"errors"
	"fmt"

	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/runes"
)

// read whitespace until the end of line
// then find the first non whitespace character
// eof and other runes are considered an error.
func findIndent(minIndent int, outIndent *int) charm.State {
	spaces := -1
	return charm.Self("findIndent", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		case runes.Space:
			if spaces >= 0 {
				spaces++
			}
			ret = self
		case runes.Newline:
			spaces = 0
			ret = self
		case runes.Eof:
			e := errors.New("unexpected eof")
			ret = charm.Error(e)

		default:
			if spaces < 0 {
				e := errors.New("expected a newline")
				ret = charm.Error(e)
			} else if spaces < minIndent {
				e := fmt.Errorf("expected at least %d spaces", minIndent)
				ret = charm.Error(e)
			} else {
				*outIndent = spaces
				// returns unhandled
			}
		}
		return
	})
}

func eatWhitespace() charm.State {
	return charm.Self("whitespace", func(self charm.State, q rune) (ret charm.State) {
		switch q {
		case runes.Space, runes.Newline:
			ret = self
		default:
			// otherwise, not whitespace so return unhandled
			// ( dont handle Eof, we want the states themselves to be allowed to handle that. )
		}
		return
	})
}
