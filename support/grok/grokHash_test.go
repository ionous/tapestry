package grok_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

// Test that the hashing properly strips spaces, lowers phrases, and produces the right ranges.
func TestHashing(t *testing.T) {
	phrases := []string{
		``,                     // 0
		`a`,                    // 1 - single letter
		`match`,                // 2
		`  a  b  `,             // 3 - single letters
		` some stuff   you`,    // 4
		` trim spaces   yay  `, // 5
		`APPLES`,               // 6
		`    `,                 // 7
		`me, and her`,          // 8
		`"quote together" you`, // 9
		`unmatched " quote`,    // 10
		`sans full stop.`,      // 11
		`nothing after. stops`, // 12
	}
	// expects that the hashes wind up matching the hashes of these exact strings.
	const expectsError = "<expects error>"
	expect := [][]string{
		nil,                       // 0
		{"a"},                     // 1
		{"match"},                 // 2
		{"a", "b"},                // 3
		{"some", "stuff", "you"},  // 4
		{"trim", "spaces", "yay"}, // 5
		{"apples"},                // 6
		nil,                       // 7
		{"me", ",", "and", "her"}, // 8
		{"quote together", "you"}, // 9
		{expectsError},            // 10
		{"sans", "full", "stop"},  // 11
		{expectsError},            // 12
	}
	if len(phrases) != len(expect) {
		panic("missing tests")
	}
	for i, w := range phrases {
		want := expect[i]
		wantsError := len(want) > 0 && want[0] == expectsError
		if have, e := grok.MakeSpan(w); wantsError != (e != nil) {
			t.Fatal("unexpected error", e)
		} else if e == nil {
			if len(have) != len(want) {
				t.Fatal(i, "mismatch len; have:", len(have))
			} else {
				for j, s := range want {
					el := have[j]
					if el.Hash() != grok.Hash(s) {
						t.Fatal(i, "mismatch el")
					}
					// lastly: test that the recorded start and end indices are correct
					// ( except for ands... which are deliberately weird )
					if strings.IndexRune(s, ',') >= 0 {
						str := strings.ToLower(el.String())
						if str != s {
							t.Fatal(i, "mismatch str", str)
						}
					}
				}
			}
		}
	}
}

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
	prefixList := groktest.PanicSpans(prefixes...)
	for i, w := range tests {
		h := groktest.PanicSpan(w)
		matched, skip := prefixList.FindPrefix(h)
		if skip == 0 { // shh...
			matched = -1
		}
		if matched != expect[i] {
			t.Fatalf("failed to match '%s', got %d instead", w, matched)
		}
	}

}
