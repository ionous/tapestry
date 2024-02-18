package inflect

import (
	"github.com/ionous/num2words"
)

// for WordsToNum, see match.
func NumToWords(n int) (ret string, okay bool) {
	if s := num2words.Convert(int(n)); len(s) > 0 {
		ret, okay = s, true
	}
	return
}
