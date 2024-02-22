package match_test

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/match"
)

// Test that the hashing properly strips spaces, lowers phrases, and produces the right ranges.
func TestHashing(t *testing.T) {
	phrases := []string{
		``,                                 // 0
		`a`,                                // 1 - single letter
		`match`,                            // 2
		`  a  b  `,                         // 3 - single letters
		` some stuff   you`,                // 4
		` trim spaces   yay  `,             // 5
		`APPLES`,                           // 6
		`    `,                             // 7
		`me, and her`,                      // 8
		`"quote together" you`,             // 9
		`unmatched " quote`,                // 10
		`sans full stop.`,                  // 11 - the full stop should end the sentence; but not be in the sentence
		`"quote stop."`,                    // 12- okay, fullstop before quote
		`nothing after. stops`,             // 13 - an error because more text after fullstop
		`"quote stop." here`,               // 14 - error
		"he said `\"quote unquote\"` what", // 15
		"unmatched ` tick",                 // 16
	}
	// expects that the hashes wind up matching the hashes of these exact strings.
	const expectsError = "<expects error>"
	expect := [][]string{
		nil,                            // 0
		{"a"},                          // 1
		{"match"},                      // 2
		{"a", "b"},                     // 3
		{"some", "stuff", "you"},       // 4
		{"trim", "spaces", "yay"},      // 5
		{"apples"},                     // 6
		nil,                            // 7
		{"me", ",", "and", "her"},      // 8
		{`"`, "quote together", "you"}, // 9
		{expectsError},                 // 10
		{"sans", "full", "stop"},       // 11
		{`"`, "quote stop."},           // 12
		{expectsError},                 // 13
		{expectsError},                 // 14
		{"he", "said", `"`, `"quote unquote"`, "what"}, // 15
		{expectsError}, // 16
	}
	if len(phrases) != len(expect) {
		panic("missing tests")
	}
	for i, w := range phrases {
		want := expect[i]
		wantsError := len(want) > 0 && want[0] == expectsError
		if have, e := match.MakeSpan(w); wantsError != (e != nil) {
			t.Fatalf("test %d has unexpected error %v ( wanted %v )", i, e, expectsError)
		} else if e == nil {
			if len(have) != len(want) {
				t.Fatal(i, "mismatch len; have:", len(have))
			} else {
				for j, s := range want {
					el := have[j]
					if a, b := el.Hash(), match.Hash(s); a != b {
						t.Fatalf("test %d mismatched %v(%s)!= %v(%s)", i, a, el.String(), b, s)
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
	prefixList := match.PanicSpans(prefixes...)
	for i, w := range tests {
		h := match.PanicSpan(w)
		matched, skip := prefixList.FindPrefix(h)
		if skip == 0 { // shh...
			matched = -1
		}
		if matched != expect[i] {
			t.Fatalf("failed to match '%s', got %d instead", w, matched)
		}
	}
}
