package jess

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/support/match"
)

// Convert from a string containing numbers as digits, or as words.
// For words, good up to "twenty".
func WordsToNum(s string) (ret int, okay bool) {
	// todo: investigate something like https://github.com/donna-legal/word2number
	if cnt, e := strconv.Atoi(s); e == nil {
		ret, okay = cnt, true
	} else {
		hash := match.Hash(s)
		for i, n := range smallNumbers {
			if hash == n[0].Hash() {
				ret, okay = i+1, true
				break
			}
		}
	}
	return
}

var smallNumbers = match.PanicSpans(
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
