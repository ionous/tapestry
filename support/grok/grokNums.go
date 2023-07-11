package grok

import (
	"strconv"
)

// Currently good up to twenty.
// maybe https://github.com/donna-legal/word2number?
func WordsToNum(s string) (ret int, okay bool) {
	if cnt, e := strconv.Atoi(s); e == nil {
		ret, okay = cnt, true
	} else {
		hash := Hash(s)
		for i, n := range smallNumbers {
			if hash == n[0].Hash() {
				ret, okay = i+1, true
				break
			}
		}
	}
	return
}

var smallNumbers = PanicSpans(
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
	"ten",
	"eleven",
	"twelve",
	"thirteen",
	"fourteen",
	"fifteen",
	"sixteen",
	"seventeen",
	"eighteen",
	"nineteen",
	"twenty",
)
