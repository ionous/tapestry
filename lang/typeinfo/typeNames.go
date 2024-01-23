package typeinfo

import (
	"strings"
	"unicode"
)

// change one_two into "One two"
// fix: move to inflect?
func FriendlyName(n string) (ret string) {
	var b strings.Builder
	cap := true
	for i, s := range strings.Split(n, "_") {
		if i > 0 {
			b.WriteRune(' ')
		}
		if !cap {
			b.WriteString(s)
		} else {
			for j, r := range s {
				if j == 0 {
					b.WriteRune(unicode.ToUpper(r))
				} else {
					b.WriteRune(r)
				}
			}
			// cap = capAllWords
		}
	}
	return b.String()
}
