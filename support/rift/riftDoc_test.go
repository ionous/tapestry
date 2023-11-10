package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift/imap"
)

func TestDoc(t *testing.T) {
	testValue(t,
		// --------------
		"Test multiple sub values", `
- -
  - 5
`, []any{[]any{nil, 5.0}},
		// -----------------------
		"Test map with nil value", `
- Field:
  Next: 5`,
		[]any{
			imap.ItemMap{
				{"Field:", nil},
				{"Next:", 5.0},
			}},

		// -----------------------
		"Test nested maps", `
- Field:
    Next: 5`,
		[]any{
			imap.ItemMap{
				{"Field:", imap.ItemMap{
					{"Next:", 5.0},
				}},
			}},

		// -----------------------
		"Test inline maps", `
- Field: Next: 5`,
		[]any{
			imap.ItemMap{{
				"Field:", imap.ItemMap{{
					"Next:", 5.0,
				}},
			}}},
	)
}
