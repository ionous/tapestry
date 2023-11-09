package rift_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/rift"
)

func xTestDoc(t *testing.T) {
	testValue(t,
		// -----------------------
		"x test multiple collections", `
cartParams:
  - name: "Leg type"
    options:
      - "Round"
      - "Square"
  - name: "Leg width"
    options: 
      - "1.00 inches"
      - "2.00 inches"`,
		rift.MapValues{
			{"cartParams:", []any{
				rift.MapValues{{
					"name:", "Leg type",
				}, {
					"options:", []any{
						"Round",
						"Square",
					},
				}}, rift.MapValues{{
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
		"x test multiple sub values", `
- -
  - 5
`, []any{[]any{nil, 5.0}},
		// -----------------------
		"test map with nil value", `
- Field:
  Next: 5`,
		[]any{
			rift.MapValues{
				{"Field:", nil},
				{"Next:", 5.0},
			}},

		// -----------------------
		"x test nested maps", `
- Field:
    Next: 5`,
		[]any{
			rift.MapValues{
				{"Field:", rift.MapValues{
					{"Next:", 5.0},
				}},
			}},

		// -----------------------
		// in yaml, inline nested maps are invalid
		// should they be here too?
		// to do, i think Value would need to examine history
		// either sniffing prior types or through a flag (ex. require newlines)
		// that it can send into NewMapping
		"x test inline maps", `
- Field: Next: 5`,
		[]any{
			rift.MapValues{{
				"Field:", rift.MapValues{{
					"Next:", 5.0,
				}},
			}}},
	)
}
