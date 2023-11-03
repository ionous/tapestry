package charmed

import (
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/charm"
)

// match the string; calls the match function when a full match or a mismatch is detected
// sends the rune being processed ( so either the last of the string, or the mismatch )
// as well as the length of the string on a full match, or index of the mismatch.
func MatchString(str string, onMatch func(r rune, i int) charm.State) charm.State {
	var i int // index in str
	return charm.Self("match "+str, func(self charm.State, r rune) (ret charm.State) {
		// if the string is empty match returns -1, and size is 0
		if match, size := utf8.DecodeRuneInString(str[i:]); match != r {
			ret = onMatch(r, i)
		} else if i += size; i < len(str) {
			ret = self // loop
		} else if i == len(str) {
			ret = onMatch(r, i)
		}
		return
	})
}
