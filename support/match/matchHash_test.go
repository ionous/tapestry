package match_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/match"
)

// Test that the finding single and multiple phrases works okay.
func TestMatching(t *testing.T) {
	// prefixes
	prefixes := []string{
		"the",                  // 0
		"match this please",    // 1
		"please",               // 2
		"please and thank you", // 3
		"apples",               // 4
	}
	// things to find in the prefixes
	tests := []string{
		"please and thank you",
		"the apple", // should match 'the'
		"APPLES",
		"please", // should match 'please' not 'please and ...'
		"humbug",
	}
	expect := []int{
		3,
		0,
		4,
		2,
		-1,
	}
	prefixList := match.PanicSpans(prefixes...)
	for i, w := range tests {
		if h, e := match.Tokenize(w); e != nil {
			t.Fatal(e)
		} else {
			matched, skip := prefixList.FindPrefixIndex(h)
			if skip == 0 { // shh...
				matched = -1
			}
			if matched != expect[i] {
				t.Fatalf("failed to match '%s', got %d instead", w, matched)
			}
		}
	}
}
