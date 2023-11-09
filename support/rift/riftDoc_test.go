package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift/imap"
)

func TestDoc(t *testing.T) {
	testValue(t,
		// -----------------------
		"Test multiple collections", `
cartParams:
  - name: "Leg type"
    options:
      - "Round"
      - "Square"
  - name: "Leg width"
    options: 
      - "1.00 inches"
      - "2.00 inches"`,
		imap.ItemMap{
			{"cartParams:", []any{
				imap.ItemMap{{
					"name:", "Leg type",
				}, {
					"options:", []any{
						"Round",
						"Square",
					},
				}}, imap.ItemMap{{
					"name:", "Leg width",
				}, {
					"options:", []any{
						"1.00 inches",
						"2.00 inches",
					}},
				},
			},
			},
		},
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
		// in yaml, inline nested maps are invalid
		// should they be here too?
		// to do, i think Value would need to examine history
		// either sniffing prior types or through a flag (ex. require newlines)
		// that it can send into NewMapping
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
